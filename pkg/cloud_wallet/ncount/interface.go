package ncount

/*
// 创建用户账户地址
	NewAccountURL = "https://ncount.hnapay.com/api/r010.htm"

	// 用户查询接口
	checkUserAccountInfoURL = "https://ncount.hnapay.com/api/q001.htm"

	// 绑卡接口
	bindCardURL = "https://ncount.hnapay.com/api/r007.htm"

	// 绑卡确认接口
	bindCardConfirmURL = "https://ncount.hnapay.com/api/r008.htm"

	// 个人用户解绑接口
	unbindCardURL = "https://ncount.hnapay.com/api/r009.htm"

	// 用户账户明细查询接口：转账详情
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

*/

// NCounter is the platform ncount interface, the interface is provided
// by the platform, and the platform implements the interface.
type NCounter interface {
	// 创建用户账户地址
	NewAccount(req *NewAccountReq) (*NewAccountResp, error)
	// 用户查询接口
	CheckUserAccountInfo(req *CheckUserAccountReq) (*CheckUserAccountResp, error)
	// 绑卡接口
	BindCard(req *BindCardReq) (*BindCardResp, error)
	// 绑卡确认接口
	BindCardConfirm(req *BindCardConfirmReq) (*BindCardConfirmResp, error)
	// 个人用户解绑接口
	UnbindCard(req *UnBindCardReq) (*UnBindCardResp, error)
	// 用户账户明细查询接口：转账详情
	CheckUserAccountDetail(req *CheckUserAccountDetailReq) (*CheckUserAccountDetailResp, error)
	// 交易查询接口
	CheckUserAccountTrans(req *CheckUserAccountTransReq) (*CheckUserAccountTransResp, error)
	// 快捷支付下单接口
	QuickPayOrder(req *QuickPayOrderReq) (*QuickPayOrderResp, error)
	// 快捷支付确认接口
	QuickPayConfirm(req *QuickPayConfirmReq) (*QuickPayConfirmResp, error)
	// 转账接口
	Transfer(req *TransferReq) (*TransferResp, error)
	// 退款接口
	Refund(req *RefundReq) (*RefundResp, error)

	// 提现接口
	Withdraw(req *WithdrawReq) (*WithdrawResp, error)
}
