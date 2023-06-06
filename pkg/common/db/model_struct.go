package db

import "time"

type Register struct {
	Account        string `gorm:"column:account;primary_key;type:char(255)" json:"account"`
	Password       string `gorm:"column:password;type:varchar(255)" json:"password"`
	Ex             string `gorm:"column:ex;size:1024" json:"ex"`
	UserID         string `gorm:"column:user_id;type:varchar(255)" json:"userID"`
	AreaCode       string `gorm:"column:area_code;type:varchar(255)"`
	InvitationCode string `gorm:"column:invitation_code;type:varchar(255)"`
	RegisterIP     string `gorm:"column:register_ip;type:varchar(255)"`
}

type Invitation struct {
	InvitationCode string    `gorm:"column:invitation_code;primary_key;type:varchar(32)"`
	CreateTime     time.Time `gorm:"column:create_time"`
	UserID         string    `gorm:"column:user_id;index:userID"`
	LastTime       time.Time `gorm:"column:last_time"`
	Status         int32     `gorm:"column:status"`
}

// message FriendInfo{
// string OwnerUserID = 1;
// string Remark = 2;
// int64 CreateTime = 3;
// UserInfo FriendUser = 4;
// int32 AddSource = 5;
// string OperatorUserID = 6;
// string Ex = 7;
// }
// crazy_server_sdk.FriendInfo(FriendUser) != imdb.Friend(FriendUserID)
type Friend struct {
	OwnerUserID    string    `gorm:"column:owner_user_id;primary_key;size:64"`
	FriendUserID   string    `gorm:"column:friend_user_id;primary_key;size:64"`
	Remark         string    `gorm:"column:remark;size:255"`
	CreateTime     time.Time `gorm:"column:create_time"`
	AddSource      int32     `gorm:"column:add_source"`
	OperatorUserID string    `gorm:"column:operator_user_id;size:64"`
	Ex             string    `gorm:"column:ex;size:1024"`
}

// message FriendRequest{
// string  FromUserID = 1;
// string ToUserID = 2;
// int32 HandleResult = 3;
// string ReqMsg = 4;
// int64 CreateTime = 5;
// string HandlerUserID = 6;
// string HandleMsg = 7;
// int64 HandleTime = 8;
// string Ex = 9;
// }
// crazy_server_sdk.FriendRequest(nickname, farce url ...) != imdb.FriendRequest
type FriendRequest struct {
	FromUserID    string    `gorm:"column:from_user_id;primary_key;size:64"`
	ToUserID      string    `gorm:"column:to_user_id;primary_key;size:64"`
	HandleResult  int32     `gorm:"column:handle_result"`
	ReqMsg        string    `gorm:"column:req_msg;size:255"`
	CreateTime    time.Time `gorm:"column:create_time"`
	HandlerUserID string    `gorm:"column:handler_user_id;size:64"`
	HandleMsg     string    `gorm:"column:handle_msg;size:255"`
	HandleTime    time.Time `gorm:"column:handle_time"`
	Ex            string    `gorm:"column:ex;size:1024"`
}

func (FriendRequest) TableName() string {
	return "friend_requests"
}

//	message GroupInfo{
//	 string GroupID = 1;
//	 string GroupName = 2;
//	 string Notification = 3;
//	 string Introduction = 4;
//	 string FaceUrl = 5;
//	 string OwnerUserID = 6;
//	 uint32 MemberCount = 8;
//	 int64 CreateTime = 7;
//	 string Ex = 9;
//	 int32 Status = 10;
//	 string CreatorUserID = 11;
//	 int32 GroupType = 12;
//	}
//
// crazy_server_sdk.GroupInfo (OwnerUserID ,  MemberCount )> imdb.Group

//
//CREATE TABLE `groups` (
//`group_id` varchar(64) NOT NULL,
//`ban_click_packet` tinyint(1) NOT NULL COMMENT ' 1为禁止抢红包',
//`name` varchar(255) DEFAULT NULL,
//`notification` varchar(255) DEFAULT NULL,
//`introduction` varchar(255) DEFAULT NULL,
//`face_url` varchar(255) DEFAULT NULL,
//`create_time` datetime(3) DEFAULT NULL,
//`ex` longtext,
//`status` int(11) DEFAULT NULL,
//`creator_user_id` varchar(64) DEFAULT NULL,
//`group_type` int(11) DEFAULT NULL,
//`need_verification` int(11) DEFAULT NULL,
//`look_member_info` int(11) DEFAULT NULL,
//`apply_member_friend` int(11) DEFAULT NULL,
//`notification_update_time` datetime(3) DEFAULT NULL,
//`notification_user_id` varchar(64) DEFAULT NULL,
//PRIMARY KEY (`group_id`),
//KEY `create_time` (`create_time`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8;

type Group struct {
	//`json:"operationID" binding:"required"`
	//`protobuf:"bytes,1,opt,name=GroupID" json:"GroupID,omitempty"` `json:"operationID" binding:"required"`
	GroupID                string    `gorm:"column:group_id;primary_key;size:64" json:"groupID" binding:"required"`
	BanClickPacket         int32     `gorm:"column:ban_click_packet" json:"banClickPacket"`
	GroupName              string    `gorm:"column:name;size:255" json:"groupName"`
	Notification           string    `gorm:"column:notification;size:255" json:"notification"`
	Introduction           string    `gorm:"column:introduction;size:255" json:"introduction"`
	FaceURL                string    `gorm:"column:face_url;size:255" json:"faceURL"`
	CreateTime             time.Time `gorm:"column:create_time;index:create_time"`
	Ex                     string    `gorm:"column:ex" json:"ex;size:1024" json:"ex"`
	Status                 int32     `gorm:"column:status"`
	CreatorUserID          string    `gorm:"column:creator_user_id;size:64"`
	GroupType              int32     `gorm:"column:group_type"`
	NeedVerification       int32     `gorm:"column:need_verification"`
	LookMemberInfo         int32     `gorm:"column:look_member_info" json:"lookMemberInfo"`
	ApplyMemberFriend      int32     `gorm:"column:apply_member_friend" json:"applyMemberFriend"`
	NotificationUpdateTime time.Time `gorm:"column:notification_update_time"`
	NotificationUserID     string    `gorm:"column:notification_user_id;size:64"`
}

