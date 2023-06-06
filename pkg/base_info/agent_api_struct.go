package base_info

// 推广员申请提交
type AgentApplyReq struct {
	Name        string `json:"name"  binding:"required"`          //用户姓名
	Mobile      string `json:"mobile"  binding:"required"`        //用户手机号码
	ChessUserId int64  `json:"chess_user_id"  binding:"required"` //互娱用户id
}

type BindAgentNumberReq struct {
	AgentNumber   int32  `json:"agent_number"  binding:"required"`   //推广员编号
	ChessUserId   int64  `json:"chess_user_id"  binding:"required"`  //互娱用户id
	ChessNickname string `json:"chess_nickname"  binding:"required"` //互娱用户昵称
}

type GetUserAgentInfoReq struct {
	ChessUserId int64 `json:"chess_user_id"  binding:"required"` //互娱用户id
}

type AgentAccountIncomeChartReq struct {
	DateType int32 `json:"date_type"` //日期类型 1(7天),2(半年) 默认7天
}

type AgentAccountRecordListReq struct {
	Date         string `json:"date"`          //日期
	BusinessType int32  `json:"business_type"` //业务类型
	Keyword      string `json:"keyword"`       //搜索关键字
	Page         int32  `json:"page"`
	Size         int32  `json:"size"`
}

type AgentGameShopBeanConfigReq struct {
	AgentNumber int32 `json:"agent_number"  binding:"required"` //推广员编号
}

type AgentBeanAccountRecordListReq struct {
	Date         string `json:"date"`          //日期
	BusinessType int32  `json:"business_type"` //业务类型
	Keyword      string `json:"keyword"`       //搜索关键字
	Page         int32  `json:"page"`
	Size         int32  `json:"size"`
}

type AgentBeanShopUpStatusReq struct {
	Status   int32 `json:"status"`    //状态(0下架、1上架)
	ConfigId int32 `json:"config_id"` //配置id
	IsAll    int32 `json:"is_all"`    //是否全部(1全部，0单个)
}

type AgentBeanShopUpdateReq struct {
	BeanShopConfig []*BeanShopConfig `json:"bean_shop_config"` //咖豆配置
}

type BeanShopConfig struct {
	BeanNumber     int64 `json:"bean_number"`
	GiveBeanNumber int32 `json:"give_bean_number"`
	Amount         int32 `json:"amount"`
	Status         int32 `json:"status"`
}

type AgentMemberListReq struct {
	Page    int32  `json:"page"`
	Size    int32  `json:"size"`
	Keyword string `json:"keyword"`  //搜索关键字
	OrderBy int32  `json:"order_by"` //排序(0默认-绑定时间倒序,1咖豆倒序,2咖豆正序,3贡献值倒序,4贡献值正序)
}

type AgentGiveMemberBeanReq struct {
	ChessUserId int64 `json:"chess_user_id"  binding:"required"`
	BeanNumber  int64 `json:"bean_number"  binding:"required"`
}

type ChessPurchaseBeanNotifyReq struct {
	OrderNo       string `json:"order_no"  binding:"required"`        //平台订单号
	NcountOrderNo string `json:"ncount_order_no"  binding:"required"` //新生支付订单号
}

type ChessShopPurchaseBeanReq struct {
	AgentNumber  int32  `json:"agent_number"  binding:"required"`   //推广员编号
	ChessUserId  int64  `json:"chess_user_id"  binding:"required"`  //互娱用户id
	ConfigId     int32  `json:"config_id"  binding:"required"`      //咖豆配置id
	ChessOrderNo string `json:"chess_order_no"  binding:"required"` //互娱订单号
}

type PlatformPurchaseBeanNotifyReq struct {
	ChessOrderNo   string `json:"chess_order_no"  binding:"required"`  //互娱订单号
	NcountOrderNo  string `json:"ncount_order_no"  binding:"required"` //新生支付订单号
	AgentNumber    int32  `json:"agent_number"  binding:"required"`    //推广员编号
	ChessUserId    int64  `json:"chess_user_id"  binding:"required"`   //互娱用户id
	BeanNumber     int64  `json:"bean_number"  binding:"required"`     //购买数量
	GiveBeanNumber int32  `json:"give_bean_number"`                    //赠送数量
	Amount         int32  `json:"amount"  binding:"required"`          //金额(单位分)
}

type PurchaseBeanReq struct {
	ConfigId int32 `json:"config_id"  binding:"required"` //咖豆配置id
}

type WithdrawReq struct {
	Amount          int32  `json:"amount"  binding:"required"`           //金额(单位分)
	PaymentPassword string `json:"payment_password"  binding:"required"` //支付密码
}

type NcountNotifyReq struct {
	OrderId    string `json:"OrderId"  binding:"required"`    //新生支付订单id
	MerOrderId string `json:"MerOrderId"  binding:"required"` //平台订单id
	Status     int32  `json:"Status"  binding:"required"`     //状态 100 未支付、200 支付成功、300支付失败
	PayTime    string `json:"PayTime"  binding:"required"`    //支付时间
	Amount     int32  `json:"Amount"  binding:"required"`     //金额(单位分)
}

type OpenAgentReq struct {
	ApplyId int32 `json:"apply_id"  binding:"required"` //申请id
}
