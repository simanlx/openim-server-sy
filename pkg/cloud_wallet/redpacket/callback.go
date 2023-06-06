package redpacket

import (
	"crazy_server/pkg/cloud_wallet/ncount"
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/db"
	commonDB "crazy_server/pkg/common/db"
	imdb "crazy_server/pkg/common/db/mysql_model/cloud_wallet"
	"crazy_server/pkg/common/log"
	"crazy_server/pkg/contrive_msg"
	"crazy_server/pkg/tools/redpacket"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

// 这里是生成红包的算法 + 回调

// 当用户发布红包发送成功的时候，调用这个回调函数进行发布红包的后续处理
func HandleSendPacketResult(redPacketID, OperateID string) error {
	//1. 查询红包信息
	redpacketInfo, err := imdb.GetRedPacketInfo(redPacketID)
	if err != nil {
		log.Error(OperateID, "get red packet info error", zap.Error(err))
		return err
	}
	if redpacketInfo == nil {
		log.Error(OperateID, "red packet info is nil")
		return errors.New("red packet info is nil")
	}

	// 2. 生成红包
	if redpacketInfo.PacketType == 2 && redpacketInfo.Number > 1 {
		// 群红包
		err = GroupPacket(redpacketInfo, redPacketID)
		if err != nil {
			return err
		}
	}

	// 3. 修改红包状态
	err = imdb.UpdateRedPacketStatus(redPacketID, imdb.RedPacketStatusNormal)
	if err != nil {
		log.Error(OperateID, "update red packet status error", zap.Error(err))
		return err
	}

	// todo 发送红包消息
	freq := &contrive_msg.FPacket{
		PacketID:        redPacketID,
		UserID:          redpacketInfo.UserID,
		PacketType:      redpacketInfo.PacketType,
		IsLucky:         redpacketInfo.IsLucky,
		ExclusiveUserID: redpacketInfo.ExclusiveUserID,
		PacketTitle:     redpacketInfo.PacketTitle,
		Amount:          redpacketInfo.Amount,
		Number:          redpacketInfo.Number,
		ExpireTime:      redpacketInfo.ExpireTime,
		MerOrderID:      redpacketInfo.MerOrderID,
		OperateID:       redpacketInfo.OperateID,
		RecvID:          redpacketInfo.RecvID,
		CreatedTime:     redpacketInfo.CreatedTime,
		UpdatedTime:     redpacketInfo.UpdatedTime,
		IsExclusive:     redpacketInfo.IsExclusive,
	}
	return contrive_msg.SendSendRedPacket(freq, int(redpacketInfo.PacketType))
}

// 给群发的红包
func GroupPacket(req *db.FPacket, redpacketID string) error {

	// 在群红包这里，
	var err error
	if req.IsLucky == 1 {
		// 如果说是手气红包，分散放入红包池
		err = spareRedPacket(redpacketID, int(req.Amount), int(req.Number))
	} else {
		// 凭手气红包
		err = spareEqualRedPacket(redpacketID, int(req.Amount), int(req.Number))
	}
	if err != nil {
		log.Error(req.OperateID, zap.Error(err))
		return err
	}
	return err
}

// 将红包放入红包池
func spareRedPacket(packetID string, amount, number int) error {
	defer func() {
		if err := recover(); err != nil {
			log.Error("spareRedPacket panic", zap.Any("err", err))
		}
	}()
	// 将发送的红包进行计算
	result := redpacket.GetRedPacket(amount, number)
	err := commonDB.DB.SetRedPacket(packetID, result)
	if err != nil {
		return err
	}
	return nil
}

// amount = 3 ,number =3
func spareEqualRedPacket(packetID string, amount, number int) error {
	result := []int{}
	for i := 0; i < number; i++ {
		result = append(result, amount)
	}
	// 将发送的红包进行计算
	err := commonDB.DB.SetRedPacket(packetID, result)
	if err != nil {
		return err
	}
	return nil
}

const (
	BusinessTypeBankcardRecharge   = 1 //银行卡充值
	BusinessTypeBankcardWithdrawal = 2 //银行卡提现
	BusinessTypeBankcardSendPacket = 3 //银行卡支付发送红包
	BusinessTypeBalanceSendPacket  = 4 //余额支付发送红包
	BusinessTypeReceivePacket      = 5 //领取红包
	BusinessTypePacketExpire       = 6 //红包超时退回
)

func BusinessTypeAttr(businessType, amount, balAmount int32) (int32, int32, int32, string, error) {
	switch businessType {
	case BusinessTypeBankcardRecharge:
		return 1, 0, balAmount + amount, "银行卡充值", nil
	case BusinessTypeBankcardWithdrawal:
		return 2, 0, balAmount - amount, "提现到银行卡", nil
	case BusinessTypeBankcardSendPacket:
		return 2, 0, balAmount, "银行卡支付发送红包", nil
	case BusinessTypeBalanceSendPacket:
		return 2, 1, balAmount, "余额支付发送红包", nil
	case BusinessTypeReceivePacket:
		return 1, 1, balAmount + amount, "领取红包", nil
	case BusinessTypePacketExpire:
		return 1, 1, balAmount + amount, "红包超时退回", nil
	default:
		return 0, 0, 0, "", errors.New("业务类型错误")
	}
}

// 银行卡充值到红包账户
func BankCardRechargePacketAccount(userId, bindCardAgrNo string, amount int32, packetID string) error {
	//获取用户账户信息
	accountInfo, err := imdb.GetNcountAccountByUserId(userId)
	if err != nil || accountInfo.Id <= 0 {
		return errors.New("账户信息不存在")
	}

	//充值支付
	accountResp, err := ncount.NewCounter().QuickPayOrder(&ncount.QuickPayOrderReq{
		MerOrderId: ncount.GetMerOrderID(),
		QuickPayMsgCipher: ncount.QuickPayMsgCipher{
			PayType:       "3", //绑卡协议号充值
			TranAmount:    cast.ToString(cast.ToFloat64(amount) / 100),
			NotifyUrl:     config.Config.Ncount.Notify.RechargeNotifyUrl,
			BindCardAgrNo: bindCardAgrNo,
			ReceiveUserId: accountInfo.PacketAccountId, //收款账户
			UserId:        accountInfo.MainAccountId,
			SubMerchantId: ncount.SUB_MERCHANT_ID, // 子商户编号
		}})
	if err != nil {
		return errors.New(fmt.Sprintf("充值失败(%s)", err.Error()))
	} else {
		if accountResp.ResultCode != "0000" {
			return errors.New(fmt.Sprintf("充值失败 (%s,%s)", accountResp.ErrorCode, accountResp.ErrorMsg))
		}
	}

	//增加账户变更日志
	err = AddNcountTradeLog(BusinessTypeBankcardSendPacket, amount, userId, accountInfo.MainAccountId, accountResp.NcountOrderId, packetID)
	if err != nil {
		return errors.New(fmt.Sprintf("增加账户变更日志失败(%s)", err.Error()))
	}

	return nil
}

// 增加账户变更日志
func AddNcountTradeLog(businessType, amount int32, userId, mainAccountId, thirdOrderNo, packetID string) (err error) {
	//获取用户余额
	accountResp, err := ncount.NewCounter().CheckUserAccountInfo(&ncount.CheckUserAccountReq{
		OrderID: ncount.GetMerOrderID(),
		UserID:  mainAccountId,
	})
	fmt.Println("accountResp Println", accountResp, err)
	if err != nil {
		return errors.New(fmt.Sprintf("查询账户信息失败(%s)", err.Error()))
	} else {
		if accountResp.ResultCode != "0000" {
			return errors.New(fmt.Sprintf("查询账户信息失败 (%s,%s)", accountResp.ErrorCode, accountResp.ErrorMsg))
		}
	}

	//余额变更、注意精度
	decimalBalAmount, _ := decimal.NewFromString(accountResp.BalAmount)
	balAmount := cast.ToInt32(decimalBalAmount.Mul(decimal.NewFromInt(100)).IntPart()) //用户余额
	changeType, ncountStatus, afterAmount, describe, err := BusinessTypeAttr(businessType, amount, balAmount)
	if err != nil {
		return err
	}

	//数据入库
	return imdb.FNcountTradeCreateData(&db.FNcountTrade{
		UserID:       userId,
		Type:         changeType,
		BusinessType: businessType,
		Describe:     describe,
		Amount:       amount,
		AfterAmount:  afterAmount,
		ThirdOrderNo: thirdOrderNo,
		NcountStatus: ncountStatus,
		PacketID:     packetID,
	})
}
