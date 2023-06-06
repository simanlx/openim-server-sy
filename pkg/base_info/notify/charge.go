package notify

// ChargeNotifyReq 充值回调

/*
version 版本号 同上送
tranCode 交易代码 同上送
merOrderId 商户订单号 同上送
merId 商户 ID 同上送
merAttach 附加数据 同上送
charset 编码方式 同上送
signType 签名类型 同上送
resultCode 交易结果 4 6.2 附录二 resultCode 0000 errorCode 异常代码 10 6.1 附录一
errorCode
errorMsg 异常描述 200 中文、字母、数字
ncountOrderId 新账通订单号 32 商户支付流程订单凭证
tranAmount 商户订单金额 1-10 格式：整数 单位：元
submitTime 提交时间 同上送
tranFinishTime 交易完成时间 14 格式： YYYYMMDDHHMMSS 本域为订单的完成时间 20111007094626
businessType 业务类型 2 03 消费，04 担保下单
feeAmount 交易手续费 格式：数字，单位：元
bankOrderId BIS 订单号 19 格式：数字 20101409271841 23217
realBankOrderId 银行单号 格式：数字
divideAcctDtl 分账订单明细divideId 为分账主订单 Id divideDtlId 为 分 账 明 细 IdledgerUserId 为分账方 Id divideStatus 为分账订单状 态 6.7 分账明细状态 [{"divideId":"123456", “divideDtlList": [{ "divideDtlId":"1234 56", "ledgerUserId":"123 456", "divideStatus":”1" }] }]
signValue 签名字符串 将报文信息用 signType 域设 置的方式加密后生成的字符

*/

//bankCode=CITIC&charset=1&bindCardAgrNo=202304030004579859&divideAcctDtl=&recvAcctAmount=10.98&ncountOrderId=2023040720951263&
//	cardType=1&resultCode=0000&errorCode=&tranFinishTime=20230407200250&checkDate=20230407&
//			version=1.0&errorMsg=&feeAmount=0&submitTime=20230407200203&shortCardNo=6812&signType=1&
//				merId=300002428690&merAttach=&tranCode=T008&businessType=03&tranAmount=0.01&merOrderId=20230407200203609866

// {"merOrderId":"20230407200203609866","resultCode":"0000","errorCode":"","errorMsg":"","ncountOrderId":"2023040720951263","tranAmount":"0.01","submitTime":"20230407200203","tranFinishTime":"20230407200250","feeAmount":"0"}

type ChargeNotifyReq struct {
	MerOrderId     string `json:"merOrderId"  form:"merOrderId"`
	ResultCode     string `json:"resultCode" form:"resultCode"`
	ErrorCode      string `json:"errorCode" form:"errorCode"`
	ErrorMsg       string `json:"errorMsg" form:"errorMsg"`
	NcountOrderId  string `json:"ncountOrderId" form:"ncountOrderId"`
	TranAmount     string `json:"tranAmount" form:"tranAmount"`
	SubmitTime     string `json:"submitTime" form:"submitTime"`
	TranFinishTime string `json:"tranFinishTime" form:"tranFinishTime"`
	FeeAmount      string `json:"feeAmount" form:"feeAmount"`
}

type ChargeNotifyResp struct {
	Code int `json:"code"`
}

/*
version 版本号 同上送
tranCode 交易代码 同上送
merOrderId 商户订单号 同上送
merId 商户 ID 同上送
merAttach 附加数据 同上送
charset 编码方式 同上送
signType 签名类型 同上送
resultCode 处理结果码 4 附录二 resultCode 9999
errorCode 异常代码 1-10 附录一 errorCode
errorMsg 异常描述 1-200 中文、字母、数字
ncountOrderId 订单号 32 新账通平台订单号
orderDate 平台订单日期 8 YYYYMMDD
tranFinishDate 订单完成日期 YYYYMMDD
signValue 签名字符串 将报文信息用
signType 域 设置的方式签名后生成的 字符串
serviceAmount 服务费 服务费，默认为 0
payAcctAmount 付款方账户余额 1-10 格式：整数 单位：元 交易成功时返回
*/

type WithdrawNotifyReq struct {
	MerOrderId     string `json:"merOrderId" form:"merOrderId"`
	ResultCode     string `json:"resultCode" form:"resultCode"`
	ErrorCode      string `json:"errorCode" form:"errorCode"`
	ErrorMsg       string `json:"errorMsg" form:"errorMsg"`
	NcountOrderId  string `json:"ncountOrderId" form:"ncountOrderId"`
	TranFinishDate string `json:"tranFinishDate" form:"tranFinishDate"`
	ServiceAmount  string `json:"serviceAmount" form:"serviceAmount"`
	PayAcctAmount  string `json:"payAcctAmount" form:"payAcctAmount"`
}
