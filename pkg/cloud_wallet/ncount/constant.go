package ncount

const (
	// 新生支付公钥
	// MER_USER_ID = "300002428690"
	//PUBLIC_KEY  = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC9vGvYjivDF5uPBNDXMtoAtjYQ2YPSsfareduDG6kHL/N3A05rFHA11Dbr+UON82Y4V0RFKAQeZFPWcTLjcy6ntZVI8XoYLpuVQBPsb0Ya+PwbzR8/TmUdUf91ru8APtJgqkULgPVrO1hhzZ1tQMznosNLTOqbknMnnMcwzB5yYwIDAQAB\n-----END PUBLIC KEY-----"
	// PRIVATE_KEY = `MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAKPjggbqm+RZOvoWMe6To3LeBlLWS8027RGqpAIJjfLaEu1HXvXX6q3Vcww7pYlzxp4fBiTcEvZ1gjTPq+N8/KyiFWRAO9Xs68HrEN2eRa92n3Gsu3XJFSJ7OeUOwAZtXQw6XlB3iIRa9XR1ueXsx8NUoGrl4mJq1rlgEvA5KGUJAgMBAAECgYBOVadH4QmkatYaxVMWgvEELYV+QLm4nAFSiWqdIq37nyqeyZdlENA2SKkV9siX24Pa/l80bRCPRvl2frDdKlem88q6D8PfdBaPRYVr950xXRLG7AAmE7YND4O6B81pQ46je28tQ/3jzwBN54/TlJJVWWQP76m5Zo/PUD3zdxQiAQJBAN7v6YaMASuUmO5DeP2C8oAnUxABhssdRgqzhbZ73bqn7kGHYdWE3TZau52UCoy+KcYyGNxTuxQr3kWTUrj9S0ECQQC8Mb9SARqKILJGVwdGRlSAS5zgnR356/0NCTdP5vws2DeXhHV50jpaYsyXBCLyFkBXNwX+2qw0+qbuOf0of4/JAkAnWBPQiPjT5h+vPP0nUGrXrxj7pClTw1DPJqucbvPMs0JbEjdz5UTdCNo/jxbli9H3hnPYvnYvsyZBBST+PMWBAkEAqPADbgrdlydYwbn4JsaVroGx9xQzx5lnlN80Dv8sWtlRtitLBauJhH/yZpJpCGafJWuYbzo/omNrnKjjsAoquQJAGMNbUXZteQJ9B0uCbSRx0KpJkw3+Ibvf/L7VRs2HCKqXQgU1xlrKxv1kgc9jhOQvwMxTGLTUrD9NOXV2w+Kapw==`
	MER_USER_ID     = "300047240928"
	SUB_MERCHANT_ID = "2305241206593380534" // 子商户号
	PUBLIC_KEY      = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC9vGvYjivDF5uPBNDXMtoAtjYQ\n2YPSsfareduDG6kHL/N3A05rFHA11Dbr+UON82Y4V0RFKAQeZFPWcTLjcy6ntZVI\n8XoYLpuVQBPsb0Ya+PwbzR8/TmUdUf91ru8APtJgqkULgPVrO1hhzZ1tQMznosNL\nTOqbknMnnMcwzB5yYwIDAQAB\n-----END PUBLIC KEY-----"
	PRIVATE_KEY     = "MIICeAIBADANBgkqhkiG9w0BAQEFAASCAmIwggJeAgEAAoGBAOK8BAwd/FWL1XQU/pGG4nz5JAlSHgY97UejJxwtO6PX8vmHQLTdudLl8eOQbl/y4fX3vjLVsFanjgKNpeAb2M0Mk+sROaJkg2GR1YhbGH8bQ6ieJz9SuP7aZ6kCGVXJHHWYwPaor37dbAOGhQ17I+fKXmNiX1KT9eoeRRDLJF9BAgMBAAECgYEA2JSUoSZ7hRPfv9TWHxjjfFFYVPb16yx4XbfBgi7LC4UaebTy4FH0UTqJRsEOTeTqZ1RRgKmSmhPPmSzJSDwRaHDEMPkykgyezjHhXPPjJSvhIvzeI81l8TEZvRsuWsTglaXwqHGNxfuOufk9r8mQQ73mA9bDAmBbfVA/gVJB7gECQQD1RSDSRZXL3hSA7qeont3jBcTgbC4z5qcI8OrZMqqS7JMERfdFkcjsWHyLJklwR9qWy5lpmqeyE2kNXwfpQ7VxAkEA7KdM6zKWFS21nn4UzFoa20pEv4RqFNxw9puDYEO5x+7rfho30ap/+1AkJJ3TwbWhj8GH8TtzGS5jFYprUDoe0QJAWVj/ZdoXgZa7HWTTCqgk6Ii3eZGvGxURED7DLrA4VyF7RPk/5MYAzahGZmJiKlbimEA++KtwH3zWrhpKRX124QJBAInf9rpYoIP6O4P5ZNih7l+wZ1lFJiC9RbsHY4UkMArBscWoNLkcoq+iQ0xp/0MuNNByKmdrAWW8VtHn8RmuouECQQDZMazOi/jZFUf0unbhf7Xhe0z9hecKGISUsHsw2lKH2/pMdTaebB3JC6J075nHaHaYOB22Z7eVUKmVtdZ8coKl"
)

const (
	// 创建用户账户地址
	NewAccountURL = "https://ncount.hnapay.com/api/r010.htm"

	// 用户查询接口
	checkUserAccountURL = "https://ncount.hnapay.com/api/q001.htm"

	// 绑卡接口
	bindCardURL = "https://ncount.hnapay.com/api/r007.htm"

	// 绑卡确认接口
	bindCardConfirmURL = "https://ncount.hnapay.com/api/r008.htm"

	// 个人用户解绑接口
	unbindCardURL = "https://ncount.hnapay.com/api/r009.htm"

	// 用户账户详情接口
	checkUserAccountDetailURL = "https://ncount.hnapay.com/api/q004.htm"

	// 交易查询接口
	checkUserAccountTransURL = "https://ncount.hnapay.com/api/q002.htm"

	// 快捷支付下单接口
	quickPayOrderURL = "https://ncount.hnapay.com/api/t007.htm"

	// 快捷支付确认接口
	quickPayConfirmURL = "https://ncount.hnapay.com/api/t008.htm"

	// 转账接口
	transferURL = "https://ncount.hnapay.com/api/t003.htm"

	// 退款接口
	refundURL = "https://ncount.hnapay.com/api/t005.htm"

	// 提现接口
	withdrawURL = "https://ncount.hnapay.com/api/t002.htm"
)

const (
	ResultCodeSuccess   = "0000" // 交易成功
	ResultCodeFail      = "4444" // 交易失败
	ResultCodeInProcess = "9999" //交易进行中
	ResultCodeNoEffect  = "7777" //交易无效

)
