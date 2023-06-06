package base_info

import (
	crazy_server_sdk "crazy_server/pkg/proto/sdk_ws"
)

type ParamsCommFriend struct {
	OperationID string `json:"operationID" binding:"required"`
	ToUserID    string `json:"toUserID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}

type AddBlacklistReq struct {
	ParamsCommFriend
}
type AddBlacklistResp struct {
	CommResp
}

type ImportFriendReq struct {
	FriendUserIDList []string `json:"friendUserIDList" binding:"required"`
	OperationID      string   `json:"operationID" binding:"required"`
	FromUserID       string   `json:"fromUserID" binding:"required"`
}
type UserIDResult struct {
	UserID string `json:"userID"`
	Result int32  `json:"result"`
}
type ImportFriendResp struct {
	CommResp
	UserIDResultList []UserIDResult `json:"data"`
}

type AddFriendReq struct {
	ParamsCommFriend
	ReqMsg string `json:"reqMsg"`
}
type AddFriendResp struct {
	VerifySwitch int32 `json:"verifySwitch"`
	CommResp
}

type AddFriendResponseReq struct {
	ParamsCommFriend
	Flag      int32  `json:"flag" binding:"required,oneof=-1 0 1"`
	HandleMsg string `json:"handleMsg"`
}
type AddFriendResponseResp struct {
	CommResp
}

type DeleteFriendReq struct {
	ParamsCommFriend
}
type DeleteFriendResp struct {
	CommResp
}

type GetBlackListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}
type GetBlackListResp struct {
	CommResp
	BlackUserInfoList []*crazy_server_sdk.PublicUserInfo `json:"-"`
	Data              []map[string]interface {
	} `json:"data" swaggerignore:"true"`
}

//type PublicUserInfo struct {
//	UserID   string `json:"userID"`
//	Nickname string `json:"nickname"`
//	FaceUrl  string `json:"faceUrl"`
//	Gender   int32  `json:"gender"`
//}

type SetFriendRemarkReq struct {
	ParamsCommFriend
	Remark string `json:"remark"`
}
type SetFriendRemarkResp struct {
	CommResp
}

type RemoveBlackListReq struct {
	ParamsCommFriend
}
type RemoveBlackListResp struct {
	CommResp
}

type IsFriendReq struct {
	ParamsCommFriend
}
type Response struct {
	Friend bool `json:"isFriend"`
}
type IsFriendResp struct {
	CommResp
	Response Response `json:"data"`
}

type GetFriendsInfoReq struct {
	ParamsCommFriend
}
type GetFriendsInfoResp struct {
	CommResp
	FriendInfoList []*crazy_server_sdk.FriendInfo `json:"-"`
	Data           []map[string]interface {
	} `json:"data" swaggerignore:"true"`
}

type GetFriendListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}
type GetFriendListResp struct {
	CommResp
	FriendInfoList []*crazy_server_sdk.FriendInfo `json:"-"`
	Data           []map[string]interface {
	} `json:"data" swaggerignore:"true"`
}

type GetFriendApplyListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}
type GetFriendApplyListResp struct {
	CommResp
	FriendRequestList []*crazy_server_sdk.FriendRequest `json:"-"`
	Data              []map[string]interface {
	} `json:"data" swaggerignore:"true"`
}

type GetSelfApplyListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}
type GetSelfApplyListResp struct {
	CommResp
	FriendRequestList []*crazy_server_sdk.FriendRequest `json:"-"`
	Data              []map[string]interface {
	} `json:"data" swaggerignore:"true"`
}
