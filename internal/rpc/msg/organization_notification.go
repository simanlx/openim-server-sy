package msg

import (
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/constant"
	imdb "crazy_server/pkg/common/db/mysql_model/im_mysql_model"
	"crazy_server/pkg/common/log"
	utils2 "crazy_server/pkg/common/utils"
	crazy_server_sdk "crazy_server/pkg/proto/sdk_ws"
	"crazy_server/pkg/utils"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func OrganizationNotificationToAll(opUserID string, operationID string) {
	err, userIDList := imdb.GetAllOrganizationUserID()
	if err != nil {
		log.Error(operationID, "GetAllOrganizationUserID failed ", err.Error())
		return
	}

	tips := crazy_server_sdk.OrganizationChangedTips{OpUser: &crazy_server_sdk.UserInfo{}}

	user, err := imdb.GetUserByUserID(opUserID)
	if err != nil {
		log.NewError(operationID, "GetUserByUserID failed ", err.Error(), opUserID)
		return
	}
	utils2.UserDBCopyOpenIM(tips.OpUser, user)

	for _, v := range userIDList {
		log.Debug(operationID, "OrganizationNotification", opUserID, v, constant.OrganizationChangedNotification, &tips, operationID)
		OrganizationNotification(config.Config.Manager.AppManagerUid[0], v, constant.OrganizationChangedNotification, &tips, operationID)
	}
}

func OrganizationNotification(opUserID string, recvUserID string, contentType int32, m proto.Message, operationID string) {
	log.Info(operationID, utils.GetSelfFuncName(), "args: ", contentType, opUserID)
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

	switch contentType {
	case constant.OrganizationChangedNotification:
		tips.DefaultTips = "OrganizationChangedNotification"

	default:
		log.Error(operationID, "contentType failed ", contentType)
		return
	}

	var n NotificationMsg
	n.SendID = opUserID
	n.RecvID = recvUserID
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
