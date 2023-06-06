package contrive_msg

import server_api_params "crazy_server/pkg/proto/sdk_ws"

// ===============================================================管理员推送单条消息===============================================================

// OpenIM上层
type ContriveData struct {
	Data        string `json:"data"`
	Description string `json:"description"`
	Extension   string `json:"extension"`
}

// 消息结构体
type ManagementSendMsg struct {
	OperationID         string `json:"operationID" binding:"required"`
	BusinessOperationID string `json:"businessOperationID"`
	SendID              string `json:"sendID" binding:"required"`
	GroupID             string `json:"groupID" `
	SenderNickname      string `json:"senderNickname" `
	SenderFaceURL       string `json:"senderFaceURL" `
	SenderPlatformID    int32  `json:"senderPlatformID"`
	//ForceList        []string                     `json:"forceList" `
	Content         ContriveData                       `json:"content" binding:"required" swaggerignore:"true"`
	ContentType     int32                              `json:"contentType" binding:"required"`
	SessionType     int32                              `json:"sessionType" binding:"required"`
	IsOnlineOnly    bool                               `json:"isOnlineOnly"`
	NotOfflinePush  bool                               `json:"notOfflinePush"`
	OfflinePushInfo *server_api_params.OfflinePushInfo `json:"offlinePushInfo"`
	// 2178158235
	RecvID string `json:"recvID" `
}

// ===============================================================下面是自定义消息体===============================================================

const (
	ContriveMessageGrapRedPacket = 1 + iota // 抢红包消息
	ContriveMessageRedPacket
)

type ContriveMessage struct {
	Data    interface{} `json:"data" binding:"required"`
	MsgType int32       `json:"msgType" binding:"required"` // 消息类型
}

const (
	MessageType_GrapRedPacket = 1 + iota
	MessageType_RedPacket

	MessageType_DismissGroup
	MessageType_RedPacketLucky
	MessageType_RedPacketReturn
)

// 发送红包消息红包结构消息
type RedPacketMessage struct {
	SendUserID       string `json:"sendUserID" binding:"required"`       // 发送方ID
	SendUserHeadImg  string `json:"sendUserHeadImg" binding:"required"`  // 发送方头像
	SendUserNickName string `json:"sendUserNickName" binding:"required"` // 发送方昵称
	RedPacketID      string `json:"redPacketID" binding:"required"`      // 红包ID
	RedPacketType    int32  `json:"redPacketType" binding:"required"`    // 红包类型 1 个人红包，2 群红包
	IsLucky          int32  `json:"isLucky" binding:"required"`          // 是否是拼手气红包 1 是 0 否
	IsExclusive      int32  `json:"isExclusive" binding:"required"`      // 是否是独享红包 1 是 0 否
	ExclusiveID      string `json:"exclusiveID" binding:"required"`      // 独享红包用户ID
	PacketTitle      string `json:"packetTitle" binding:"required"`      // 红包标题
}

// 抢红包消息结构体 // 谁抢了我的红包 ｜ 我抢了谁的红包
type RedPacketGrabMessage struct {
	RedPacketID   string `json:"redPacketID" binding:"required"`   // 红包ID
	SendUserID    string `json:"sendUserID" binding:"required"`    // 发送方ID
	ClickUserID   string `json:"clickUserID" binding:"required"`   // 抢红包的用户ID
	SendUserName  string `json:"sendUserName" binding:"required"`  // 发送方红包的姓名
	ClickUserName string `json:"clickUserName" binding:"required"` // 抢红包的用户ID
}

// 群解散消息
type GroupDismissMessage struct {
	GroupID string `json:"groupID" binding:"required"` // 群ID
}

// 最佳手气红包消息
type RedPacketLuckyMessage struct {
	RedPacketID string `json:"redPacketID" binding:"required"` // 红包ID
	UserName    string `json:"userName" binding:"required"`    // 用户昵称
	SpendTime   int64  `json:"spendTime" binding:"required"`   // 总的花费时间
}

// 红包退回消息
type RedPacketBackMessage struct {
	RedPacketID string `json:"redPacketID" binding:"required"` // 红包ID
	Content     string `json:"content" binding:"required"`     // 退回原因
}
