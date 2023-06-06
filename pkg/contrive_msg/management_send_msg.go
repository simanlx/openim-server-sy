package contrive_msg

import (
	"bytes"
	imdb "crazy_server/pkg/common/db/mysql_model/im_mysql_model"
	"crazy_server/pkg/common/log"
	server_api_params "crazy_server/pkg/proto/sdk_ws"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type FPacket struct {
	ID              int64  `gorm:"column:id;primary_key;AUTO_INCREMENT;not null" json:"id"`
	PacketID        string `gorm:"column:packet_id;not null" json:"packet_id"`
	UserID          string `gorm:"column:user_id;not null" json:"user_id"`
	PacketType      int32  `gorm:"column:packet_type;not null" json:"packet_type"`
	IsLucky         int32  `gorm:"column:is_lucky;not null" json:"is_lucky"`
	ExclusiveUserID string `gorm:"column:exclusive_user_id;not null" json:"exclusive_user_id"`
	PacketTitle     string `gorm:"column:packet_title;not null" json:"packet_title"`
	Amount          int64  `gorm:"column:amount;not null" json:"amount"`
	Number          int32  `gorm:"column:number;not null" json:"number"`
	ExpireTime      int64  `gorm:"column:expire_time;not null" json:"expire_time"`
	MerOrderID      string `gorm:"column:mer_order_id;not null" json:"mer_order_id"`
	OperateID       string `gorm:"column:operate_id;not null" json:"operate_id"`
	RecvID          string `gorm:"column:recv_id;not null" json:"recv_id"`
	CreatedTime     int64  `gorm:"column:created_time;not null" json:"created_time"`
	UpdatedTime     int64  `gorm:"column:updated_time;not null" json:"updated_time"`
	Status          int32  `gorm:"column:status;not null" json:"status"` // 0 创建未生效，1 为红包正在领取中，2为红包领取完毕，3为红包过期
	IsExclusive     int32  `gorm:"column:is_exclusive;not null" json:"is_exclusive"`
}

const (
	Url = "http://server.jiadengni.com:10002/msg/manage_send_msg"
)

func SendGrabPacket(sendID, recevieID string, sessionID int32, OperateID, remark_click, remark_send, redPacketID string) {
	// 发送红包消息
	content1, content2 := NewManagementSendMsg_ClickPacket(sendID, recevieID, sessionID, OperateID, remark_click, remark_send, redPacketID)
	// 将消息发送给用户
	err := SendMessage(OperateID, content1)
	if err != nil {
		// todo  这里发送消息应该必须是可以重试的
		log.Error(OperateID, "发送消息失败", err)
	}

	err = SendMessage(OperateID, content2)
	if err != nil {
		// todo  这里发送消息应该必须是可以重试的
		log.Error(OperateID, "发送消息失败", err)
	}
}

// 发送红包消息
func SendSendRedPacket(f *FPacket, sessionID int) error {
	content, err := NewManagementSendMsg_RedMsg(f, f.OperateID, sessionID)
	if err != nil {
		log.Error(f.OperateID, "发送红包消息失败", err)
		return err
	}
	fmt.Println(string(content))
	// 将消息发送给用户
	err = SendMessage(f.OperateID, content)
	if err != nil {
		// todo  这里发送消息应该必须是可以重试的
		log.Error(f.OperateID, "发送消息失败", err)
		return err
	}
	return nil
}

// 发送两条抢红包消息，一条给发送方，一条给抢红包方
func NewManagementSendMsg_ClickPacket(sendID, recevieID string, sessionID int32, OperateID, clickUsername, sendUsername, redPacketID string) ([]byte, []byte) {
	// 创建msg 给对方发一条消息
	msg1 := ContriveMessage{
		Data: RedPacketGrabMessage{
			RedPacketID: redPacketID,
		},
		MsgType: ContriveMessageGrapRedPacket,
	}
	msg2 := ContriveMessage{
		Data: RedPacketGrabMessage{
			RedPacketID: redPacketID,
		},
		MsgType: ContriveMessageGrapRedPacket,
	}
	co1, _ := json.Marshal(msg1)
	co2, _ := json.Marshal(msg2)
	remarkSendContent := ContriveData{
		Data:        string(co2),
		Description: "你的红包被人抢了",
		Extension:   "",
	}
	remarkClickContent := ContriveData{
		Data:        string(co1),
		Description: "你抢到了红包",
		Extension:   "",
	}

	// 1. 发送一条给发送方的消息
	msgSend := newGrabRedPacket(sendID, recevieID, sessionID, OperateID, remarkSendContent)
	// 2. 发送一条给抢红包方的消息
	msgReceive := newGrabRedPacket(recevieID, sendID, sessionID, OperateID, remarkClickContent)

	return msgSend, msgReceive
}

// 创建发红包消息
func NewManagementSendMsg_RedMsg(f *FPacket, OperateID string, sessionID int) ([]byte, error) {
	usr, err := imdb.GetUserByUserID(f.UserID)
	if err != nil {
		return nil, err
	}

	contriveData := RedPacketMessage{
		SendUserID:       usr.UserID,
		SendUserHeadImg:  usr.FaceURL,
		SendUserNickName: usr.Nickname,
		RedPacketID:      f.PacketID,
		RedPacketType:    f.PacketType,
		IsLucky:          f.IsLucky,
		IsExclusive:      f.IsExclusive,
		ExclusiveID:      f.ExclusiveUserID,
		PacketTitle:      f.PacketTitle,
	}

	wrap := &ContriveMessage{
		MsgType: ContriveMessageRedPacket,
		Data:    contriveData,
	}

	co, _ := json.Marshal(wrap)

	res := &ManagementSendMsg{
		OperationID:         OperateID,
		BusinessOperationID: OperateID,
		SendID:              f.UserID,
		SenderPlatformID:    1,
		SenderFaceURL:       usr.FaceURL,
		SenderNickname:      usr.Nickname,
		Content: ContriveData{
			Data:        string(co),
			Description: "红包消息",
			Extension:   "",
		},
		ContentType:     110,              // 自定义消息
		SessionType:     int32(sessionID), // 1 单聊 2 群聊
		IsOnlineOnly:    false,
		NotOfflinePush:  false,
		GroupID:         f.RecvID, // 接收方ID 群聊
		RecvID:          f.RecvID, // 接收方ID 用户
		OfflinePushInfo: &server_api_params.OfflinePushInfo{},
	}
	co1, _ := json.Marshal(res)
	return co1, nil
}

// 创建抢红包消息
func newGrabRedPacket(sendID, ReceID string, SessionType int32, OperateID string, remark ContriveData) []byte {
	// 1. 发送一条给发送方的消息
	msg := &ManagementSendMsg{
		OperationID:         OperateID,
		BusinessOperationID: OperateID,
		SendID:              sendID,
		GroupID:             ReceID,
		Content:             remark,
		ContentType:         110, // 自定义消息
		SessionType:         SessionType,
		IsOnlineOnly:        false,
		NotOfflinePush:      false,
		OfflinePushInfo: &server_api_params.OfflinePushInfo{
			Title:         "红包来了",
			Desc:          "红包来了",
			IOSBadgeCount: false,
		},
		RecvID: ReceID,
	}
	co, _ := json.Marshal(msg)
	return co
}

// 红包退回消息

// 调用post发送消息
func SendMessage(OperateID string, content []byte) error {
	// 发送请求
	resp, err := http.Post(Url, "application/json", bytes.NewBuffer(content))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// 如果响应返回的不是200，则表示发送失败
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("send message failed: %s", string(body))
	}
	log.Error(OperateID, string(body))
	return nil
}