// message GroupMemberFullInfo {
// string GroupID = 1 ;
// string UserID = 2 ;
// int32 roleLevel = 3;
// int64 JoinTime = 4;
// string NickName = 5;
// string FaceUrl = 6;
// int32 JoinSource = 8;
// string OperatorUserID = 9;
// string Ex = 10;
// int32 AppMangerLevel = 7; //if >0
// }  crazy_server_sdk.GroupMemberFullInfo(AppMangerLevel) > imdb.GroupMember
type GroupMember struct {
	GroupID        string    `gorm:"column:group_id;primary_key;size:64"`
	UserID         string    `gorm:"column:user_id;primary_key;size:64"`
	Nickname       string    `gorm:"column:nickname;size:255"`
	FaceURL        string    `gorm:"column:user_group_face_url;size:255"`
	RoleLevel      int32     `gorm:"column:role_level"`
	JoinTime       time.Time `gorm:"column:join_time"`
	JoinSource     int32     `gorm:"column:join_source"`
	InviterUserID  string    `gorm:"column:inviter_user_id;size:64"`
	OperatorUserID string    `gorm:"column:operator_user_id;size:64"`
	MuteEndTime    time.Time `gorm:"column:mute_end_time"`
	Ex             string    `gorm:"column:ex;size:1024"`
}

// message GroupRequest{
// string UserID = 1;
// string GroupID = 2;
// string HandleResult = 3;
// string ReqMsg = 4;
// string  HandleMsg = 5;
// int64 ReqTime = 6;
// string HandleUserID = 7;
// int64 HandleTime = 8;
// string Ex = 9;
// }crazy_server_sdk.GroupRequest == imdb.GroupRequest
type GroupRequest struct {
	UserID        string    `gorm:"column:user_id;primary_key;size:64"`
	GroupID       string    `gorm:"column:group_id;primary_key;size:64"`
	HandleResult  int32     `gorm:"column:handle_result"`
	ReqMsg        string    `gorm:"column:req_msg;size:1024"`
	HandledMsg    string    `gorm:"column:handle_msg;size:1024"`
	ReqTime       time.Time `gorm:"column:req_time"`
	HandleUserID  string    `gorm:"column:handle_user_id;size:64"`
	HandledTime   time.Time `gorm:"column:handle_time"`
	JoinSource    int32     `gorm:"column:join_source"`
	InviterUserID string    `gorm:"column:inviter_user_id;size:64"`
	Ex            string    `gorm:"column:ex;size:1024"`
}

// string UserID = 1;
// string Nickname = 2;
// string FaceUrl = 3;
// int32 Gender = 4;
// string PhoneNumber = 5;
// string Birth = 6;
// string Email = 7;
// string Ex = 8;
// string CreateIp = 9;
// int64 CreateTime = 10;
// int32 AppMangerLevel = 11;
// crazy_server_sdk.User == imdb.User
type User struct {
	UserID           string    `gorm:"column:user_id;primary_key;size:64"`
	Nickname         string    `gorm:"column:name;size:255"`
	Unionid          string    `gorm:"column:unionid"`
	FaceURL          string    `gorm:"column:face_url;size:255"`
	Gender           int32     `gorm:"column:gender"`
	PhoneNumber      string    `gorm:"column:phone_number;size:32"`
	Birth            time.Time `gorm:"column:birth"`
	Email            string    `gorm:"column:email;size:64"`
	Ex               string    `gorm:"column:ex;size:1024"`
	status           int32     `gorm:"column:status"`
	AppMangerLevel   int32     `gorm:"column:app_manger_level"`
	GlobalRecvMsgOpt int32     `gorm:"column:global_recv_msg_opt"`
	CreateTime       time.Time `gorm:"column:create_time;index:create_time"`
}

type UserIpRecord struct {
	UserID        string    `gorm:"column:user_id;primary_key;size:64"`
	CreateIp      string    `gorm:"column:create_ip;size:15"`
	LastLoginTime time.Time `gorm:"column:last_login_time"`
	LastLoginIp   string    `gorm:"column:last_login_ip;size:15"`
	LoginTimes    int32     `gorm:"column:login_times"`
}

// ip limit login
type IpLimit struct {
	Ip            string    `gorm:"column:ip;primary_key;size:15"`
	LimitRegister int32     `gorm:"column:limit_register;size:1"`
	LimitLogin    int32     `gorm:"column:limit_login;size:1"`
	CreateTime    time.Time `gorm:"column:create_time"`
	LimitTime     time.Time `gorm:"column:limit_time"`
}

// ip login
type UserIpLimit struct {
	UserID     string    `gorm:"column:user_id;primary_key;size:64"`
	Ip         string    `gorm:"column:ip;primary_key;size:15"`
	CreateTime time.Time `gorm:"column:create_time"`
}

