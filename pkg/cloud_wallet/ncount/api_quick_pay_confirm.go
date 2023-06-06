package ncount

import (
	"github.com/pkg/errors"
)

// 支付确认接口

/*
ncountOrderId 新账通订 单号 32 支付下单请求接口新账通 响应的支付订单号 不可 例 如 ： 20121214324
smsCode 短信验证 码 1-8 格式：字母，数字 可空 例如：7889
merUserIp 商户用户 IP 0-128 商户用户支付时所在的机 器 IP 地址 可空 例 如 ： 211.12.38.88
paymentTermin alInfo 付款方终 端信息 0-100 付款方交易终端设备类型、 编码等 注：包含两个子域（终端设 备类型，终端设备编码）， 子域值间用“|”隔开，格式 如下：子域值 1|子域值 2 详情参见 6.9.1 付款方终端及 设备信息说明 不可 例 如 ： 01|10001 receiverTermina lInfo 收款方终 端信息 0-100 收款方交易终端设备的编 码、类型、国家、地区等 注：包含四个子域（终端类 型、终端编码、国家编码、 地区编码），子域值间用“|” 隔开，格式如下：子域值 1|子域值 2|子域值 3|子域 不可 例 如 ： 01|00001|CN |110000
*/
type QuickPayConfirmMsgCipher struct {
	NcountOrderId        string `json:"ncountOrderId" binding:"required"`
	SmsCode              string `json:"smsCode" `
	PaymentTerminalInfo  string `json:"paymentTerminalInfo" binding:"required"`
	ReceiverTerminalInfo string `json:"receiverTerminalInfo" binding:"required"`
	/*
		deviceInfo 设备信息 0-200 绑卡设备信息，如：IP、 MAC、IMEI、IMSI、ICCID、 Wif 如 ： 192.168.0.1|E 1E2E3E4E5E6| 12345678901 2345|20000|1 23456789012 34567890|H1 H2H3H4H5H 6|A.BCDEFG,- H.IJKLMN
		businessType 业务类型 2 03 消费，04 担保下单。默 认 03 可空 例如：03
		feeType 手续费内 扣外扣 1 0 外扣，1 内扣，默认外扣， 本期业务类型 04 担保下单 仅支持内扣 可空 例如：1
		divideAcctDtl 分账明细 用户 ID，金额列表。可空。 若业务类型 04 时，必填。 用户 id，不大于 12 位， 分账金额算上小数点不大 于 10 位，金额小数点后 2 位，单位元 可空 例如： [{\"ledgerUserId \": \"22000000139 0\",\"amount\": \"50\"}, {\"ledgerUserId\ ": \"22000000140 8\",\"amount\": \"50\"}]
		feeAmountUser 手续费承 担方 id 12
	*/
	DeviceInfo string `json:"deviceInfo" binding:"required"`
	//BusinessType  string `json:"businessType" `
	//FeeType       string `json:"feeType" `
	//DivideAcctDtl string `json:"divideAcctDtl" `
	//FeeAmountUser string `json:"feeAmountUser" `
}

func (q *QuickPayConfirmMsgCipher) Valid() error {
	if q.NcountOrderId == "" {
		return errors.New("ncountOrderId is required")
	}
	if q.PaymentTerminalInfo == "" {
		return errors.New("paymentTerminalInfo is required")
	}
	if q.ReceiverTerminalInfo == "" {
		return errors.New("receiverTerminalInfo is required")
	}
	if q.DeviceInfo == "" {
		return errors.New("deviceInfo is required")
	}
	return nil
}

type QuickPayConfirmReq struct {
	MerOrderId               string `json:"merOrderId" binding:"required"`
	QuickPayConfirmMsgCipher QuickPayConfirmMsgCipher
}

func (q *QuickPayConfirmReq) Valid() error {
	if q.MerOrderId == "" {
		return errors.New("merOrderId is required")
	}
	return q.QuickPayConfirmMsgCipher.Valid()
}

type QuickPayConfirmResp struct {
	/*
		resultCode 交易结果 4 详情参见 6.2 附录二 resultCode 0000
		errorCode 异常代码 10 详情参见 6.1 附录一 errorCode
		errorMsg 异常描述 200 中文、字母、数字
		ncountOrderId 新账通订单 32 商户支付流程订单凭证

		tranAmount 商户订单金 额 1-10 格式：数字 单位：元
		checkDate 对账日期 14 格式： YYYYMMDD
		submitTime 提交时间 同上送
		tranFinishTime 交易完成时 间 14 格式： YYYYMMDDHHMMSS 本域为订单的完成时间 20111007094626
	*/
	ResultCode     string      `json:"resultCode" binding:"required"`
	ErrorCode      string      `json:"errorCode" `
	ErrorMsg       string      `json:"errorMsg" `
	NcountOrderId  string      `json:"ncountOrderId" binding:"required"`
	TranAmount     interface{} `json:"tranAmount" binding:"required"`
	CheckDate      string      `json:"checkDate" binding:"required"`
	SubmitTime     string      `json:"submitTime" binding:"required"`
	TranFinishTime string      `json:"tranFinishTime" binding:"required"`

	/*
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
	*/
	BankCode          string `json:"bankCode" binding:"required"`
	CardType          string `json:"cardType" binding:"required"`
	ShortCardNo       string `json:"shortCardNo" binding:"required"`
	BindCardAgrNo     string `json:"bindCardAgrNo" binding:"required"`
	BusinessType      string `json:"businessType" binding:"required"`
	FeeAmount         string `json:"feeAmount" binding:"required"`
	DivideAcctDtl     string `json:"divideAcctDtl" binding:"required"`
	SignValue         string `json:"signValue" binding:"required"`
	RecvAcctAmount    string `json:"recvAcctAmount" `
	InstalmentNum     string `json:"instalmentNum" `
	PayableFeeAmt     string `json:"payableFeeAmt" `
	PayableFeeType    string `json:"payableFeeType" `
	FirstPeriodFeeAmt string `json:"firstPeriodFeeAmt" `
	EachPeriodFeeAmt  string `json:"eachPeriodFeeAmt" `
	/*
		firstPeriodPayAmt 首期还款金 额 0-12 注：仅在下单时分期期数 不为空才返回 例如：5
		instalmentRate 商户分期实 际贴息费率 6 注：仅在下单时分期期数 不为空才返回 例如：001212
		instalmentAmt 商户分期实 际贴息费用 0-12 注：仅在下单时分期期数 不为空才返回 例如：1.23
	*/

	FirstPeriodPayAmt string `json:"firstPeriodPayAmt" `
	InstalmentRate    string `json:"instalmentRate" `
	InstalmentAmt     string `json:"instalmentAmt" `
}
