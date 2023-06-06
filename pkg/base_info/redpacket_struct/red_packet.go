package redpacket_struct

// 发送红包
type SendRedPacket struct {
	UserId          string `json:"userId" binding:"required"`      //用户id
	PacketType      int32  `json:"PacketType" binding:"required"`  //红包类型(1个人红包、2群红包)
	IsLucky         int32  `json:"IsLucky" `                       //是否为拼手气红包
	IsExclusive     int32  `json:"IsExclusive"`                    //是否为专属红包(0不是、1是)
	ExclusiveUserID string `json:"ExclusiveUserID"`                //专属红包接收者 和isExclusive
	PacketTitle     string `json:"PacketTitle" binding:"required"` //红包标题
	Amount          int64  `json:"Amount" binding:"required"`      //红包金额 单位：分
	Number          int32  `json:"Number" binding:"required"`      //红包个数
	SendType        int32  `json:"SendType" binding:"required"`    //发送方式(1钱包余额、2银行卡)
	BankCardID      int64  `json:"BankCardID"`                     //银行卡id
	OperationID     string `json:"operationID" binding:"required"` //链路跟踪id
	RecvID          string `json:"recvID" binding:"required"`      //接收者id
	Password        string `json:"password" binding:"required"`    //支付密码

	// 	BindCardAgrNo string `protobuf:"bytes,13,opt,name=bindCardAgrNo,proto3" json:"bindCardAgrNo,omitempty"` //绑卡协议号
	BindCardAgrNo string `json:"bindCardAgrNo"` //绑卡协议号
}

// 发送红包响应
type SendRedPacketResp struct {
}

// 点击抢红包
type ClickRedPacketReq struct {
	UserId      string `json:"userId"`      //用户id
	RedPacketID string `json:"redPacketID"` //红包id
	OperateID   string `json:"operateID"`   //链路跟踪id
}

// 确认发送红包
type ConfirmSendRedPacketReq struct {
	Code        string `json:"code" binding:"required"`        //验证码
	RedPacketID string `json:"redPacketID" binding:"required"` //红包id
	OperateID   string `json:"operateID"`                      //链路跟踪id
}

// 红包领取明细列表记录
type RedPacketReceiveDetailReq struct {
	StartTime   string `json:"start_time" binding:"required"`
	EndTime     string `json:"end_time" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}

// 红包详情接口
type RedPacketInfoReq struct {
	UserId      string `json:"user_id"`                        //用户id
	PacketId    string `json:"packet_id"  binding:"required"`  //红包id
	OperationID string `json:"operationID" binding:"required"` //链路跟踪id
}

// 禁止用户抢红包 (禁止用户抢红包：未做)
type BanRedPacketReq struct {
	UserId      string `json:"user_id"`     //用户id
	Forbid      int32  `json:"forbid"`      //是否禁止抢红包(0不禁止、1禁止)
	GroupId     string `json:"group_id"`    //群id
	OperationID string `json:"operationID"` //链路跟踪id
}

// json : {"forbid":0,"group_id":"123456","operationID":"123456","user_id":"123456"}

// 获取声网token (获取声网token： 未做)
type AgoraTokenReq struct {
	Channel_name string `json:"ChannelName"` // 频道名称，如果是个人，就用用户的ID= UserID：single ，如果是群，就用群的ID= GroupID：group
	Role         uint32 `json:"role"`        // 用户角色，1是发起者，2是接受者
	OperationID  string `json:"operationID"` // 链路跟踪id
}

type AgoraTokenResp struct {
	Token string `json:"token"`
	AppID string `json:"appID"`
}

// 腾讯云进行消息转义 (腾讯云进行消息转义： 未做)
type TencentMsgEscapeReq struct {
	ContentUrl  string `json:"content_url" binding:"required"` //消息内容
	OperationID string `json:"operationID"`                    // 链路跟踪id
}

type GetVersionReq struct {
	// param out : 最新版本号、下载地址、更新内容、是否强制更新
	VersionCode string `json:"version_code" binding:"required"` //版本号
	OperationID string `json:"operationID" binding:"required"`  // 链路跟踪id
}

const (
	ReFoundPacketSecret = "dksidaksdkxhdkd"
)

type ReFoundPacketReq struct {
	// 红包退还接口，不对外开放
	Secret      string `json:"secret" binding:"required"`      //秘钥
	OperationID string `json:"operationID" binding:"required"` // 链路跟踪id
}

type ThirdPayReq struct {
	// 暴露第三方接口
	OprationID       string `json:"opration_id" binding:"required"` //链路跟踪id
	OrderNo          string `json:"order_no" binding:"required"`    //本平台订单号
	Password         string `json:"password" binding:"required"`    //支付密码
	SendType         int32  `json:"send_type" binding:"required"`   //发送方式(1钱包余额、2银行卡)
	BankcardProtocol string `json:"bankcard_protocol"`              //协议号
}

type CreateThirdPayOrder struct {
	// 创建第三方订单
	MerchantID string `json:"merchant_id" binding:"required"`  //商户号，需要向平台申请
	MerOrderID string `json:"mer_order_id" binding:"required"` //商户订单号 ，全局唯一，不能重复
	NotifyURL  string `json:"notify_url" binding:"required"`   //异步通知地址
	Amount     int32  `json:"amount" binding:"required"`       //金额，单位分
	Remark     string `json:"remark"`                          //备注
}

type GetThirdPayOrder struct {
	// 获取第三方订单
	OrderNO     string `json:"order_no" binding:"required"`     //订单号
	OperationID string `json:"operation_id" binding:"required"` // 链路跟踪id
}

// 通过业务ID来区别确认支付类型
type ThirdPayConfirm struct {
	OrderNo      string `json:"order_no" binding:"required"`      //订单号
	Code         string `json:"code" binding:"required"`          //验证码
	BusinessType int32  `json:"business_type" binding:"required"` //业务类型
	OperationID  string `json:"operation_id" binding:"required"`  // 链路跟踪id
}

// 充值回调
type ThirdPayCallback struct {
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

// 提现，咖豆提现到云钱包
type ThirdWithdrawReq struct {
	ThirdOrderId string `json:"thirdOrderId" binding:"required"` //第三方订单号
	Amount       int32  `json:"amount" binding:"required"`       //金额，单位分
	Commission   int32  `json:"commission" binding:"required"`   //手续费，单位分
	Password     string `json:"password" binding:"required"`     //支付密码
	OperationID  string `json:"operationID" binding:"required"`  // 链路跟踪id
}