// message BlackInfo{
// string OwnerUserID = 1;
// int64 CreateTime = 2;
// PublicUserInfo BlackUserInfo = 4;
// int32 AddSource = 5;
// string OperatorUserID = 6;
// string Ex = 7;
// }
// crazy_server_sdk.BlackInfo(BlackUserInfo) != imdb.Black (BlockUserID)
type Black struct {
	OwnerUserID    string    `gorm:"column:owner_user_id;primary_key;size:64"`
	BlockUserID    string    `gorm:"column:block_user_id;primary_key;size:64"`
	CreateTime     time.Time `gorm:"column:create_time"`
	AddSource      int32     `gorm:"column:add_source"`
	OperatorUserID string    `gorm:"column:operator_user_id;size:64"`
	Ex             string    `gorm:"column:ex;size:1024"`
}

type ChatLog struct {
	ServerMsgID      string    `gorm:"column:server_msg_id;primary_key;type:char(64)" json:"serverMsgID"`
	ClientMsgID      string    `gorm:"column:client_msg_id;type:char(64)" json:"clientMsgID"`
	SendID           string    `gorm:"column:send_id;type:char(64);index:send_id,priority:2" json:"sendID"`
	RecvID           string    `gorm:"column:recv_id;type:char(64);index:recv_id,priority:2" json:"recvID"`
	SenderPlatformID int32     `gorm:"column:sender_platform_id" json:"senderPlatformID"`
	SenderNickname   string    `gorm:"column:sender_nick_name;type:varchar(255)" json:"senderNickname"`
	SenderFaceURL    string    `gorm:"column:sender_face_url;type:varchar(255);" json:"senderFaceURL"`
	SessionType      int32     `gorm:"column:session_type;index:session_type,priority:2;index:session_type_alone" json:"sessionType"`
	MsgFrom          int32     `gorm:"column:msg_from" json:"msgFrom"`
	ContentType      int32     `gorm:"column:content_type;index:content_type,priority:2;index:content_type_alone" json:"contentType"`
	Content          string    `gorm:"column:content;type:varchar(3000)" json:"content"`
	Status           int32     `gorm:"column:status" json:"status"`
	SendTime         time.Time `gorm:"column:send_time;index:sendTime;index:content_type,priority:1;index:session_type,priority:1;index:recv_id,priority:1;index:send_id,priority:1" json:"sendTime"`
	CreateTime       time.Time `gorm:"column:create_time" json:"createTime"`
	Ex               string    `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}

func (ChatLog) TableName() string {
	return "chat_logs"
}

type BlackList struct {
	UserId           string    `gorm:"column:uid"`
	BeginDisableTime time.Time `gorm:"column:begin_disable_time"`
	EndDisableTime   time.Time `gorm:"column:end_disable_time"`
}
type Conversation struct {
	OwnerUserID           string `gorm:"column:owner_user_id;primary_key;type:char(128)" json:"OwnerUserID"`
	ConversationID        string `gorm:"column:conversation_id;primary_key;type:char(128)" json:"conversationID"`
	ConversationType      int32  `gorm:"column:conversation_type" json:"conversationType"`
	UserID                string `gorm:"column:user_id;type:char(64)" json:"userID"`
	GroupID               string `gorm:"column:group_id;type:char(128)" json:"groupID"`
	RecvMsgOpt            int32  `gorm:"column:recv_msg_opt" json:"recvMsgOpt"`
	UnreadCount           int32  `gorm:"column:unread_count" json:"unreadCount"`
	DraftTextTime         int64  `gorm:"column:draft_text_time" json:"draftTextTime"`
	IsPinned              bool   `gorm:"column:is_pinned" json:"isPinned"`
	IsPrivateChat         bool   `gorm:"column:is_private_chat" json:"isPrivateChat"`
	BurnDuration          int32  `gorm:"column:burn_duration;default:30" json:"burnDuration"`
	GroupAtType           int32  `gorm:"column:group_at_type" json:"groupAtType"`
	IsNotInGroup          bool   `gorm:"column:is_not_in_group" json:"isNotInGroup"`
	UpdateUnreadCountTime int64  `gorm:"column:update_unread_count_time" json:"updateUnreadCountTime"`
	AttachedInfo          string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	Ex                    string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}

func (Conversation) TableName() string {
	return "conversations"
}

type Department struct {
	DepartmentID   string    `gorm:"column:department_id;primary_key;size:64" json:"departmentID"`
	FaceURL        string    `gorm:"column:face_url;size:255" json:"faceURL"`
	Name           string    `gorm:"column:name;size:256" json:"name" binding:"required"`
	ParentID       string    `gorm:"column:parent_id;size:64" json:"parentID" binding:"required"` // "0" or Real parent id
	Order          int32     `gorm:"column:order" json:"order" `                                  // 1, 2, ...
	DepartmentType int32     `gorm:"column:department_type" json:"departmentType"`                //1, 2...
	RelatedGroupID string    `gorm:"column:related_group_id;size:64" json:"relatedGroupID"`
	CreateTime     time.Time `gorm:"column:create_time" json:"createTime"`
	Ex             string    `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}

func (Department) TableName() string {
	return "departments"
}

type OrganizationUser struct {
	UserID      string    `gorm:"column:user_id;primary_key;size:64"`
	Nickname    string    `gorm:"column:nickname;size:256"`
	EnglishName string    `gorm:"column:english_name;size:256"`
	FaceURL     string    `gorm:"column:face_url;size:256"`
	Gender      int32     `gorm:"column:gender"` //1 ,2
	Mobile      string    `gorm:"column:mobile;size:32"`
	Telephone   string    `gorm:"column:telephone;size:32"`
	Birth       time.Time `gorm:"column:birth"`
	Email       string    `gorm:"column:email;size:64"`
	CreateTime  time.Time `gorm:"column:create_time"`
	Ex          string    `gorm:"column:ex;size:1024"`
}

