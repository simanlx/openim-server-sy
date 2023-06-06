package ncount

import (
	"github.com/pkg/errors"
)

/*
ncountOrder Id 签约订单号 1-32 签约请求下单接口新账 通响应的签约订单号 不可 例 如 ： 20170323
新账通平台商户接入规范 2021/11/11 Copyright 新生支付有限公司 - 31 - 12244422 smsCode 签约短信验 证码 1-8 格式：字母，数字 不可 例 如 ： 788900
merUserIp 商户用户 IP 0-128 商户用户签约时所在的 机器 IP 地址 可空 例 如 ： 211.12.38. 88
*/

type BindCardConfirmMsgCipherText struct {
	NcountOrderId string `json:"ncountOrderId" binding:"required"`
	SmsCode       string `json:"smsCode" binding:"required"`
	MerUserIp     string `json:"merUserIp" binding:"required"`
}

func (b *BindCardConfirmMsgCipherText) Valid() error {
	if b.NcountOrderId == "" {
		return errors.New("ncountOrderId is empty")
	}
	if b.SmsCode == "" {
		return errors.New("smsCode is empty")
	}
	if b.MerUserIp == "" {
		return errors.New("merUserIp is empty")
	}
	return nil
}

type BindCardConfirmReq struct {
	MerOrderId                   string `json:"merOrderId" binding:"required"`
	BindCardConfirmMsgCipherText BindCardConfirmMsgCipherText
}

func (b *BindCardConfirmReq) Valid() error {
	if b.MerOrderId == "" {
		return errors.New("merOrderId is empty")
	}
	return b.BindCardConfirmMsgCipherText.Valid()
}

/*

version 版本号 同上送
tranCode 交易代码 同上送
merOrderId 商户订单号 同上送
merId 商户 ID 同上送
merAttach 附加数据 同上送
charset 编码方式 同上送
signType 签名类型 同上送
resultCode 处理结果码 4 详情参见 6.2 附录二resultCode 0000
errorCode 异常代码 100 详情参见 6.1 附录一 errorCode
errorMsg 异常描述 中文、字母、数字
bindCardAgrNo 绑卡协议号 30 新账通返回的用户签约协议 号
bankCode 签约银行简码 8 附录 6.3 附录三 银行简码
cardType 支付银行卡卡 类型 1 1:借记卡(DEBITCARD) 2:信用卡(CREDITCARD) 例如：1
shortCardNo 签约银行卡后 四位 4 支付签约的银行卡后四位
signValue 签名字符串 将报文信息用
*/

type BindCardConfirmResp struct {
	Version       string `json:"version"`
	TranCode      string `json:"tranCode"`
	MerOrderId    string `json:"merOrderId"`
	MerId         string `json:"merId"`
	MerAttach     string `json:"merAttach"`
	Charset       string `json:"charset"`
	SignType      string `json:"signType"`
	ResultCode    string `json:"resultCode"`
	ErrorCode     string `json:"errorCode"`
	ErrorMsg      string `json:"errorMsg"`
	BindCardAgrNo string `json:"bindCardAgrNo"`
	BankCode      string `json:"bankCode"`
	CardType      string `json:"cardType"`
	ShortCardNo   string `json:"shortCardNo"`
	SignValue     string `json:"signValue"`
}
