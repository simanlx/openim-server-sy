package cloud_wallet

import (
	"context"
	"crazy_server/pkg/cloud_wallet/ncount"
	imdb "crazy_server/pkg/common/db/mysql_model/cloud_wallet"
	"crazy_server/pkg/common/log"
	pb "crazy_server/pkg/proto/cloud_wallet"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type NcountPay struct {
	count ncount.NCounter
}

func NewNcountPay() *NcountPay {
	return &NcountPay{
		count: ncount.NewCounter(),
	}
}

type PayResult struct {
	ErrMsg        string
	ErrCode       int
	NcountOrderID string
}

// 余额支付
func (np *NcountPay) payByBalance(operationId, payAccountID, ReceiveAccountId, MerOrderId, totalAmount string) *PayResult {
	var (
		resp = &PayResult{
			ErrCode: 0,
			ErrMsg:  "支付成功",
		}
	)
	req := &ncount.TransferReq{
		MerOrderId: MerOrderId,
		TransferMsgCipher: ncount.TransferMsgCipher{
			PayUserId:     payAccountID,
			ReceiveUserId: ReceiveAccountId,
			TranAmount:    totalAmount, //分转元
		},
	}

	escap := time.Now()
	transferResult, err := np.count.Transfer(req)
	log.Info(operationId, "transfer req", req, "耗费时间:", time.Since(escap))
	if err != nil {
		//这里是网络层面的错误
		log.Error(operationId, "调用第三方支付出现网络错误", err)
		resp.ErrMsg = "调用第三方支付出现网络错误"
		resp.ErrCode = 400
		return resp
	}
	// 备注； 现在一般httpcode是自有逻辑实现，不用在专门判断
	//=========================成功返回=========================
	if transferResult.ResultCode != ncount.ResultCodeSuccess {
		remark := "余额支付失败： "
		co, _ := json.Marshal(transferResult)
		err := imdb.CreateErrorLog(remark, operationId, MerOrderId, transferResult.ErrorMsg, transferResult.ErrorCode, string(co))
		if err != nil {
			log.Error(operationId, "创建错误日志失败", err)
		}
		resp.ErrCode = 400
	}
	resp.ErrMsg = transferResult.ErrorMsg
	resp.NcountOrderID = transferResult.NcountOrderId
	return resp
}

// 银行卡支付
func (np *NcountPay) payByBankCard(operationId, payAccountID, ReceiveAccountId, MerOrderId, totalAmount, BankProtocol, NotifyUrl string) *PayResult {
	var (
		resp = &PayResult{
			ErrCode: 0,
			ErrMsg:  "支付成功",
		}
	)
	//充值支付
	transferResult, err := np.count.QuickPayOrder(&ncount.QuickPayOrderReq{
		MerOrderId: MerOrderId,
		QuickPayMsgCipher: ncount.QuickPayMsgCipher{
			PayType:       "3", //绑卡协议号充值
			TranAmount:    totalAmount,
			NotifyUrl:     NotifyUrl,
			BindCardAgrNo: BankProtocol,
			ReceiveUserId: ReceiveAccountId, //收款账户
			UserId:        payAccountID,
			SubMerchantId: ncount.SUB_MERCHANT_ID, // 子商户编号
		}})
	if err != nil {
		//这里是网络层面的错误
		log.Error(operationId, "调用第三方支付出现网络错误", err)
		resp.ErrMsg = "调用第三方支付出现网络错误"
		resp.ErrCode = 400
		return resp
	}
	if transferResult.ResultCode != ncount.ResultCodeSuccess {
		remark := "余额发送红包失败： "
		co, _ := json.Marshal(transferResult)
		err := imdb.CreateErrorLog(remark, operationId, MerOrderId, transferResult.ErrorMsg, transferResult.ErrorCode, string(co))
		if err != nil {
			log.Error(operationId, "创建错误日志失败", err)
		}
		resp.ErrCode = 400
	}
	resp.NcountOrderID = transferResult.NcountOrderId
	resp.ErrMsg = transferResult.ErrorMsg
	return resp
}

// 支付确认
func (np *NcountPay) payComfirm(OrderNo, Code string) *PayResult {
	if OrderNo == "" || Code == "" {
		return &PayResult{
			ErrCode: 400,
			ErrMsg:  "参数错误",
		}
	}
	var (
		resp = &PayResult{
			ErrCode: 0,
			ErrMsg:  "支付成功",
		}
	)

	//新生支付确认接口
	accountResp, err := np.count.QuickPayConfirm(&ncount.QuickPayConfirmReq{
		MerOrderId: ncount.GetMerOrderID(),
		QuickPayConfirmMsgCipher: ncount.QuickPayConfirmMsgCipher{
			NcountOrderId:        OrderNo,
			SmsCode:              Code,
			PaymentTerminalInfo:  "02|AA01BB",
			ReceiverTerminalInfo: "01|00001|CN|469023",
			DeviceInfo:           "192.168.0.1|E1E2E3E4E5E6|123456789012345|20000|898600MFSSYYGXXXXXXP|H1H2H3H4H5H6|AABBCC",
		},
	})
	if err != nil {
		//这里是网络层面的错误
		log.Error("调用第三方支付出现网络错误", err)
		fmt.Println(err)
		resp.ErrMsg = "调用第三方支付出现网络错误"
		resp.ErrCode = 400
		return resp
	}
	if accountResp.ResultCode == ncount.ResultCodeFail {
		remark := "确认支付失败： "
		co, _ := json.Marshal(accountResp)
		err := imdb.CreateErrorLog(remark, "", OrderNo, accountResp.ErrorMsg, accountResp.ErrorCode, string(co))
		if err != nil {
			log.Error("创建错误日志失败", err)
		}
		resp.ErrCode = 400
	}
	if accountResp.ResultCode == ncount.ResultCodeInProcess {
		resp.ErrCode = 0
		resp.ErrMsg = "该笔订单正在处理中，请耐心等待"
		return resp
	}
	resp.NcountOrderID = accountResp.NcountOrderId
	resp.ErrMsg = accountResp.ErrorMsg
	return resp
}

// 支付确认：rpc
func (cl *CloudWalletServer) PayConfirm(ctx context.Context, in *pb.PayConfirmReq) (*pb.CommonResp, error) {

	var (
		resp = &pb.CommonResp{
			ErrCode: 0,
			ErrMsg:  "支付成功",
		}
	)
	// 通过平台订单号码查询订单
	err, outTrade := imdb.GetThirdPayOrderNo(in.OrderNo)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resp.ErrCode = 400
			resp.ErrMsg = "订单不存在"
			return resp, nil
		} else {
			return nil, err
		}
	}

	if outTrade.Id == 0 {
		resp.ErrCode = 400
		resp.ErrMsg = "订单不存在"
		return resp, nil
	}

	// 梳理逻辑
	nc := NewNcountPay()
	payConfirmResp := nc.payComfirm(outTrade.NcountTureNo, in.Code)

	resp.ErrMsg = payConfirmResp.ErrMsg
	resp.ErrCode = pb.CloudWalletErrCode(payConfirmResp.ErrCode)
	return resp, nil
}

