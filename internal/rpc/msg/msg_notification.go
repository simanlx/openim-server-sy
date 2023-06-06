package msg

import (
	"crazy_server/pkg/common/constant"
	"crazy_server/pkg/common/log"
	crazy_server_sdk "crazy_server/pkg/proto/sdk_ws"
	"crazy_server/pkg/utils"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func DeleteMessageNotification(opUserID, userID string, seqList []uint32, operationID string) {
	DeleteMessageTips := crazy_server_sdk.DeleteMessageTips{OpUserID: opUserID, UserID: userID, SeqList: seqList}
	MessageNotification(operationID, userID, userID, constant.DeleteMessageNotification, &DeleteMessageTips)
}

func MessageNotification(operationID, sendID, recvID string, contentType int32, m proto.Message) {
	log.Debug(operationID, utils.GetSelfFuncName(), "args: ", m.String(), contentType)
	var err error
	var tips crazy_server_sdk.TipsComm
	tips.Detail, err = proto.Marshal(m)
	if err != nil {
		log.Error(operationID, "Marshal failed ", err.Error(), m.String())
		return
	}

	marshaler := jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: false,
	}

	tips.JsonDetail, _ = marshaler.MarshalToString(m)
	var n NotificationMsg
	n.SendID = sendID
	n.RecvID = recvID
	n.ContentType = contentType
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = operationID
	n.Content, err = proto.Marshal(&tips)
	if err != nil {
		log.Error(operationID, "Marshal failed ", err.Error(), tips.String())
		return
	}
	Notification(&n)
}
