package msg

import (
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/constant"
	"crazy_server/pkg/common/log"
	crazy_server_sdk "crazy_server/pkg/proto/sdk_ws"
	"crazy_server/pkg/utils"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func SetConversationNotification(operationID, sendID, recvID string, contentType int, m proto.Message, tips crazy_server_sdk.TipsComm) {
	log.NewInfo(operationID, "args: ", sendID, recvID, contentType, m.String(), tips.String())
	var err error
	tips.Detail, err = proto.Marshal(m)
	if err != nil {
		log.NewError(operationID, "Marshal failed ", err.Error(), m.String())
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
	n.ContentType = int32(contentType)
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = operationID
	n.Content, err = proto.Marshal(&tips)
	if err != nil {
		log.Error(operationID, utils.GetSelfFuncName(), "Marshal failed ", err.Error(), tips.String())
		return
	}
	Notification(&n)
}

// SetPrivate调用
func ConversationSetPrivateNotification(operationID, sendID, recvID string, isPrivateChat bool) {
	log.NewInfo(operationID, utils.GetSelfFuncName())
	conversationSetPrivateTips := &crazy_server_sdk.ConversationSetPrivateTips{
		RecvID:    recvID,
		SendID:    sendID,
		IsPrivate: isPrivateChat,
	}
	var tips crazy_server_sdk.TipsComm
	var tipsMsg string
	if isPrivateChat == true {
		tipsMsg = config.Config.Notification.ConversationSetPrivate.DefaultTips.OpenTips
	} else {
		tipsMsg = config.Config.Notification.ConversationSetPrivate.DefaultTips.CloseTips
	}
	tips.DefaultTips = tipsMsg
	SetConversationNotification(operationID, sendID, recvID, constant.ConversationPrivateChatNotification, conversationSetPrivateTips, tips)
}

// 会话改变
func ConversationChangeNotification(operationID, userID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName())
	ConversationChangedTips := &crazy_server_sdk.ConversationUpdateTips{
		UserID: userID,
	}
	var tips crazy_server_sdk.TipsComm
	tips.DefaultTips = config.Config.Notification.ConversationOptUpdate.DefaultTips.Tips
	SetConversationNotification(operationID, userID, userID, constant.ConversationOptChangeNotification, ConversationChangedTips, tips)
}

//会话未读数同步
func ConversationUnreadChangeNotification(operationID, userID, conversationID string, updateUnreadCountTime int64) {
	log.NewInfo(operationID, utils.GetSelfFuncName())
	ConversationChangedTips := &crazy_server_sdk.ConversationUpdateTips{
		UserID:                userID,
		ConversationIDList:    []string{conversationID},
		UpdateUnreadCountTime: updateUnreadCountTime,
	}
	var tips crazy_server_sdk.TipsComm
	tips.DefaultTips = config.Config.Notification.ConversationOptUpdate.DefaultTips.Tips
	SetConversationNotification(operationID, userID, userID, constant.ConversationUnreadNotification, ConversationChangedTips, tips)
}