func (OrganizationUser) TableName() string {
	return "organization_users"
}

type DepartmentMember struct {
	UserID       string    `gorm:"column:user_id;primary_key;size:64"`
	DepartmentID string    `gorm:"column:department_id;primary_key;size:64"`
	Order        int32     `gorm:"column:order" json:"order"` //1,2
	Position     string    `gorm:"column:position;size:256" json:"position"`
	Leader       int32     `gorm:"column:leader" json:"leader"` //-1, 1
	Status       int32     `gorm:"column:status" json:"status"` //-1, 1
	CreateTime   time.Time `gorm:"column:create_time"`
	Ex           string    `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}

func (DepartmentMember) TableName() string {
	return "department_members"
}

type AppVersion struct {
	Version     string `gorm:"column:version;size:64" json:"version"`
	Type        int    `gorm:"column:type;primary_key" json:"type"`
	UpdateTime  int    `gorm:"column:update_time" json:"update_time"`
	ForceUpdate bool   `gorm:"column:force_update" json:"force_update"`
	FileName    string `gorm:"column:file_name" json:"file_name"`
	YamlName    string `gorm:"column:yaml_name" json:"yaml_name"`
	UpdateLog   string `gorm:"column:update_log" json:"update_log"`
}

func (AppVersion) TableName() string {
	return "app_version"
}

type RegisterAddFriend struct {
	UserID string `gorm:"column:user_id;primary_key;size:64"`
}

func (RegisterAddFriend) TableName() string {
	return "register_add_friend"
}

type ClientInitConfig struct {
	DiscoverPageURL string `gorm:"column:discover_page_url;size:64" json:"version"`
}

func (ClientInitConfig) TableName() string {
	return "client_init_config"
}

type FNcountAccount struct {
	Id              int32     `gorm:"column:id;type:int(10) unsigned;not null;primary_key;auto_increment;comment:'主键'" json:"id"`
	UserID          string    `gorm:"column:user_id;type:varchar(64);not null;comment:'用户id'" json:"userID"`
	MainAccountId   string    `gorm:"column:main_account_id;type:varchar(20);default:null;comment:'新生支付主账号id'" json:"mainAccountId"`
	PacketAccountId string    `gorm:"column:packet_account_id;type:varchar(20);default:null;comment:'新生支付红包账户id'" json:"packetAccountId"`
	Mobile          string    `gorm:"column:mobile;type:varchar(15);not null;comment:'手机号码'" json:"mobile"`
	RealName        string    `gorm:"column:realname;type:varchar(20);not null;comment:'真实姓名'" json:"realName"`
	IdCard          string    `gorm:"column:id_card;type:varchar(30);not null;comment:'身份证'" json:"idCard"`
	PaySwitch       int32     `gorm:"column:pay_switch;type:tinyint(4);default:1;comment:'支付开关(0关闭、1默认开启)'" json:"paySwitch"`
	BodPaySwitch    int32     `gorm:"column:bod_pay_switch;type:tinyint(4);default:0;comment:'指纹支付/人脸支付开关(0默认关闭、1开启)'" json:"bodPaySwitch"`
	PaymentPassword string    `gorm:"column:payment_password;type:varchar(32);default:null;comment:'支付密码(md5加密)'" json:"paymentPassword"`
	OpenStatus      int32     `gorm:"column:open_status;type:tinyint(4);default:0;comment:'开通状态'" json:"openStatus"`
	OpenStep        int32     `gorm:"column:open_step;type:tinyint(4);default:1;comment:'开通认证步骤(1身份证认证、2支付密码、3绑定银行卡或者刷脸)'" json:"openStep"`
	CreatedTime     time.Time `gorm:"column:created_time;type:datetime;default:null" json:"createdTime"`
	UpdatedTime     time.Time `gorm:"column:updated_time;type:datetime;default:null" json:"updatedTime"`
}

func (FNcountAccount) TableName() string {
	return "f_ncount_account"
}

/*  `ncount_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '账户类型(1主账户，2红包账户)',*/
// 用户银行卡绑定表
type FNcountBankCard struct {
	Id                int32     `gorm:"column:id" json:"id"`
	UserId            string    `gorm:"column:user_id" json:"user_id"`                         //用户id
	NcountUserId      string    `gorm:"column:ncount_user_id" json:"ncount_user_id"`           //新生支付用户id
	MerOrderId        string    `gorm:"column:mer_order_id" json:"mer_order_id"`               //平台订单号
	NcountOrderId     string    `gorm:"column:ncount_order_id" json:"ncount_order_id"`         //第三方签约订单号
	BindCardAgrNo     string    `gorm:"column:bind_card_agr_no" json:"bind_card_agr_no"`       //第三方绑卡协议号
	Mobile            string    `gorm:"column:mobile" json:"mobile"`                           //手机号码
	CardOwner         string    `gorm:"column:card_owner" json:"card_owner"`                   //持卡者名字
	BankCardNumber    string    `gorm:"column:bank_card_number" json:"bank_card_number"`       //银行卡号
	CardAvailableDate string    `gorm:"column:card_available_date" json:"card_available_date"` //信用卡有效期
	Cvv2              string    `gorm:"column:cvv2" json:"cvv2"`                               //cvv2
	BankCode          string    `gorm:"column:bank_code" json:"bank_code"`                     //银行简写
	IsBind            int       `gorm:"column:is_bind" json:"is_bind"`                         //是否绑定成功(0预提交、1绑定成功)
	IsDelete          int       `gorm:"column:is_delete" json:"is_delete"`                     //是否删除(0未删除，1已删除)
	CreatedTime       time.Time `gorm:"column:created_time" json:"created_time"`               //
	UpdatedTime       time.Time `gorm:"column:updated_time" json:"updated_time"`               //
}

func (FNcountBankCard) TableName() string {
	return "f_ncount_bank_card"
}

type FNcountTrade struct {
	ID           int32  `gorm:"column:id;primary_key;AUTO_INCREMENT;not null" json:"id"`
	UserID       string `gorm:"column:user_id;not null" json:"user_id"`    //用户id
	Type         int32  `gorm:"column:type" json:"type"`                   //收支类型(1收入、2支出)
	BusinessType int32  `gorm:"column:business_type" json:"business_type"` //业务类型(见枚举)
	Describe     string `gorm:"column:describe" json:"describe"`           //描述
	Amount       int32  `gorm:"column:amount;not null" json:"amount"`      //变更金额(单位：分)
	//BeferAmount  int32     `gorm:"column:befer_amount" json:"befer_amount"`     //变更前金额(单位：分)
	AfterAmount  int32     `gorm:"column:after_amount" json:"after_amount"`     //变更后金额(单位：分)
	MerOrderId   string    `gorm:"column:mer_order_id" json:"mer_order_id"`     //平台订单号
	ThirdOrderNo string    `gorm:"column:third_order_no" json:"third_order_no"` //第三方订单号
	NcountStatus int32     `gorm:"column:ncount_status" json:"ncount_status"`   //异步通知状态（0未生效，1生效）
	PacketID     string    `gorm:"column:packet_id" json:"packet_id"`           //红包id
	CreatedTime  time.Time `gorm:"column:created_time" json:"created_time"`
	UpdatedTime  time.Time `gorm:"column:updated_time" json:"updated_time"`
}

func (FNcountTrade) TableName() string {
	return "f_ncount_trade"
}

type GroupHistoryMembers struct {
	Id              int       `gorm:"column:id" json:"id"`                                 //id
	GroupId         string    `gorm:"column:group_id" json:"group_id"`                     //群id
	UserId          string    `gorm:"column:user_id" json:"user_id"`                       //用户id
	LastSendMsgTime int       `gorm:"column:last_send_msg_time" json:"last_send_msg_time"` //最后发送群消息时间
	CreatedTime     time.Time `gorm:"column:created_time" json:"created_time"`             //加群时间
}

func (GroupHistoryMembers) TableName() string {
	return "group_history_members"
}

type UserCollect struct {
	Id             int32     `gorm:"column:id" json:"id"`                           //id
	UserId         string    `gorm:"column:user_id" json:"user_id"`                 //用户id
	CollectType    int32     `gorm:"column:collect_type" json:"collect_type"`       //收藏类型
	CollectContent string    `gorm:"column:collect_content" json:"collect_content"` //收藏内容
	CreatedTime    time.Time `gorm:"column:created_time" json:"created_time"`       //收藏时间
}

func (UserCollect) TableName() string {
	return "user_collect"
}

//CREATE TABLE `f_error_log` (
//`id` int(11) NOT NULL AUTO_INCREMENT,
//`remark` varchar(255) DEFAULT NULL COMMENT '我们的备注，',
//`oprationID` varchar(255) DEFAULT NULL COMMENT '平台唯一操作ID',
//`mer_order_id` varchar(255) DEFAULT NULL COMMENT '新生支付那边的订单',
//`err_msg` varchar(255) DEFAULT NULL COMMENT '新生支付的错误',
//`err_code` varchar(255) DEFAULT NULL COMMENT '错误码',
//`all_msg` text COMMENT '完整信息集合',
//`create_time` int(11) DEFAULT NULL COMMENT '创建时间',
//PRIMARY KEY (`id`)
//) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

type FErrorLog struct {
	Id         int32  `gorm:"column:id" json:"id"`                     //id
	Remark     string `gorm:"column:remark" json:"remark"`             //我们的备注，
	OprationID string `gorm:"column:oprationID" json:"oprationID"`     //平台唯一操作ID
	MerOrderId string `gorm:"column:mer_order_id" json:"mer_order_id"` //新生支付那边的订单
	ErrMsg     string `gorm:"column:err_msg" json:"err_msg"`           //新生支付的错误
	ErrCode    string `gorm:"column:err_code" json:"err_code"`         //错误码
	AllMsg     string `gorm:"column:all_msg" json:"all_msg"`           //完整信息集合
	CreateTime int64  `gorm:"column:create_time" json:"create_time"`   //创建时间
}

func (FErrorLog) TableName() string {
	return "f_error_log"
}

//CREATE TABLE `f_packet` (
//`id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
//`packet_id` varchar(255) DEFAULT NULL COMMENT '红包ID',
//`submit_time` varchar(255) DEFAULT NULL COMMENT '下单时间，用于退款',
//`remark` varchar(255) DEFAULT NULL COMMENT '红包状态描述',
//`status` int(2) NOT NULL COMMENT '红包状态： 1 为创建 、2 为正常、3为异常  ,100 为退回，200 为退回异常',
//`user_id` varchar(255) NOT NULL COMMENT '红包发起者',
//`user_redpacket_account` varchar(255) DEFAULT NULL COMMENT '发送红包的用户的账户',
//`packet_type` tinyint(1) NOT NULL COMMENT '红包类型(1个人红包、2群红包)',
//`is_lucky` tinyint(1) DEFAULT '0' COMMENT '是否为拼手气红包',
//`is_exclusive` tinyint(1) NOT NULL COMMENT '是否为专属红包： 0为否，1为是',
//`exclusive_user_id` varchar(255) DEFAULT '0' COMMENT '专属用户id',
//`packet_title` varchar(100) NOT NULL COMMENT '红包标题',
//`amount` int(11) NOT NULL COMMENT '单个红包金额，如果说是',
//`number` tinyint(3) NOT NULL COMMENT '红包个数',
//`total_amount` int(11) DEFAULT NULL COMMENT '发红包的总金额 == remain_amount的初始值',
//`expire_time` int(11) DEFAULT NULL COMMENT '红包过期时间',
//`mer_order_id` varchar(255) DEFAULT NULL COMMENT '红包第三方的请求ID',
//`ncount_order_id` varchar(255) DEFAULT NULL COMMENT '新生支付id',
//`operate_id` varchar(255) DEFAULT NULL COMMENT '链路追踪ID',
//`recv_id` varchar(255) DEFAULT NULL COMMENT '被发送用户的ID',
//`send_type` tinyint(11) DEFAULT NULL COMMENT '红包发送方式： 1：钱包余额，2是银行卡',
//`bind_card_agr_no` varchar(255) DEFAULT NULL COMMENT '银行卡绑定协议号',
//`remain` int(11) DEFAULT NULL COMMENT '剩余红包数量',
//`remain_amout` int(11) NOT NULL DEFAULT '0' COMMENT '剩余红包金额',
//`refound_amout` int(11) DEFAULT NULL COMMENT '退款金额',
//`lucky_user_id` varchar(255) NOT NULL DEFAULT '' COMMENT '最佳手气红包用户ID',
//`luck_user_amount` int(11) NOT NULL DEFAULT '0' COMMENT '最大红包的值： account amount  分为单位',
//`created_time` int(11) DEFAULT NULL,
//`updated_time` int(11) DEFAULT NULL,
//PRIMARY KEY (`id`) USING BTREE,
//KEY `idx_user_id` (`user_id`) USING BTREE,
//KEY `idx_packet_id` (`packet_id`) USING BTREE,
//KEY `idx_expire` (`expire_time`) USING BTREE
//) ENGINE=InnoDB AUTO_INCREMENT=372 DEFAULT CHARSET=utf8mb4 COMMENT='用户红包表';
type FPacket struct {
	ID                   int64  `gorm:"column:id;primary_key;AUTO_INCREMENT;not null" json:"id"`
	PacketID             string `gorm:"column:packet_id;not null" json:"packet_id"`
	SubmitTime           string `gorm:"column:submit_time;not null" json:"submit_time"` // 调用新生支付的时间
	UserID               string `gorm:"column:user_id;not null" json:"user_id"`
	UserRedpacketAccount string `gorm:"column:user_redpacket_account;not null" json:"user_redpacket_account"`
	PacketType           int32  `gorm:"column:packet_type;not null" json:"packet_type"`
	IsLucky              int32  `gorm:"column:is_lucky;not null" json:"is_lucky"`
	ExclusiveUserID      string `gorm:"column:exclusive_user_id;not null" json:"exclusive_user_id"`
	PacketTitle          string `gorm:"column:packet_title;not null" json:"packet_title"`
	Amount               int64  `gorm:"column:amount;not null" json:"amount"`
	Number               int32  `gorm:"column:number;not null" json:"number"`
	TotalAmount          int64  `gorm:"column:total_amount;not null" json:"total_amount"`
	ExpireTime           int64  `gorm:"column:expire_time;not null" json:"expire_time"`
	MerOrderID           string `gorm:"column:mer_order_id;not null" json:"mer_order_id"`
	NcountOrderID        string `gorm:"column:ncount_order_id;not null" json:"ncount_order_id"`
	SendType             int32  `gorm:"column:send_type;not null" json:"send_type"`
	BindCardAgrNo        string `gorm:"column:bind_card_agr_no;not null" json:"bind_card_agr_no"`
	OperateID            string `gorm:"column:operate_id;not null" json:"operate_id"`
	RecvID               string `gorm:"column:recv_id;not null" json:"recv_id"`
	Remain               int64  `gorm:"column:remain;not null" json:"remain"`                     // 剩余红包数量
	RemainAmout          int64  `gorm:"column:remain_amout;not null" json:"remain_amout"`         // 剩余红包金额
	RefoundAmout         int64  `gorm:"column:refound_amout;not null" json:"refound_amout"`       // 退款金额
	LuckyUserID          string `gorm:"column:lucky_user_id;not null" json:"lucky_user_id"`       // 最佳手气红包用户ID
	LuckUserAmount       int64  `gorm:"column:luck_user_amount;not null" json:"luck_user_amount"` // 最大红包的值： account amount  分为单位
	CreatedTime          int64  `gorm:"column:created_time;not null" json:"created_time"`
	UpdatedTime          int64  `gorm:"column:updated_time;not null" json:"updated_time"`
	Status               int32  `gorm:"column:status;not null" json:"status"` // 0 创建未生效，1 为红包正在领取中，2为红包领取完毕，3为红包过期
	Remark               string `gorm:"column:remark;not null" json:"remark"`
	IsExclusive          int32  `gorm:"column:is_exclusive;not null" json:"is_exclusive"`
}

func (FPacket) TableName() string {
	return "f_packet"
}

type FVersion struct {
	ID            int64  `gorm:"column:id;primary_key;AUTO_INCREMENT;not null" json:"id"`
	AppType       int32  `gorm:"column:app_type" json:"app_type"` //app类型(1、安卓，2、ios)
	VersionCode   string `gorm:"column:version_code;not null" json:"version_code"`
	DownloadUrl   string `gorm:"column:download_url;not null" json:"download_url"`
	UpdateContent string `gorm:"column:update_content;not null" json:"update_content"`
	IsForce       int32  `gorm:"column:is_force;not null" json:"is_force"`
	Status        int32  `gorm:"column:status" json:"status"` //app类型(1、安卓，2、ios)
	CreateTime    int64  `gorm:"column:create_time;not null" json:"create_time"`
}

func (FVersion) TableName() string {
	return "f_version"
}

type UserAttributeSwitch struct {
	Id                    int32     `gorm:"column:id" json:"id"`                                             //id
	UserId                string    `gorm:"column:user_id" json:"user_id"`                                   //用户id
	AddFriendVerifySwitch int32     `gorm:"column:add_friend_verify_switch" json:"add_friend_verify_switch"` //加好友验证开关
	AddFriendGroupSwitch  int32     `gorm:"column:add_friend_group_switch" json:"add_friend_group_switch"`   //加好友群组开关
	AddFriendQrcodeSwitch int32     `gorm:"column:add_friend_qrcode_switch" json:"add_friend_qrcode_switch"` //加好友二维码开关
	AddFriendCardSwitch   int32     `gorm:"column:add_friend_card_switch" json:"add_friend_card_switch"`     //加好友名片开关
	UpdatedTime           time.Time `gorm:"column:updated_time" json:"updated_time"`                         //时间
}

func (UserAttributeSwitch) TableName() string {
	return "user_attribute_switch"
}

/*
CREATE TABLE `third_pay_order` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `order_no` varchar(255) DEFAULT NULL COMMENT '订单ID，用于提供支付的',
  `mer_order_no` varchar(255) DEFAULT NULL COMMENT '商户订单ID',
  `ncount_order_no` varchar(255) DEFAULT NULL COMMENT '调用新生支付的mer_order_id',
  `ncount_ture_order_no` varchar(255) DEFAULT NULL COMMENT '新生支付的orderID',
  `order_type` int(11) DEFAULT '100' COMMENT '100 是支付 ；200是提现 ；',
  `mer_id` varchar(255) DEFAULT NULL COMMENT '商户ID',
  `amount` int(11) DEFAULT NULL COMMENT '订单金额',
  `recieve_account` varchar(255) DEFAULT NULL COMMENT '收款方的新生支付的acount',
  `status` int(5) DEFAULT NULL COMMENT '100： 创建， 200:支付成功，400是支付失败',
  `remark` varchar(255) DEFAULT NULL COMMENT '支付备注',
  `notify_url` varchar(255) DEFAULT NULL COMMENT '回调地址',
  `is_notify` varchar(255) DEFAULT NULL COMMENT '是否对用户提供的地址进行回调',
  `notify_count` varchar(255) DEFAULT NULL COMMENT '回调的总次数',
  `last_notify_time` datetime DEFAULT NULL COMMENT '最后一次回调时间',
  `pay_time` datetime DEFAULT NULL COMMENT '订单支付时间',
  `add_time` datetime DEFAULT NULL,
  `edit_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=79 DEFAULT CHARSET=utf8;
*/

type ThirdPayOrder struct {
	Id             int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;not null" json:"id"`
	OrderNo        string    `gorm:"column:order_no;not null" json:"order_no"`
	MerOrderNo     string    `gorm:"column:mer_order_no;not null" json:"mer_order_no"`
	NcountOrderNo  string    `gorm:"column:ncount_order_no;not null" json:"ncount_order_no"`
	NcountTureNo   string    `gorm:"column:ncount_ture_order_no;not null" json:"ncount_ture_order_no"`
	OderType       int32     `gorm:"column:order_type;not null" json:"order_type"`
	MerId          string    `gorm:"column:mer_id;not null" json:"mer_id"`
	Status         int32     `gorm:"column:status;not null" json:"status" `
	Amount         int64     `gorm:"column:amount;not null" json:"amount"`
	RecieveAccount string    `gorm:"column:recieve_account;not null" json:"recieve_account"`
	Remark         string    `gorm:"column:remark;not null" json:"remark"`
	NotifyUrl      string    `gorm:"column:notify_url;not null" json:"notify_url"`
	IsNotify       int32     `gorm:"column:is_notify;not null" json:"is_notify"`
	NotifyCount    int32     `gorm:"column:notify_count;not null" json:"notify_count"`
	LastNotifyTime time.Time `gorm:"column:last_notify_time;not null" json:"last_notify_time"`
	PayTime        time.Time `gorm:"column:pay_time;not null" json:"pay_time"`
	AddTime        time.Time `gorm:"column:add_time;not null" json:"add_time"`
	EditTime       time.Time `gorm:"column:edit_time;not null" json:"edit_time"`
}

func (ThirdPayOrder) TableName() string {
	return "third_pay_order"
}

/*CREATE TABLE `third_pay_merchant` (
`id` int(11) NOT NULL,
`merchant_id` varchar(255) DEFAULT NULL COMMENT '商户号',
`name` varchar(255) DEFAULT NULL COMMENT '商户名称',
`ncount_account` varchar(255) DEFAULT NULL COMMENT '新生支付的账号',
`add_time` datetime DEFAULT NULL COMMENT '创建时间',
`edit_time` datetime DEFAULT NULL COMMENT '最后修改时间',
PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;*/

type ThirdPayMerchant struct {
	Id            int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;not null" json:"id"`
	MerchantId    string    `gorm:"column:merchant_id;not null" json:"merchant_id"`
	Name          string    `gorm:"column:name;not null" json:"name"`
	NcountAccount string    `gorm:"column:ncount_account;not null" json:"ncount_account"`
	AddTime       time.Time `gorm:"column:add_time;not null" json:"add_time"`
	EditTime      time.Time `gorm:"column:edit_time;not null" json:"edit_time"`
}

func (ThirdPayMerchant) TableName() string {
	return "third_pay_merchant"
}

type ThirdWithdraw struct {
	Id            int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;not null" json:"id"`
	UserId        string    `gorm:"column:user_id;not null" json:"user_id"`
	MerOrderId    string    `gorm:"column:mer_order_id;not null" json:"mer_order_id"`
	NcountOrderId string    `gorm:"column:ncount_order_id;not null" json:"ncount_order_id"`
	ThirdOrderId  string    `gorm:"column:third_order_id;not null" json:"third_order_id"`
	Account       string    `gorm:"column:account;not null" json:"account"`
	Amount        int64     `gorm:"column:amount;not null" json:"amount"`
	RecevieAmount int64     `gorm:"column:recevie_amount;not null" json:"recevie_amount"`
	Commission    int64     `gorm:"column:commission;not null" json:"commission"`
	Status        int32     `gorm:"column:status;not null" json:"status"`
	Remark        string    `gorm:"column:remark;not null" json:"remark"`
	AddTime       time.Time `gorm:"column:add_time;not null" json:"add_time"`
	UpdateTime    time.Time `gorm:"column:update_time;not null" json:"update_time"`
}

func (ThirdWithdraw) TableName() string {
	return "third_withdraw"
}

type AppWgtVersion struct {
	Id          int32     `json:"id"`
	AppId       string    `json:"app_id"`  // 应用appid
	Version     string    `json:"version"` // 版本号
	Url         string    `json:"url"`     // 应用url地址
	Remarks     string    `json:"remarks"`
	Status      int8      `json:"status"` // 状态
	CreatedTime time.Time `json:"created_time"`
}

func (AppWgtVersion) TableName() string {
	return "app_wgt_version"
}

//CREATE TABLE `help_feedback` (
//  `id` int(11) NOT NULL,
//  `user_id` varchar(255) DEFAULT NULL COMMENT '用户的ID',
//  `type` int(1) DEFAULT NULL COMMENT '1:建议、2是功能问题、3是违法问题',
//  `content` varchar(500) DEFAULT NULL COMMENT '问题描述（最大存储100个字节）',
//  `contact` varchar(255) DEFAULT NULL COMMENT '联系人电话',
//  `status` int(2) DEFAULT NULL COMMENT '1是正常状态，2是关闭状态',
//  `add_time` datetime DEFAULT NULL,
//  `update_time` datetime DEFAULT NULL,
//  PRIMARY KEY (`id`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='帮助反馈表';
type HelpFeedback struct {
	Id         int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;not null" json:"id"`
	UserId     string    `gorm:"column:user_id;not null" json:"user_id"`
	Type       int32     `gorm:"column:type;not null" json:"type"`
	Content    string    `gorm:"column:content;not null" json:"content"`
	Contact    string    `gorm:"column:contact;not null" json:"contact"`
	Status     int32     `gorm:"column:status;not null" json:"status"`
	AddTime    time.Time `gorm:"column:add_time;not null" json:"add_time"`
	UpdateTime time.Time `gorm:"column:update_time;not null" json:"update_time"`
}

//CREATE TABLE `help_normal_question` (
//`id` int(11) NOT NULL AUTO_INCREMENT,
//`title` varchar(255) DEFAULT NULL,
//`content` text,
//`solved` int(11) DEFAULT '0',
//`unsolved` int(11) DEFAULT '0',
//`status` tinyint(1) DEFAULT '1' COMMENT '1是正常，2是删除',
//`ord` int(11) DEFAULT '0' COMMENT '排序',
//`add_time` datetime DEFAULT NULL,
//`update_time` datetime DEFAULT NULL,
//PRIMARY KEY (`id`)
//) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

type HelpNormalQuestion struct {
	Id         int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;not null" json:"id"`
	Title      string    `gorm:"column:title;not null" json:"title"`
	Content    string    `gorm:"column:content;not null" json:"content"`
	Solved     int32     `gorm:"column:solved;not null" json:"solved"`
	Unsolved   int32     `gorm:"column:unsolved;not null" json:"unsolved"`
	Status     int32     `gorm:"column:status;not null" json:"status"`
	Ord        int32     `gorm:"column:ord;not null" json:"ord"`
	AddTime    time.Time `gorm:"column:add_time;not null" json:"add_time"`
	UpdateTime time.Time `gorm:"column:update_time;not null" json:"update_time"`
}

//CREATE TABLE `blockword` (
//`id` int(11) NOT NULL AUTO_INCREMENT,
//`word` varchar(100) DEFAULT NULL COMMENT '敏感字串',
//PRIMARY KEY (`id`)
//) ENGINE=InnoDB AUTO_INCREMENT=6097 DEFAULT CHARSET=utf8;

type Blockword struct {
	Id   int64  `gorm:"column:id;primary_key;AUTO_INCREMENT;not null" json:"id"`
	Word string `gorm:"column:word;not null" json:"word"`
}
