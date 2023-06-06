package contrive_msg

import (
	"crazy_server/pkg/common/log"
	server_api_params "crazy_server/pkg/proto/sdk_ws"
	"encoding/json"
)

// 解散群聊消息推送
func DismissGroup(OperateID, UserID, GroupID string) error {
	GroupDismissMsg := &ContriveMessage{
		Data:    &GroupDismissMessage{GroupID: GroupID},
		MsgType: 11,
	}
	co, _ := json.Marshal(GroupDismissMsg)
	res := &ManagementSendMsg{
		OperationID:         OperateID,
		BusinessOperationID: OperateID,
		SendID:              UserID,
		SenderPlatformID:    1,
		Content: ContriveData{
			Data:        string(co),
			Description: "解散群聊消息",
			Extension:   "",
		},
		ContentType:     110, // 自定义消息
		SessionType:     2,   // 1 单聊 2 群聊
		IsOnlineOnly:    false,
		NotOfflinePush:  false,
		GroupID:         GroupID, // 接收方ID 群聊
		OfflinePushInfo: &server_api_params.OfflinePushInfo{},
	}

	coo, _ := json.Marshal(res)
	return SendMessage(OperateID, coo)
}

// 群聊： 推送抢红包消息
func RedPacketGrabPushToGroup(OperateID, SendPacketUserID, ClickPacketUserID, RedPacketID, SendUserName, ClickUserName, GroupID string) error {
	GroupDismissMsg := &ContriveMessage{
		Data: &RedPacketGrabMessage{
			RedPacketID:   RedPacketID,
			ClickUserID:   ClickPacketUserID,
			SendUserID:    SendPacketUserID,
			SendUserName:  SendUserName,
			ClickUserName: ClickUserName,
		},
		MsgType: MessageType_GrapRedPacket, // 抢红包消息
	}
	co, _ := json.Marshal(GroupDismissMsg)
	res := &ManagementSendMsg{
		OperationID:         OperateID,
		BusinessOperationID: OperateID,
		SendID:              SendPacketUserID,
		SenderPlatformID:    1,
		Content: ContriveData{
			Data:        string(co),
			Description: "群聊：抢红包推送",
			Extension:   "",
		},
		ContentType:     110, // 自定义消息
		SessionType:     2,   // 1 单聊 2 群聊
		IsOnlineOnly:    false,
		NotOfflinePush:  false,
		GroupID:         GroupID, // 接收方ID 群聊
		OfflinePushInfo: &server_api_params.OfflinePushInfo{},
	}
	coo, _ := json.Marshal(res)
	return SendMessage(OperateID, coo)
}

// 个人： 推送抢红包消息
func RedPacketGrabPushToUser(OperateID, SendMessageUserID, SendPacketUserID, RedPacketID, SendUserName, ClickUserName, ReceiveID string) error {
	GroupDismissMsg := &ContriveMessage{
		Data: &RedPacketGrabMessage{
			RedPacketID:   RedPacketID,
			SendUserID:    SendPacketUserID,
			SendUserName:  SendUserName,
			ClickUserName: ClickUserName,
		},
		MsgType: MessageType_GrapRedPacket, // 抢红包消息
	}
	co, _ := json.Marshal(GroupDismissMsg)
	res := &ManagementSendMsg{
		OperationID:         OperateID,
		BusinessOperationID: OperateID,
		SendID:              SendMessageUserID,
		SenderPlatformID:    1,
		Content: ContriveData{
			Data:        string(co),
			Description: "单聊：抢红包推送",
			Extension:   "",
		},
		ContentType:     110, // 自定义消息
		SessionType:     1,   // 1 单聊 2 群聊
		IsOnlineOnly:    false,
		NotOfflinePush:  false,
		RecvID:          ReceiveID,
		OfflinePushInfo: &server_api_params.OfflinePushInfo{},
	}
	coo, _ := json.Marshal(res)
	return SendMessage(OperateID, coo)
}

// 红包退回： 个人+ 群聊 发送红包退回消息
func SendRebackMessage(OperationID, redPacketID, content string, sessionID int, SenderID, ReciveID string) error {
	msg := ContriveMessage{
		Data: RedPacketBackMessage{
			RedPacketID: redPacketID,
			Content:     content,
		},
		MsgType: MessageType_RedPacketReturn,
	}
	co, _ := json.Marshal(msg)
	res := &ManagementSendMsg{
		OperationID:         OperationID,
		BusinessOperationID: OperationID,
		SendID:              SenderID,
		SenderPlatformID:    1,
		Content: ContriveData{
			Data:        string(co),
			Description: "红包退回消息",
			Extension:   "",
		},
		ContentType:     110,              // 自定义消息
		SessionType:     int32(sessionID), // 1 单聊 2 群聊
		IsOnlineOnly:    false,
		NotOfflinePush:  false,
		GroupID:         ReciveID, // 接收方ID 群聊
		RecvID:          ReciveID,
		OfflinePushInfo: &server_api_params.OfflinePushInfo{},
	}
	coo, _ := json.Marshal(res)
	log.Error(OperationID, string(coo))
	return SendMessage(OperationID, coo)
}

// 推送最佳手气红包消息 - 未接入
func SendRedPacketLuckyMessage(OperateID, SendPacketUserID, RedPacketID, LuckyUserName, GroupID string, spendTime int64) error {
	GroupDismissMsg := &ContriveMessage{
		Data: &RedPacketLuckyMessage{
			RedPacketID: RedPacketID,
			UserName:    LuckyUserName,
			SpendTime:   spendTime,
		},
		MsgType: MessageType_RedPacketLucky, // 最佳手气红包消息
	}
	co, _ := json.Marshal(GroupDismissMsg)
	res := &ManagementSendMsg{
		OperationID:         OperateID,
		BusinessOperationID: OperateID,
		SendID:              SendPacketUserID,
		SenderPlatformID:    1,
		Content: ContriveData{
			Data:        string(co),
			Description: "单聊：抢红包推送",
			Extension:   "",
		},
		ContentType:     110, // 自定义消息
		SessionType:     2,   // 1 单聊 2 群聊
		IsOnlineOnly:    false,
		NotOfflinePush:  false,
		GroupID:         GroupID, // 接收方ID 群聊
		OfflinePushInfo: &server_api_params.OfflinePushInfo{},
	}
	coo, _ := json.Marshal(res)
	return SendMessage(OperateID, coo)
}
