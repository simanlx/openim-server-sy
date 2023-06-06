package ncount

import (
	"github.com/pkg/errors"
)

// 快捷支付
/*
	tranAmount 支付金额 1-12 格式：数字（以元为单 位） 不可 例如： 100
	payType 支付方式 1 2:银行卡卡号 3:绑卡协议号 不可 例如：2
	cardNo 支付银行 卡卡号 0-30 payType=2 不可空 可空 例 如 ： 611888 812128
	holderName 持卡人姓 名 0-40payType=2 不可空 可空
	cardAvailableDate 信用卡有 效期 0-4 payType=2，且为 信用卡时不可空 可空 例 如 ： 0320 含 义 ： 2020 年 03 月 cvv2 信 用 卡
	CVV2 0-3 payType=2，且为 信用卡时不可空 可空 例 如 ： 318
	mobileNo 银行签约 手机号 011 payType=2 不可空 可空
	identityType 证件类型 0-2
	payType=2 不可空 暂仅支持 1 身份证 可空
	identityCode 证件号码 0-50
	payType=2 不可空 可空
	bindCardAgrNo 绑卡协议 号 30
	payType=3 不可空 可空
	notifyUrl 商户异步 通知地址 1-255 后台通知地址 不可 例 如 ： https:/ /www.x
*/
type QuickPayMsgCipher struct {
	TranAmount    string `json:"tranAmount" binding:"required"`    // 支付金额
	PayType       string `json:"payType" binding:"required"`       // 支付方式
	NotifyUrl     string `json:"notifyUrl" binding:"required"`     // 商户异步通知地址
	BindCardAgrNo string `json:"bindCardAgrNo" binding:"required"` // 绑卡协议号
	UserId        string `json:"userId" binding:"required"`        // 用户编号
	ReceiveUserId string `json:"receiveUserId" binding:"required"` // 收款方ID
	SubMerchantId string `json:"subMerchantId" binding:"required"` // 商户渠道进件ID
}

func (q *QuickPayMsgCipher) Valid() error {
	if q.TranAmount == "" {
		return errors.New("支付金额不能为空")
	}
	if q.PayType == "" {
		q.PayType = "3" //2是银行卡支付，3.协议号支付
	}
	if q.PayType == "3" {
		if q.BindCardAgrNo == "" {
			return errors.New("绑卡协议号不能为空")
		}
	}
	if q.NotifyUrl == "" {
		return errors.New("异步通知地址不能为空")
	}
	if q.SubMerchantId == "" {
		return errors.New("商户渠道进件ID不能为空")
	}
	if q.ReceiveUserId == "" {
		return errors.New("收款方ID不能为空")
	}
	return nil
}

type QuickPayOrderReq struct {
	MerOrderId        string `json:"merOrderId" binding:"required"`
	QuickPayMsgCipher QuickPayMsgCipher
}

func (q *QuickPayOrderReq) Valid() error {
	if q.MerOrderId == "" {
		return errors.New("商户订单号不能为空")
	}
	return q.QuickPayMsgCipher.Valid()
}

/*
resultCode 处理结果码 4 详情参见附录二 resultCode 9999
errorCode 异常代码 1-10 详情参见附录一 errorCode
errorMsg 异常描述 1-200 中文、字母、数字
ncountOrderId 新账通订单 号 32 新账通平台交易订单号
submitTime 商户请求时 间 同上送
signValue 签名字符串 将报文信息用
signType 域设 置的方式签名后生成的字符 串
*/
type QuickPayOrderResp struct {
	BaseReturnParam
	ResultCode    string `json:"resultCode"`
	ErrorCode     string `json:"errorCode"`
	ErrorMsg      string `json:"errorMsg"`
	NcountOrderId string `json:"ncountOrderId"`
	SubmitTime    string `json:"submitTime"`
	SignValue     string `json:"signValue"`
}