// 第三方支付回调
func (cl *CloudWalletServer) PayCallback(ctx context.Context, in *pb.PayCallbackReq) (*pb.PayCallbackResp, error) {
	var (
		resp = &pb.PayCallbackResp{
			CommonResp: &pb.CommonResp{
				ErrCode: 0,
				ErrMsg:  "回调成功",
			},
		}
	)
	// 支付回调属于业务层面的内容了
	var err error
	switch in.BusinessType {
	case pb.PayCallbackBusinessType_P_ThirdPay: // 第三方支付回调
		err = ThirdPayCallBack(in)
	}

	if err != nil {
		return nil, err
	}
	return resp, nil
}

// 第三方支付回调业务处理
func ThirdPayCallBack(in *pb.PayCallbackReq) error {
	// 查询订单，通过mer_order_id
	err, outTrade := imdb.GetThirdPayNcountMerOrderID(in.NcountOrderId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("订单不存在")
		} else {
			return err
		}
		return err
	}

	// 修改订单状态
	if in.ResultCode == ncount.ResultCodeSuccess {
		// 支付成功
		outTrade.Status = 200
		outTrade.PayTime = time.Now()
		err = imdb.UpdateThirdPayOrder(outTrade, outTrade.Id)
		if err != nil {
			return err
		}
	}

	// 修改交易记录状态
	err = imdb.FNcountTradeUpdateStatusbyThirdOrderNo(in.NcountOrderId)
	if err != nil {
		log.Error("修改订单状态失败", err, in)
		return err
	}

	log.Info("新生支付回调成功", in)

	// 传入的ID是MerID
	notifyThirdPay(in.MerOrderId)
	return nil
}
