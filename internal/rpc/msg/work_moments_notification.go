package msg

import (
	"crazy_server/pkg/common/constant"
	"crazy_server/pkg/common/log"
	pbOffice "crazy_server/pkg/proto/office"
	sdk "crazy_server/pkg/proto/sdk_ws"
	"crazy_server/pkg/utils"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func WorkMomentSendNotification(operationID, recvID string, notificationMsg *pbOffice.WorkMomentNotificationMsg) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), recvID, notificationMsg)
	WorkMomentNotification(operationID, recvID, recvID, notificationMsg)
}

func WorkMomentNotification(operationID, sendID, recvID string, m proto.Message) {
	var tips sdk.TipsComm
	var err error
	marshaler := jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: false,
	}
	tips.JsonDetail, _ = marshaler.MarshalToString(m)
	n := &NotificationMsg{
		SendID:      sendID,
		RecvID:      recvID,
		MsgFrom:     constant.UserMsgType,
		ContentType: constant.WorkMomentNotification,
		SessionType: constant.SingleChatType,
		OperationID: operationID,
	}
	n.Content, err = proto.Marshal(&tips)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "proto.Marshal failed")
		return
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), string(n.Content))
	Notification(n)
}
