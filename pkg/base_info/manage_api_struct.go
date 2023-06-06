package base_info

import (
	pbRelay "crazy_server/pkg/proto/relay"
	server_api_params "crazy_server/pkg/proto/sdk_ws"
	pbUser "crazy_server/pkg/proto/user"
)

type DeleteUsersReq struct {
	OperationID      string   `json:"operationID" binding:"required"`
	DeleteUserIDList []string `json:"deleteUserIDList" binding:"required"`
}
type DeleteUsersResp struct {
	CommResp
	FailedUserIDList []string `json:"data"`
}
type GetAllUsersUidReq struct {
	OperationID string `json:"operationID" binding:"required"`
}
type GetAllUsersUidResp struct {
	CommResp
	UserIDList []string `json:"data"`
}
type GetUsersOnlineStatusReq struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required,lte=200"`
}
type GetUsersOnlineStatusResp struct {
	CommResp
	SuccessResult []*pbRelay.GetUsersOnlineStatusResp_SuccessResult `json:"data"`
}
type AccountCheckReq struct {
	OperationID     string   `json:"operationID" binding:"required"`
	CheckUserIDList []string `json:"checkUserIDList" binding:"required,lte=100"`
}
type AccountCheckResp struct {
	CommResp
	ResultList []*pbUser.AccountCheckResp_SingleUserStatus `json:"data"`
}

type ManagementSendMsg struct {
	OperationID         string `json:"operationID" binding:"required"`
	BusinessOperationID string `json:"businessOperationID"`
	SendID              string `json:"sendID" binding:"required"`
	GroupID             string `json:"groupID" `
	SenderNickname      string `json:"senderNickname" `
	SenderFaceURL       string `json:"senderFaceURL" `
	SenderPlatformID    int32  `json:"senderPlatformID"`
	//ForceList        []string                     `json:"forceList" `
	Content         map[string]interface{}             `json:"content" binding:"required" swaggerignore:"true"`
	ContentType     int32                              `json:"contentType" binding:"required"`
	SessionType     int32                              `json:"sessionType" binding:"required"`
	IsOnlineOnly    bool                               `json:"isOnlineOnly"`
	NotOfflinePush  bool                               `json:"notOfflinePush"`
	OfflinePushInfo *server_api_params.OfflinePushInfo `json:"offlinePushInfo"`
}

type ManagementSendMsgReq struct {
	ManagementSendMsg
	RecvID string `json:"recvID" `
}

type ManagementSendMsgResp struct {
	CommResp
	ResultList server_api_params.UserSendMsgResp `json:"data"`
}

type ManagementBatchSendMsgReq struct {
	ManagementSendMsg
	IsSendAll  bool     `json:"isSendAll"`
	RecvIDList []string `json:"recvIDList"`
}

type ManagementBatchSendMsgResp struct {
	CommResp
	Data struct {
		ResultList   []*SingleReturnResult `json:"resultList"`
		FailedIDList []string
	} `json:"data"`
}
type SingleReturnResult struct {
	ServerMsgID string `json:"serverMsgID"`
	ClientMsgID string `json:"clientMsgID"`
	SendTime    int64  `json:"sendTime"`
	RecvID      string `json:"recvID" `
}

type CheckMsgIsSendSuccessReq struct {
	OperationID string `json:"operationID"`
}

type CheckMsgIsSendSuccessResp struct {
	CommResp
	Status int32 `json:"status"`
}

type GetUsersReq struct {
	OperationID string `json:"operationID" binding:"required"`
	UserName    string `json:"userName"`
	UserID      string `json:"userID"`
	Content     string `json:"content"`
	PageNumber  int32  `json:"pageNumber" binding:"required"`
	ShowNumber  int32  `json:"showNumber" binding:"required"`
}

type CMSUser struct {
	UserID           string `json:"userID"`
	Nickname         string `json:"nickname"`
	FaceURL          string `json:"faceURL"`
	Gender           int32  `json:"gender"`
	PhoneNumber      string `json:"phoneNumber"`
	Birth            uint32 `json:"birth"`
	Email            string `json:"email"`
	Ex               string `json:"ex"`
	CreateIp         string `json:"createIp"`
	CreateTime       uint32 `json:"createTime"`
	LastLoginIp      string `json:"LastLoginIp"`
	LastLoginTime    uint32 `json:"LastLoginTime"`
	AppMangerLevel   int32  `json:"appMangerLevel"`
	GlobalRecvMsgOpt int32  `json:"globalRecvMsgOpt"`
	IsBlock          bool   `json:"isBlock"`
}

type GetUsersResp struct {
	CommResp
	Data struct {
		UserList    []*CMSUser `json:"userList"`
		TotalNum    int32      `json:"totalNum"`
		CurrentPage int32      `json:"currentPage"`
		ShowNumber  int32      `json:"showNumber"`
	} `json:"data"`
}

type AttributeSwitchReq struct {
	OperationID string `json:"operationID" binding:"required"`
}

type AttributeSwitchSetReq struct {
	SetType     int32  `json:"set_type"`
	SetValue    int32  `json:"set_value"`
	OperationID string `json:"operationID" binding:"required"`
}

type AttributeSwitchSetResp struct {
	CommResp
}

type AttributeMenuReq struct {
	OperationID string `json:"operationID" binding:"required"`
}

type WgtVersionReq struct {
	AppId       string `json:"app_id" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}

// 用户反馈、常见问题
type FeedbackReq struct {
	OperationID     string `json:"operationID" binding:"required"`
	FeedbackType    int32  `json:"feedbackType" binding:"required"`
	FeedbackContent string `json:"feedbackContent" binding:"required"`
	FeedbackContact string `json:"feedbackContact" `
}

type CommonQuestionReq struct {
	OperationID string `json:"operationID" binding:"required"`
}

type CommonQuestionFeedbackReq struct {
	OperationID string `json:"operationID" binding:"required"`
	QuestionID  int64  `json:"questionID" binding:"required"`
	Solved      int32  `json:"solved" ` // 1 . 0 未解决 1 已解决
}

type LatestVersionReq struct {
	AppType int32 `json:"app_type"  binding:"required"` //app类型 (1安卓、2ios、3....)
}

type FilterContentReq struct {
	OperationID string `json:"operationID" binding:"required"`
	Content     string `json:"content" binding:"required"`
}
