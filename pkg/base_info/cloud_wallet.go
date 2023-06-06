package base_info

type CheckUserHaveAccountResp struct {
	HaveAccount bool `json:"have_account"`
	Step        int  `json:"step"` // 1:未实名认证 2:未绑定银行卡 3:未设置交易密码 4:已开户
}

/*
message idCardRealNameAuthReq{
  uint32 userID = 1; //用户id
  string idCard = 2; //身份证
  string realName = 3; //真实姓名
  string mobile = 4; //手机号码
}
*/
type CreateUserAccount struct {
	IdCard   string `json:"id_card"`
	RealName string `json:"real_name"`
	Mobile   string `json:"mobile"`
}

type CreateUserAccountresp struct {
	Code int    `json:"code"` // 0:成功 1:失败 11: 未实名认证 12: 未绑定银行卡 13: 未设置交易密码
	Msg  string `json:"msg"`  // 失败原因
}

// ======================================================== 快速充值接口 ========================================

// 用户点击充值按钮
type QuickPayReq struct {
	Amount string `json:"amount"` // 充值金额
}

type QuickPayResp struct {
	CommResp     // 需要优先返回这个字段
	Ret      int `json:"ret"` // 0:成功 1:失败
}

// ======================================================== 发红包========================================
/*
  string userId = 1; //用户id
  int32 PacketType = 2; //红包类型(1个人红包、2群红包)
  int32 IsLucky = 3; //是否为拼手气红包
  int32 IsExclusive = 4; //是否为专属红包(0不是、1是)
  int32 ExclusiveUserID = 5; //专属红包接收者 和isExclusive
  string PacketTitle = 6; //红包标题
  float Amount = 7; //红包金额 单位：分
  int32 Number = 8; //红包个数

*/
// 用户点击充值按钮
type SendRedPacketReq struct {
	UserId      string `json:"user_id"`
	PacketType  int    `json:"packet_type"`
	IsLucky     int    `json:"is_lucky"`
	IsExclusive int    `json:"is_exclusive"`
	ExclusiveID string `json:"exclusive_id"`
	PacketTitle string `json:"packet_title"`
	Amount      string `json:"amount"`
	Number      int    `json:"number"`
	OperationID string `json:"operation_id"` // 链路追踪ID
}

type SendRedPacketResp struct {
	CommResp     // 需要优先返回这个字段
	Ret      int `json:"ret"` // 0:成功 1:失败
}
