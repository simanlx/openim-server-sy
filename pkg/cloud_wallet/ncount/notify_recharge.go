package ncount

import "crazy_server/pkg/common/log"

// 充值回调

/*
resultCode 交易结果 4 详情参见 6.2 附录二 resultCode 0000
errorCode 异常代码 10 详情参见 6.1 附录一 errorCode
errorMsg 异常描述 200 中文、字母、数字
ncountOrderId 新账通订单 32 商户支付流程订单凭证

tranAmount 商户订单金 额 1-10 格式：数字 单位：元
checkDate 对账日期 14 格式： YYYYMMDD
submitTime 提交时间 同上送
tranFinishTime 交易完成时 间 14 格式： YYYYMMDDHHMMSS 本域为订单的完成时间 20111007094626
bankCode 签约银行简 码 8 附录三 银行简码
cardType 支付银行卡 卡类型 1 1:借记卡(DEBITCARD) 2:信用卡(CREDITCARD) 1
shortCardNo 支付银行卡 后四位 4 9876
bindCardAgrNo 绑卡协议号 30 商户绑卡确认后获得的 绑卡协议号
businessType 业务类型 2 03 消费，04 担保下单
feeAmount 交易手续费 格式：数字，单位：元
divideAcctDtl 分账订单明 细 divideId 为分账主订 单 Id divideDtlId 为分账明 细 Id ledgerUserId 为分账 方 Id divideStatus 为 分账 订 单状态 6.7 分账明细状态 [{"divideId":"123456", “divideDtlList": [{ "divideDtlId":"123456", "ledgerUserId":"123456", "divideStatus":”1" }] }]
signValue 签名字符串 将报文信息用 signType 域设置的方式加密后生 成的字符串
recvAcctAmount 收款方账户 余额 1-10 格式：数字 单位：元 交易成功时返回
instalmentNum 分期期数 2 注：仅在下单时分期期数 不为空才返回 例如：12
payableFeeAmt 分期应付手 续费 0-12 注：仅在下单时分期期数 不为空才返回 例如：2
payableFeeType 手续费支付 方式 1 注：仅在下单时分期期数 不为空才返回 例如：1
firstPeriodFeeAmt 首期手续费 0-12 注：仅在下单时分期期数 不为空才返回 例如：5
eachPeriodFeeAmt 每期手续费 0-12 注：仅在下单时分期期数 例如：5
firstPeriodPayAmt 首期还款金 额 0-12 注：仅在下单时分期期数 不为空才返回 例如：5
instalmentRate 商户分期实 际贴息费率 6 注：仅在下单时分期期数 不为空才返回 例如：001212
instalmentAmt 商户分期实 际贴息费用 0-12 注：仅在下单时分期期数 不为空才返回 例如：1.23
*/

// 这里是回调请求
type NotifyRechargeReq struct {
	ResultCode    string `json:"resultCode"`
	ErrorCode     string `json:"errorCode"`
	ErrorMsg      string `json:"errorMsg"`
	NcountOrderID string `json:"ncountOrderId"`

	TranAmount        string `json:"tranAmount"`
	CheckDate         string `json:"checkDate"`
	SubmitTime        string `json:"submitTime"`
	TranFinishTime    string `json:"tranFinishTime"`
	BankCode          string `json:"bankCode"`
	CardType          string `json:"cardType"`
	ShortCardNo       string `json:"shortCardNo"`
	BindCardAgrNo     string `json:"bindCardAgrNo"`
	BusinessType      string `json:"businessType"`
	FeeAmount         string `json:"feeAmount"`
	DivideAcctDtl     string `json:"divideAcctDtl"`
	SignValue         string `json:"signValue"`
	RecvAcctAmount    string `json:"recvAcctAmount"`
	InstalmentNum     string `json:"instalmentNum"`
	PayableFeeAmt     string `json:"payableFeeAmt"`
	PayableFeeType    string `json:"payableFeeType"`
	FirstPeriodFeeAmt string `json:"firstPeriodFeeAmt"`
	EachPeriodFeeAmt  string `json:"eachPeriodFeeAmt"`
	FirstPeriodPayAmt string `json:"firstPeriodPayAmt"`
	InstalmentRate    string `json:"instalmentRate"`
}

// 如果说请求回调走到这里来了

func (n *NotifyRechargeReq) VarifyNotify(req *NotifyRechargeReq) {
	if ResultCodeFail == req.ResultCode {
		// 记录数据库 ： 充值失败
		log.Error("回调通知", req)
		// 记录数据库充值失败
	}
	if ResultCodeSuccess == req.ResultCode {
		// 记录数据库 ： 充值成功
		log.Error("回调通知", req)
		// 记录数据库充值成功
	}
}
