package ncount

import (
	"github.com/pkg/errors"
)

/*
	payUserId 付 款方 用 户编号 12 用户编号 不可 例如：11000000111
	receiveUserId 收款方 ID 12 用户编号 不可 例如：11000000112
	tranAmount 转账金额 1-10 单位：元，保留小数点 2 位 不可 例如：1.11
	businessType 业务类型 2 02 转账，默认 02 可空 例如：02
*/

/*
payUserId 付 款方 用 户编号 12 用户编号 不可 例如：11000000111
receiveUserId 收款方 ID 12 用户编号 不可 例如：11000000112
tranAmount 转账金额 1-10 单位：元，保留小数点 2 位 不可 例如：1.11
businessType 业务类型 2 02 转账，默认 02 可空 例如：02
*/
type TransferMsgCipher struct {
	PayUserId     string `json:"payUserId" binding:"required"`
	ReceiveUserId string `json:"receiveUserId" binding:"required"`
	TranAmount    string `json:"tranAmount" binding:"required"`
	BusinessType  string `json:"businessType" binding:"required"`
}

func (t *TransferMsgCipher) Valid() error {
	if t.PayUserId == "" {
		return errors.New("payUserId is empty")
	}
	if t.ReceiveUserId == "" {
		return errors.New("receiveUserId is empty")
	}
	if t.TranAmount == "" {
		return errors.New("tranAmount is empty")
	}
	return nil
}

type TransferReq struct {
	MerOrderId        string `json:"merOrderId" binding:"required"`
	TransferMsgCipher TransferMsgCipher
}

func (t *TransferReq) Valid() error {
	if t.MerOrderId == "" {
		return errors.New("merOrderId is empty")
	}
	return t.TransferMsgCipher.Valid()
}

/*
ncountOrderId 订单号 32 新账通平台订单号
orderDate 平 台 订 单 日 期 8 YYYYMMDD
resultCode 处理结果码 4 附录二 resultCode 9999
errorCode 异常代码 1-10 附录一 errorCode
errorMsg 异常描述 1-200 中文、字母、数字
payAcctAmount 付 款 方 账 户 余额 1-10 格式：整数 单位：元 交易成功时返回
recvAcctAmount 收 款 方 账 户 余额 1-10 格式：整数 单位：元 交易成功时返回
businessType 业务类型 2 02 转账 02
signValue 签名字符串 将报文信息用signType 域设 置的方式签名后生成的字符串
*/
type TransferResp struct {
	BaseReturnParam
	NcountOrderId  string `json:"ncountOrderId" binding:"required"`
	OrderDate      string `json:"orderDate" binding:"required"`
	ResultCode     string `json:"resultCode" binding:"required"`
	ErrorCode      string `json:"errorCode" binding:"required"`
	ErrorMsg       string `json:"errorMsg" binding:"required"`
	PayAcctAmount  string `json:"payAcctAmount" binding:"required"`
	RecvAcctAmount string `json:"recvAcctAmount" binding:"required"`
	BusinessType   string `json:"businessType" binding:"required"`
	SignValue      string `json:"signValue" binding:"required"`
}
