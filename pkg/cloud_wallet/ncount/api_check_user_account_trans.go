package ncount

import (
	"github.com/pkg/errors"
)

// 交易详情记录

/*
tranMerOrderId 商 户 订 单号 1-32 商户请求交易时的商户 订单号 不可 例 如 ： 11000000111
queryType 交 易 大 类 1-32 收款：RECV 付款：PAY 退款：REFUND 转账：TRAN 不可 例如：RECV
*/
type CheckUserAccountTransMsgCipher struct {
	TranMerOrderId string `json:"tranMerOrderId" binding:"required"`
	QueryType      string `json:"queryType" binding:"required"`
}

func (c *CheckUserAccountTransMsgCipher) Valid() error {
	if c.TranMerOrderId == "" {
		return errors.New("tranMerOrderId is empty")
	}
	if c.QueryType == "" {
		return errors.New("queryType is empty")
	}
	return nil
}

type CheckUserAccountTransReq struct {
	MerOrderId                     string `json:"merOrderId" binding:"required"`
	CheckUserAccountTransMsgCipher CheckUserAccountTransMsgCipher
}

func (c *CheckUserAccountTransReq) Valid() error {
	if c.MerOrderId == "" {
		return errors.New("merOrderId is empty")
	}
	return c.CheckUserAccountTransMsgCipher.Valid()
}

/*
	version 版本号 同上送
	tranCode 交易代码 同上送
	merOrderId 商户订单号 同上送
	merId 商户 ID 同上送
	merAttach 附加数据 同上送
	charset 编码方式 同上送
	signType 签名类型 同上送
	resultCode 处理结果码 4 附录二resultCode 9999
	errorCode 异常代码 1-10 附录一 errorCode
	errorMsg 异常描述 1-200 中文、字母、数字
	tranMerOrderId 商户订单号 同上送
	ncountOrderId 新账通订单 号 新账通平台交易订单号
	orderStatus 订单状态 0：进行中 1：成功 2：失败 3：失效
	tranAmount 交易金额 格式：数字，单位：元
	feeAmount 交易手续费 格式：数字，单位：元
	bankOrderId BIS 订单号 19 格式：数字 20101409271841 23217
	realBankOrderId 银行单号 格式：数字
	signValue 签名字符串 将报文信息用
	signType 域设 置的方式签名后生成的字符 串
	businessType 业务类型 01:充值，02:转账，03：消费， 04：担保下单，08:绑定卡提 现，09：同名非绑定卡提现，
	divideAcctDtl 分账订单明 细 divideId 为分账主订单 Id divideDtlId 为 分 账 明 细 IdledgerUserId 为分账方 Id divideStatus 为分账订单状 态 6.7 分账明细状态 [{"divideId":"123456", “divideDtlList": [{ "divideDtlId":"1234 56", "ledgerUserId":"123 456", "divideStatus":”1" }] }]收款订单出现
	unconfirmedAmount 待确认金额 格式：数字，单位：元 收款订单出现
	refundableAmount 可退款金额 格式：数字，单位：元 收款订单出现
	serviceAmount 服务费金额 格式：数字，单位：元 付款订单出现
	orderFinishTm 交易完成时 间 格式:YYYYMMDDHHMMSS
	instalmentNum 分期期数 2 注：仅在下单时分期期数不为 空才返回 12
	payableFeeAmt 分期应付手 续费 0-12 注：仅在下单时分期期数不为 空才返回 2
	payableFeeType 手续费支付 方式 1 注：仅在下单时分期期数不为 空才返回 1
	firstPeriodFeeAm t 首期手续费 0-12 注：仅在下单时分期期数不为 空才返回
	eachPeriodFeeA mt 每期手续费 0-12 注：仅在下单时分期期数不为 空才返回
	firstPeriodPayAm t 首期还款金 额 0-12 注：仅在下单时分期期数不为 空才返回
	instalmentRate 商户分期实 际贴息费率 6 注：仅在下单时分期期数不为 空才返回
	instalmentAmt 商户分期实 际贴息费用 0-12 注：仅在下单时分期期数不为
*/
type CheckUserAccountTransResp struct {
	Version           string `json:"version"`
	TranCode          string `json:"tranCode"`
	MerOrderId        string `json:"merOrderId"`
	MerId             string `json:"merId"`
	MerAttach         string `json:"merAttach"`
	Charset           string `json:"charset"`
	SignType          string `json:"signType"`
	ResultCode        string `json:"resultCode"`
	ErrorCode         string `json:"errorCode"`
	ErrorMsg          string `json:"errorMsg"`
	TranMerOrderId    string `json:"tranMerOrderId"`
	NcountOrderId     string `json:"ncountOrderId"`
	OrderStatus       string `json:"orderStatus"`
	TranAmount        string `json:"tranAmount"`
	FeeAmount         string `json:"feeAmount"`
	BankOrderId       string `json:"bankOrderId"`
	RealBankOrderId   string `json:"realBankOrderId"`
	SignValue         string `json:"signValue"`
	BusinessType      string `json:"businessType"`
	DivideAcctDtl     string `json:"divideAcctDtl"`
	UnconfirmedAmount string `json:"unconfirmedAmount"`
	RefundableAmount  string `json:"refundableAmount"`
	ServiceAmount     string `json:"serviceAmount"`
	OrderFinishTm     string `json:"orderFinishTm"`
	InstalmentNum     string `json:"instalmentNum"`
	PayableFeeAmt     string `json:"payableFeeAmt"`
	PayableFeeType    string `json:"payableFeeType"`
	FirstPeriodFeeAmt string `json:"firstPeriodFeeAmt"`
	EachPeriodFeeAmt  string `json:"eachPeriodFeeAmt"`
	FirstPeriodPayAmt string `json:"firstPeriodPayAmt"`
	InstalmentRate    string `json:"instalmentRate"`
	InstalmentAmt     string `json:"instalmentAmt"`
}
