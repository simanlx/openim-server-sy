package ncount

import (
	"github.com/pkg/errors"
)

// 用户账户详情

/*
userId 用户编号 12 用户编号 不可 例 如 ： 11000000111
acctType 账户类型 2 1 现金户 3 待清算账户 不可 例 如 ： 11000000112
startDate 开始时间 8 格式：YYYYMMDD 不可 例如：20181230
endDate 结束时间 8 格式：YYYYMMDD 不可 例如：20181231
pageNum 请求页数 4 默认 1 可空 例如：1
*/

type CheckUserAccountDetailMsgCipher struct {
	UserId    string `json:"userId" binding:"required"`
	AcctType  string `json:"acctType" binding:"required"`
	StartDate string `json:"startDate" binding:"required"`
	EndDate   string `json:"endDate" binding:"required"`
	PageNum   string `json:"pageNum" binding:"required"`
}

func (c *CheckUserAccountDetailMsgCipher) Valid() error {
	if c.UserId == "" {
		return errors.New("userId is empty")
	}
	if c.AcctType == "" {
		return errors.New("acctType is empty")
	}
	if c.StartDate == "" {
		return errors.New("startDate is empty")
	}
	if c.EndDate == "" {
		return errors.New("endDate is empty")
	}
	return nil
}

type CheckUserAccountDetailReq struct {
	MerOrderId                      string `json:"merOrderId" binding:"required"`
	CheckUserAccountDetailMsgCipher CheckUserAccountDetailMsgCipher
}

func (c *CheckUserAccountDetailReq) Valid() error {
	if c.MerOrderId == "" {
		return errors.New("merOrderId is empty")
	}
	return c.CheckUserAccountDetailMsgCipher.Valid()
}

/*
version 版本号 同上送
tranCode 交易代码 同上送
merOrderId 商户订单号 同上送
merId 商户 ID 同上送
merAttach 附加数据 同上送
charset 编码方式 同上送新账通平台商户接入规范 2021/11/11
signType 签名类型 同上送
resultCode 处理结果码 4 详情参见附录二 resultCode 9999
errorCode 异常代码 1-10 详情参见附录一 errorCode
errorMsg 异常描述 1-200 中文、字母、数字
userId 用户编号 12 acctType 账户类型 1
count 总条数
transDetailList 账户明细

// list
transType 交易类型 转帐、提现等 transAmt 交易金额
ncountOrderId 新 账 通 订 单 号 transTime 交易日期 YYYYMMDD
accountingTime 交易时间 YYYYMMDDHHMMSS
balance 交易后余额
remark 预留字段
summary 账户摘要 描述交易类型和手续费扣减信 息
businessType 业务类型 业务类型
ieType 借贷标识 I 代表+，E

// list
signValue 签名字符串 将报文信息用 signType 域设 置的方式签名后生成的字符串
*/
type CheckUserAccountDetailResp struct {
	Version         string        `json:"version" binding:"required"`
	TranCode        string        `json:"tranCode" binding:"required"`
	MerOrderId      string        `json:"merOrderId" binding:"required"`
	MerId           string        `json:"merId" binding:"required"`
	MerAttach       string        `json:"merAttach" binding:"required"`
	Charset         string        `json:"charset" binding:"required"`
	SignType        string        `json:"signType" binding:"required"`
	ResultCode      string        `json:"resultCode" binding:"required"`
	ErrorCode       string        `json:"errorCode" binding:"required"`
	ErrorMsg        string        `json:"errorMsg" binding:"required"`
	UserId          string        `json:"userId" binding:"required"`
	AcctType        string        `json:"acctType" binding:"required"`
	Count           int32         `json:"count" binding:"required"`
	SignValue       string        `json:"signValue" binding:"required"`
	TransDetailList []interface{} `json:"transDetailList" binding:"required"`
}

/*
	// list
	transType 交易类型 转帐、提现等 transAmt 交易金额
	ncountOrderId 新 账 通 订 单 号 transTime 交易日期 YYYYMMDD
	accountingTime 交易时间 YYYYMMDDHHMMSS
	balance 交易后余额
	remark 预留字段
	summary 账户摘要 描述交易类型和手续费扣减信 息
	businessType 业务类型 业务类型
	ieType 借贷标识 I 代表+，E
*/

type TransDetailList struct {
	TransType      string `json:"transType" binding:"required"`
	NcountOrderId  string `json:"ncountOrderId" binding:"required"`
	AccountingTime string `json:"accountingTime" binding:"required"`
	Balance        string `json:"balance" binding:"required"`
	Remark         string `json:"remark" binding:"required"`
	Summary        string `json:"summary" binding:"required"`
	BusinessType   string `json:"businessType" binding:"required"`
	IeType         string `json:"ieType" binding:"required"`
}
