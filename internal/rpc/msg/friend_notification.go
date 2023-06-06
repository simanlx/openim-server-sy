package msg

import (
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/constant"
	imdb "crazy_server/pkg/common/db/mysql_model/im_mysql_model"
	"crazy_server/pkg/common/log"
	utils2 "crazy_server/pkg/common/utils"
	pbFriend "crazy_server/pkg/proto/friend"
	crazy_server_sdk "crazy_server/pkg/proto/sdk_ws"
	"crazy_server/pkg/utils"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func getFromToUserNickname(fromUserID, toUserID string) (string, string, error) {
	from, err := imdb.GetUserByUserID(fromUserID)
	if err != nil {
		return "", "", utils.Wrap(err, "")
	}
	to, err := imdb.GetUserByUserID(toUserID)
	if err != nil {
		return "", "", utils.Wrap(err, "")
	}
	return from.Nickname, to.Nickname, nil
}

func friendNotification(commID *pbFriend.CommID, contentType int32, m proto.Message) {
	log.Info(commID.OperationID, utils.GetSelfFuncName(), "args: ", commID, contentType)
	var err error
	var tips crazy_server_sdk.TipsComm
	tips.Detail, err = proto.Marshal(m)
	if err != nil {
		log.Error(commID.OperationID, "Marshal failed ", err.Error(), m.String())
		return
	}

	marshaler := jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: false,
	}

	tips.JsonDetail, _ = marshaler.MarshalToString(m)

	fromUserNickname, toUserNickname, err := getFromToUserNickname(commID.FromUserID, commID.ToUserID)
	if err != nil {
		log.Error(commID.OperationID, "getFromToUserNickname failed ", err.Error(), commID.FromUserID, commID.ToUserID)
		return
	}
	cn := config.Config.Notification
	switch contentType {
	case constant.FriendApplicationNotification:
		tips.DefaultTips = fromUserNickname + cn.FriendApplication.DefaultTips.Tips
	case constant.FriendApplicationApprovedNotification:
		tips.DefaultTips = fromUserNickname + cn.FriendApplicationApproved.DefaultTips.Tips
	case constant.FriendApplicationRejectedNotification:
		tips.DefaultTips = fromUserNickname + cn.FriendApplicationRejected.DefaultTips.Tips
	case constant.FriendAddedNotification:
		tips.DefaultTips = cn.FriendAdded.DefaultTips.Tips
	case constant.FriendDeletedNotification:
		tips.DefaultTips = cn.FriendDeleted.DefaultTips.Tips + toUserNickname
	case constant.FriendRemarkSetNotification:
		tips.DefaultTips = fromUserNickname + cn.FriendRemarkSet.DefaultTips.Tips
	case constant.BlackAddedNotification:
		tips.DefaultTips = cn.BlackAdded.DefaultTips.Tips
	case constant.BlackDeletedNotification:
		tips.DefaultTips = cn.BlackDeleted.DefaultTips.Tips + toUserNickname
	case constant.UserInfoUpdatedNotification:
		tips.DefaultTips = cn.UserInfoUpdated.DefaultTips.Tips
	case constant.FriendInfoUpdatedNotification:
		tips.DefaultTips = cn.FriendInfoUpdated.DefaultTips.Tips + toUserNickname
	default:
		log.Error(commID.OperationID, "contentType failed ", contentType)
		return
	}

	var n NotificationMsg
	n.SendID = commID.FromUserID
	n.RecvID = commID.ToUserID
	n.ContentType = contentType
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = commID.OperationID
	n.Content, err = proto.Marshal(&tips)
	if err != nil {
		log.Error(commID.OperationID, "Marshal failed ", err.Error(), tips.String())
		return
	}
	Notification(&n)
}

func FriendApplicationNotification(req *pbFriend.AddFriendReq) {
	FriendApplicationTips := crazy_server_sdk.FriendApplicationTips{FromToUserID: &crazy_server_sdk.FromToUserID{}}
	FriendApplicationTips.FromToUserID.FromUserID = req.CommID.FromUserID
	FriendApplicationTips.FromToUserID.ToUserID = req.CommID.ToUserID
	friendNotification(req.CommID, constant.FriendApplicationNotification, &FriendApplicationTips)
}

func FriendApplicationApprovedNotification(req *pbFriend.AddFriendResponseReq) {
	FriendApplicationApprovedTips := crazy_server_sdk.FriendApplicationApprovedTips{FromToUserID: &crazy_server_sdk.FromToUserID{}}
	FriendApplicationApprovedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	FriendApplicationApprovedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	FriendApplicationApprovedTips.HandleMsg = req.HandleMsg
	friendNotification(req.CommID, constant.FriendApplicationApprovedNotification, &FriendApplicationApprovedTips)
}

