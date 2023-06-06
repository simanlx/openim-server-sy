package base_info

import (
	crazy_server_sdk "crazy_server/pkg/proto/sdk_ws"
)

type GetUsersInfoReq struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required"`
}
type GetUsersInfoResp struct {
	CommResp
	UserInfoList []*crazy_server_sdk.PublicUserInfo `json:"-"`
	Data         []map[string]interface{}           `json:"data" swaggerignore:"true"`
}

type UpdateSelfUserInfoReq struct {
	ApiUserInfo
	OperationID string `json:"operationID" binding:"required"`
}
type SetGlobalRecvMessageOptReq struct {
	OperationID      string `json:"operationID" binding:"required"`
	GlobalRecvMsgOpt *int32 `json:"globalRecvMsgOpt" binding:"omitempty,oneof=0 1 2"`
}
type SetGlobalRecvMessageOptResp struct {
	CommResp
}
type UpdateUserInfoResp struct {
	CommResp
}

type GetSelfUserInfoReq struct {
	OperationID string `json:"operationID" binding:"required"`
	UserID      string `json:"userID" binding:"required"`
}
type GetSelfUserInfoResp struct {
	CommResp
	UserInfo *crazy_server_sdk.UserInfo `json:"-"`
	Data     map[string]interface{}     `json:"data" swaggerignore:"true"`
}

type GetFriendIDListFromCacheReq struct {
	OperationID string `json:"operationID" binding:"required"`
}

type GetFriendIDListFromCacheResp struct {
	CommResp
	UserIDList []string `json:"userIDList" binding:"required"`
}

type GetBlackIDListFromCacheReq struct {
	OperationID string `json:"operationID" binding:"required"`
}

type GetBlackIDListFromCacheResp struct {
	CommResp
	UserIDList []string `json:"userIDList" binding:"required"`
}