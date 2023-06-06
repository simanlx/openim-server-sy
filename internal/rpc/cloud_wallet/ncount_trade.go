package cloud_wallet

import (
	"crazy_server/pkg/cloud_wallet/ncount"
	"crazy_server/pkg/common/db"
	imdb "crazy_server/pkg/common/db/mysql_model/cloud_wallet"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

const (
	BusinessTypeBankcardRecharge   = 1 //银行卡充值
	BusinessTypeBankcardWithdrawal = 2 //银行卡提现
	BusinessTypeBankcardSendPacket = 3 //银行卡支付发送红包
	BusinessTypeBalanceSendPacket  = 4 //余额支付发送红包
	BusinessTypeReceivePacket      = 5 //领取红包
	BusinessTypePacketExpire       = 6 //红包超时退回

	BusinessTypeBankcardThirdPay = 7 //银行卡第三方支付
	BusinessTypeBalanceThirdPay  = 8 //余额第三方支付

	BusinessTypeThirdWithDraw = 9 //第三方提现到余额
)

func BusinessTypeAttr(businessType, amount, balAmount int32) (int32, int32, int32, string, error) {
	switch businessType {
	case BusinessTypeBankcardRecharge:
		return 1, 0, balAmount + amount, "银行卡充值", nil
	case BusinessTypeBankcardWithdrawal:
		return 2, 0, balAmount, "提现到银行卡", nil
	case BusinessTypeBankcardSendPacket:
		return 2, 0, balAmount, "银行卡支付发送红包", nil
	case BusinessTypeBalanceSendPacket:
		return 2, 1, balAmount, "余额支付发送红包", nil
	case BusinessTypeReceivePacket:
		return 1, 1, balAmount, "领取红包", nil
	case BusinessTypePacketExpire:
		return 1, 1, balAmount, "红包超时退回", nil
	case BusinessTypeBankcardThirdPay:
		return 2, 0, balAmount, "银行卡支付第三方", nil // 支付用户第三方内容
	case BusinessTypeBalanceThirdPay:
		return 2, 1, balAmount, "余额支付第三方", nil // 支付用户第三方内容
	case BusinessTypeThirdWithDraw:
		return 1, 1, balAmount + amount, "第三方提现到余额", nil // 第三方提现到余额
	default:
		return 0, 0, 0, "", errors.New("业务类型错误")
	}
}

// 增加账户变更日志
func AddNcountTradeLog(businessType, amount int32, userId, mainAccountId, merOrderNo, thirdOrderNo, packetID string) (err error) {
	if len(mainAccountId) < 1 {
		//获取用户账户信息
		accountInfo, err := imdb.GetNcountAccountByUserId(userId)
		if err != nil || accountInfo.Id <= 0 {
			return errors.New(fmt.Sprintf("查询账户数据失败 %s,error:%s", userId, err.Error()))
		}
		mainAccountId = accountInfo.MainAccountId
	}

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
		MerOrderId:   merOrderNo,
	})
}
