package ncount

import (
	"github.com/pkg/errors"
)

/*
cardNo 支 付 银 行 卡 卡号 1-30 格式：数字 不可
holderName 持卡人姓名 0-40 持卡人姓名 不可
cardAvailableDate 信 用 卡 有 效 期 4 贷记卡有效期 信用卡 不可空 例如：0320 含义：2020 年 03 月
cvv2 信用卡 CVV2 3 格式:数字 信用卡 不可空 例如：318
mobileNo 银 行 签 约 手 机号 11 格式：数字 不可
identityType 证件类型 2 暂仅支持 1 : 身份证 不可
identityCode 证件号码 0-50 格式：数字、英 文字母 不可
userId 用户 ID 1-32 格式：数字，字 母，下划线，竖 划线，中划线 不可 例如：102121
merUserIp 商户用户 IP 0-128 商户 用户签 约 时所 在的机 器 IP 地址 可空 例 如 ： 211.12.38.88
*/

type BindCardMsgCipherText struct {
	CardNo            string `json:"cardNo" binding:"required"`
	HolderName        string `json:"holderName" binding:"required"`
	CardAvailableDate string `json:"cardAvailableDate" binding:"required"`
	Cvv2              string `json:"cvv2" binding:"required"`
	MobileNo          string `json:"mobileNo" binding:"required"`
	IdentityType      string `json:"identityType" binding:"required"`
	IdentityCode      string `json:"identityCode" binding:"required"`
	UserId            string `json:"userId" binding:"required"`
}

func (b *BindCardMsgCipherText) Valid() error {
	if b.CardNo == "" {
		return errors.New("cardNo is empty")
	}
	if b.HolderName == "" {
		return errors.New("holderName is empty")
	}
	if b.MobileNo == "" {
		return errors.New("mobileNo is empty")
	}
	if b.IdentityType == "" {
		return errors.New("identityType is empty")
	}
	if b.IdentityCode == "" {
		return errors.New("identityCode is empty")
	}
	if b.UserId == "" {
		return errors.New("userId is empty")
	}
	return nil
}

type BindCardReq struct {
	MerOrderId            string `json:"merOrderId" binding:"required"`
	BindCardMsgCipherText BindCardMsgCipherText
}

func (b *BindCardReq) Valid() error {
	if b.MerOrderId == "" {
		return errors.New("merOrderId is empty")
	}
	return b.BindCardMsgCipherText.Valid()
}

/*
version 版本号 同上送
tranCode 交易代码 同上送
merOrderId 商户订单号 同上送
merId 商户 ID 同上送
merAttach 附加数据 同上送
charset 编码方式 同上送
signType 签名类型 同上送
resultCode 处理结果码 4 格式：数字 详情参见 6.2 附录二
resultCode 说明 0000
errorCode 异常代码 100 格式：数字 详情参见 6.1 附录一
errorCode 说明
errorMsg 异常描述 中文、字母、数字
ncountOrderId 签约订单号 32 签约下单成功后新账通平台生 成的订单号（签约后续流程使 用） 2012844332232
signValue 签名字符串 将报文信息用 signType 域设 置的方式签名后生成的字符串
*/
type BindCardResp struct {
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
	NcountOrderId string `json:"ncountOrderId"`
	SignValue     string `json:"signValue"`
}
