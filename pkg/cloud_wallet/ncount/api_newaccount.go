package ncount

import (
	"errors"
)

/*
	4.1.1.1 商户请求参数列表（POST）
	version 版本号 5 目前必须为 1.0 不可 例如：1.0
	tranCode 交易代码 5 此交易只能为 R010 不可 例如：R010
	merId 商户 ID 12 新账通平台提供给商户 的唯一 ID 不可 例 如 ： 00000100000
	merOrderId 商 户 订 单 号 1 - 32 格式：数字，字母，下 划线，竖划线，中划线 不可 例 如 ： aa201612011 102
	submitTime 请 求 提 交 时间 14 格 式 ： YYYYMMDDHHMMS S 不可 例 如 ： 20161201110 233
	msgCiphert ext 报文密文 1-4000 用平台公钥 RSA 加密后 base64 的编码值 不可
	signType 签名类型 1 1：RSA 不可 例如：1
	signValue 签 名 密 文 串 将 报 文 信 息 用signType 域设 置的方 式签名后生成的字符串 不可
	merAttach 附加数据 0-80 格式：英文字母/汉字 可空
	charset 编码方式 1 1：UTF8 不可 例如：1
*/
type NewAccountReq struct {
	OrderID       string `json:"orderId" binding:"required"`
	MsgCipherText *NewAccountMsgCipherText
}

func (n *NewAccountReq) Vaild() error {
	if n.OrderID == "" {
		return errors.New("orderId is empty")
	}
	if n.MsgCipherText == nil {
		return errors.New("msgCipherText is empty")
	}
	return n.MsgCipherText.Vaild()
}

/*
	merUserId 商 户 用户 唯 一标识 1-30 按照商户侧规则 同一个平台商户唯一 字符串类型 不可 例 如 ： xsqianyi1_ 148080610 02
	mobile 用户手机号 11 用户本人在运营商已 实名的手机号 不可 例如： 138000000 00
	userName 真实姓名 1-30 和身份证上的姓名保 持一致 不可 例如：张三
	certNo 身份证号 18 身份证号码 不可 例如： 110210199 008123456
*/
type NewAccountMsgCipherText struct {
	MerUserId string `json:"merUserId" binding:"required"`
	Mobile    string `json:"mobile" binding:"required"`
	UserName  string `json:"userName" binding:"required"`
	CertNo    string `json:"certNo" binding:"required"`
}

func (n *NewAccountMsgCipherText) Vaild() error {
	if n.MerUserId == "" {
		return errors.New("merUserId is empty")
	}
	if n.Mobile == "" {
		return errors.New("mobile is empty")
	}
	if n.UserName == "" {
		return errors.New("userName is empty")
	}
	if n.CertNo == "" {
		return errors.New("certNo is empty")
	}
	return nil
}

/*
tranCode 交易代码 同上送 merOrderId 商户订单号 同上送 merId 商户 ID 同上送 merAttach 附加数据 同上送 charset 编码方式 同上送 signType 签名类型 同上送
*/

// 调用新生接口创建账户，拿到的返回结果
type NewAccountResp struct {
	Version    string `json:"version" binding:"required"`
	TranCode   string `json:"tranCode" binding:"required"`
	MerOrderId string `json:"merOrderId" binding:"required"`
	MerId      string `json:"merId" binding:"required"`
	MerAttach  string `json:"merAttach" binding:"required"`
	Charset    string `json:"charset" binding:"required"`
	SignType   string `json:"signtype" binding:"required"`

	ResultCode string `json:"resultCode" binding:"required"`
	ErrorCode  string `json:"errorCode" binding:"required"`
	ErrorMsg   string `json:"errorMsg" binding:"required"`
	UserId     string `json:"userId" binding:"required"`
	SignValue  string `json:"signValue" binding:"required"`
}