func FriendApplicationRejectedNotification(req *pbFriend.AddFriendResponseReq) {
	FriendApplicationApprovedTips := crazy_server_sdk.FriendApplicationApprovedTips{FromToUserID: &crazy_server_sdk.FromToUserID{}}
	FriendApplicationApprovedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	FriendApplicationApprovedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	FriendApplicationApprovedTips.HandleMsg = req.HandleMsg
	friendNotification(req.CommID, constant.FriendApplicationRejectedNotification, &FriendApplicationApprovedTips)
}

func FriendAddedNotification(operationID, opUserID, fromUserID, toUserID string) {
	friendAddedTips := crazy_server_sdk.FriendAddedTips{Friend: &crazy_server_sdk.FriendInfo{}, OpUser: &crazy_server_sdk.PublicUserInfo{}}
	user, err := imdb.GetUserByUserID(opUserID)
	if err != nil {
		log.NewError(operationID, "GetUserByUserID failed ", err.Error(), opUserID)
		return
	}
	utils2.UserDBCopyOpenIMPublicUser(friendAddedTips.OpUser, user)
	friend, err := imdb.GetFriendRelationshipFromFriend(fromUserID, toUserID)
	if err != nil {
		log.NewError(operationID, "GetFriendRelationshipFromFriend failed ", err.Error(), fromUserID, toUserID)
		return
	}
	utils2.FriendDBCopyOpenIM(friendAddedTips.Friend, friend)
	commID := pbFriend.CommID{FromUserID: fromUserID, ToUserID: toUserID, OpUserID: opUserID, OperationID: operationID}
	friendNotification(&commID, constant.FriendAddedNotification, &friendAddedTips)
}

func FriendDeletedNotification(req *pbFriend.DeleteFriendReq) {
	friendDeletedTips := crazy_server_sdk.FriendDeletedTips{FromToUserID: &crazy_server_sdk.FromToUserID{}}
	friendDeletedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	friendDeletedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	friendNotification(req.CommID, constant.FriendDeletedNotification, &friendDeletedTips)
}

func FriendRemarkSetNotification(operationID, opUserID, fromUserID, toUserID string) {
	friendInfoChangedTips := crazy_server_sdk.FriendInfoChangedTips{FromToUserID: &crazy_server_sdk.FromToUserID{}}
	friendInfoChangedTips.FromToUserID.FromUserID = fromUserID
	friendInfoChangedTips.FromToUserID.ToUserID = toUserID
	commID := pbFriend.CommID{FromUserID: fromUserID, ToUserID: toUserID, OpUserID: opUserID, OperationID: operationID}
	friendNotification(&commID, constant.FriendRemarkSetNotification, &friendInfoChangedTips)
}

func BlackAddedNotification(req *pbFriend.AddBlacklistReq) {
	blackAddedTips := crazy_server_sdk.BlackAddedTips{FromToUserID: &crazy_server_sdk.FromToUserID{}}
	blackAddedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	blackAddedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	friendNotification(req.CommID, constant.BlackAddedNotification, &blackAddedTips)
}

func BlackDeletedNotification(req *pbFriend.RemoveBlacklistReq) {
	blackDeletedTips := crazy_server_sdk.BlackDeletedTips{FromToUserID: &crazy_server_sdk.FromToUserID{}}
	blackDeletedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	blackDeletedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	friendNotification(req.CommID, constant.BlackDeletedNotification, &blackDeletedTips)
}

//send to myself
func UserInfoUpdatedNotification(operationID, opUserID string, changedUserID string) {
	selfInfoUpdatedTips := crazy_server_sdk.UserInfoUpdatedTips{UserID: changedUserID}
	commID := pbFriend.CommID{FromUserID: opUserID, ToUserID: changedUserID, OpUserID: opUserID, OperationID: operationID}
	friendNotification(&commID, constant.UserInfoUpdatedNotification, &selfInfoUpdatedTips)
}

func FriendInfoUpdatedNotification(operationID, changedUserID string, needNotifiedUserID string, opUserID string) {
	selfInfoUpdatedTips := crazy_server_sdk.UserInfoUpdatedTips{UserID: changedUserID}
	commID := pbFriend.CommID{FromUserID: opUserID, ToUserID: needNotifiedUserID, OpUserID: opUserID, OperationID: operationID}
	friendNotification(&commID, constant.FriendInfoUpdatedNotification, &selfInfoUpdatedTips)
}
