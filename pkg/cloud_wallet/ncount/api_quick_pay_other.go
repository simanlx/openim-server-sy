package ncount

// 这里是为了解决银行卡充值问题的渠道

//
// tranAmount 支付金额 1-12 格式：数字（以元为单 位） 不可 例如： 100
//payType 支付方式 1 2:银行卡卡号 3:绑卡协议号 不可 例如：2
//cardNo 支付银行 卡卡号 0-30 payType=2 不可空 可空 例 如 ： 611888 812128
//holderName 持卡人姓 名 0-40 payType=2 不可空 可空
//cardAvailableDate 信用卡有 效期 0-4 payType=2，且为 信用卡时不可空 可空 例 如 ： 0320 含 义 ： 2020 年 03 月 cvv2 信 用 卡 CVV2 0-3 payType=2，且为 信用卡时不可空 可空 例 如 ： 318
//mobileNo 银行签约 手机号 011 payType=2 不可空 可空
//identityType 证件类型 0-2 payType=2 不可空 暂仅支持 1 身份证 可空
//identityCode 证件号码 0-50 payType=2 不可空 可空
//bindCardAgrNo 绑卡协议 号 30 payType=3 不可空 可空
//notifyUrl 商户异步 通知地址 1-255 后台通知地址 不可 例 如 ： https:/ /www.x
//- xx.com /respo nse.do
//orderExpireTime 订单过期 时长 0-1440 订单过期时长（单位：分 钟） 可空
//userId 用户编号 1-32 协议支付时，必填，要素 支付时，可空 可空 例 如 ： 102121
//receiveUserId 收款方 ID 1-32 消费交易时，填收款方 ID 担保交易时，填商户 ID 不可 例如： 102121
//merUserIp 商户用户 IP 0-128 商户用户签约时所在的机 器 IP 地址 可空 例 如 ： 211.12. 38.88
//riskExpand 风控扩展 信息 0-80 风控扩展信息 可空
//goodsInfo 商品信息 0-80 商品信息 可空
//subMerchantId 商户渠道 进件 ID 0-100 商户渠道进件 ID 不可

type NAccountQuickPayOtherOther struct {
	TranAmount    string `json:"tranAmount" binding:"required"`   // 0.5
	PayType       string `json:"payType" binding:"required"`      // 2
	CardNo        string `json:"cardNo" binding:"required"`       // 支付银行卡卡号
	HolderName    string `json:"holderName" binding:"required"`   // 持卡人姓名
	MobileNo      string `json:"mobileNo" binding:"required"`     // 银行签约手机号
	IdentityType  string `json:"identityType" binding:"required"` // 1 身份证
	IdentityCode  string `json:"identityCode" binding:"required"` // 身份证号码
	NotifyUrl     string `json:"notifyUrl" binding:"required"`
	ReceiveUserId string `json:"receiveUserId" binding:"required"`
	SubMerchantId string `json:"subMerchantId" binding:"required"`
}
