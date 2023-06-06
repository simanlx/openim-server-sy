package msg

import (
	"crazy_server/pkg/common/constant"
	"crazy_server/pkg/common/log"
	//sdk "crazy_server/pkg/proto/sdk_ws"
	"crazy_server/pkg/utils"
	//"github.com/golang/protobuf/jsonpb"
	//"github.com/golang/protobuf/proto"
)

func SuperGroupNotification(operationID, sendID, recvID string) {
	n := &NotificationMsg{
		SendID:      sendID,
		RecvID:      recvID,
		MsgFrom:     constant.SysMsgType,
		ContentType: constant.SuperGroupUpdateNotification,
		SessionType: constant.SingleChatType,
		OperationID: operationID,
	}

	log.NewInfo(operationID, utils.GetSelfFuncName(), string(n.Content))
	Notification(n)
}
