package api

import (
	"crazy_server/internal/api/agent"
	apiAuth "crazy_server/internal/api/auth"
	clientInit "crazy_server/internal/api/client_init"
	"crazy_server/internal/api/cloud_wallet/account"
	"crazy_server/internal/api/cloud_wallet/notify"
	"crazy_server/internal/api/cloud_wallet/redpacket"
	"crazy_server/internal/api/conversation"
	"crazy_server/internal/api/filter"
	"crazy_server/internal/api/friend"
	"crazy_server/internal/api/group"
	"crazy_server/internal/api/manage"
	"crazy_server/internal/api/middleware"
	apiChat "crazy_server/internal/api/msg"
	"crazy_server/internal/api/office"
	"crazy_server/internal/api/organization"
	"crazy_server/internal/api/system"
	apiThird "crazy_server/internal/api/third"
	"crazy_server/internal/api/user"
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/constant"
	"crazy_server/pkg/common/log"
	"crazy_server/pkg/utils"
	"io"
	"os"

	promePkg "crazy_server/pkg/common/prometheus"

	"github.com/gin-gonic/gin"
)

func NewGinRouter() *gin.Engine {
	log.NewPrivateLog(constant.LogFileName)
	gin.SetMode(gin.DebugMode)
	f, _ := os.Create("./logs/api.log")
	gin.DefaultWriter = io.MultiWriter(f)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(utils.CorsHandler())
	log.Info("load config: ", config.Config)
	if config.Config.Prometheus.Enable {
		promePkg.NewApiRequestCounter()
		promePkg.NewApiRequestFailedCounter()
		promePkg.NewApiRequestSuccessCounter()
		r.Use(promePkg.PromeTheusMiddleware)
		r.GET("/metrics", promePkg.PrometheusHandler())
	}

	filterGroup := r.Group("/filter")
	{
		filterGroup.GET("/varify", filter.Filter)
	}

	//推广计划(鉴权)
	agentGroup := r.Group("/agent")
	agentGroup.Use(middleware.JWTAuth())
	{
		//新互娱接口 - start
		agentGroup.POST("user_agent_info", agent.GetUserAgentInfo) //获取当前用户的推广员信息以及绑定关系
		agentGroup.POST("apply", agent.AgentApply)                 //推广员申请提交
		//新互娱接口 - end

		agentGroup.POST("main", agent.AgentMainInfo)                           //推广员主页信息
		agentGroup.POST("account/income_chart", agent.AgentAccountIncomeChart) //账户明细收益趋势图
		agentGroup.POST("account/record_list", agent.AgentAccountRecordList)   //账户明细详情列表

		//咖豆管理
		agentGroup.POST("bean/platform_config", agent.PlatformBeanShopConfig)         //获取平台咖豆商城配置
		agentGroup.POST("bean/config", agent.AgentDiyBeanShopConfig)                  //推广员自定义咖豆商城配置
		agentGroup.POST("bean_account/record_list", agent.AgentBeanAccountRecordList) //咖豆账户明细详情列表
		agentGroup.POST("bean/config_status", agent.AgentBeanShopUpStatus)            //咖豆管理上下架
		agentGroup.POST("bean/config_up", agent.AgentBeanShopUpdate)                  //咖豆配置管理

		agentGroup.POST("member_list", agent.AgentMemberList)          //推广下属用户列表
		agentGroup.POST("give_member_bean", agent.AgentGiveMemberBean) //赠送下属成员咖豆
		agentGroup.POST("purchase_bean", agent.PurchaseBean)           //推广员购买咖豆
		agentGroup.POST("withdraw", agent.Withdraw)                    //推广员余额提现
	}

	//推广系统(不需要鉴权)
	agentCallbackGroup := r.Group("/agent")
	{
		agentCallbackGroup.POST("open", agent.OpenAgent)                                           //推广员开通
		agentCallbackGroup.POST("bind_agent_number", agent.BindAgentNumber)                        //绑定推广员
		agentCallbackGroup.POST("game_shop/purchase_bean", agent.ChessShopPurchaseBean)            //互娱商城购买咖豆下单(预提交)
		agentCallbackGroup.POST("notify/agent_purchase_bean", agent.ChessPurchaseBeanNotify)       //推广员成员购买咖豆回调(推广员商城) - 互娱回调
		agentCallbackGroup.POST("notify/platform_purchase_bean", agent.PlatformPurchaseBeanNotify) //推广员成员购买咖豆回调(平台商城) - 互娱回调
		agentCallbackGroup.POST("notify/recharge", agent.RechargeNotify)                           //推广员充值咖豆 - 新生支付回调
		agentCallbackGroup.POST("game_shop/bean_config", agent.AgentGameShopBeanConfig)            //获取推广员游戏商城咖豆配置
	}

	// CloudWallet
	cloudWalletGroup := r.Group("/cloudWallet")
	{
		// 用户账户管理
		cloudWalletGroup.POST("/account", account.Account)                                //获取账户信息
		cloudWalletGroup.POST("/id_card/real_name/auth", account.IdCardRealNameAuth)      //身份证实名认证
		cloudWalletGroup.POST("/set_payment_secret", account.SetPaymentSecret)            // 设置支付密码
		cloudWalletGroup.POST("/check_payment_secret", account.CheckPaymentSecret)        // 校验支付密码
		cloudWalletGroup.POST("/cloud_wallet/record_list", account.CloudWalletRecordList) // 云钱包明细：云钱包收支情况
		cloudWalletGroup.POST("/cloud_wallet/record_del", account.CloudWalletRecordDel)   // 删除云钱包明细

		//用户银行卡管理
		cloudWalletGroup.POST("/bind_user_bankcard", account.BindUserBankCard)                //绑定银行卡(预提交)
		cloudWalletGroup.POST("/bind_user_bankcard/confirm", account.BindUserBankcardConfirm) //确认绑定银行卡-code验证
		cloudWalletGroup.POST("/unbinding/user_bankcard", account.UnBindUserBankcard)         //解绑银行卡

		//账户
		cloudWalletGroup.POST("/charge_account", account.ChargeAccount)                //账户充值
		cloudWalletGroup.POST("/charge_account/confirm", account.ChargeAccountConfirm) //账户充值code确认
		cloudWalletGroup.POST("/draw_account", account.DrawAccount)                    //提现

		// 回调接口
		cloudWalletGroup.POST("/charge_account_callback", notify.ChargeNotify) // 充值回调
		cloudWalletGroup.POST("/draw_account_callback", notify.WithDrawNotify) //提现回调

		// 红包管理
		cloudWalletGroup.POST("/send_red_packet", redpacket.SendRedPacket)   //发送红包
		cloudWalletGroup.POST("/click_red_packet", redpacket.ClickRedPacket) // 抢红包接口
		// 确认发送红包
		cloudWalletGroup.POST("/send_red_packet/confirm", redpacket.SendRedPacketConfirm) //确认发送红包

		cloudWalletGroup.POST("/red_packet/receive_detail", redpacket.RedPacketReceiveDetail)  // 红包领取明细
		cloudWalletGroup.POST("/red_packet/info", redpacket.GetRedPacketInfo)                  // 红包详情
		cloudWalletGroup.POST("/ban_gourp_click_red_packet", redpacket.BanGroupClickRedPacket) // 禁止群抢红包
		cloudWalletGroup.POST("/refound_packet", redpacket.RefoundPacket)                      // 红包退款，退款到云钱包

		// 生成声网token
		cloudWalletGroup.POST("/getAgoraToken", redpacket.GetAgoraToken)   // 获取声网token
		cloudWalletGroup.POST("/translateVideo", redpacket.TranslateVideo) // 翻译文字
		cloudWalletGroup.POST("/getVersion", redpacket.GetVersion)         // 获取版本

		// 这里是做第三方支付
		cloudWalletGroup.POST("/create_third_pay_order", redpacket.CreateThirdPayOrder) // 创建第三方订单 - 竞技使用
		cloudWalletGroup.POST("/get_third_pay_order", redpacket.GetThirdPayOrder)       // 查询第三方订单
		cloudWalletGroup.POST("/third_pay", redpacket.ThirdPay)                         // 第三方支付
		cloudWalletGroup.POST("/third_withdraw", redpacket.ThirdWithdraw)               // 提现到用户的银行卡

		// 这里做新生支付的统一封装
		cloudWalletGroup.POST("/pay_callback", redpacket.ThirdPayCallback) // 第三方支付
		cloudWalletGroup.POST("/pay_confirm", redpacket.PayConfirm)        // 第三方支付

		// 这里临时给检测使用
		cloudWalletGroup.GET("/check_status", func(context *gin.Context) {
			context.JSON(200, gin.H{
				"code": 0,
				"msg":  "ok",
			})
		}) // 抢红包接口
	}

	// user routing group, which handles user registration and login services
	userRouterGroup := r.Group("/user")
	{
		userRouterGroup.POST("/update_user_info", user.UpdateUserInfo) //1
		userRouterGroup.POST("/set_global_msg_recv_opt", user.SetGlobalRecvMessageOpt)
		userRouterGroup.POST("/get_users_info", user.GetUsersPublicInfo)            //1
		userRouterGroup.POST("/get_self_user_info", user.GetSelfUserInfo)           //1
		userRouterGroup.POST("/get_users_online_status", user.GetUsersOnlineStatus) //1
		userRouterGroup.POST("/get_users_info_from_cache", user.GetUsersInfoFromCache)
		userRouterGroup.POST("/get_user_friend_from_cache", user.GetFriendIDListFromCache)
		userRouterGroup.POST("/get_black_list_from_cache", user.GetBlackIDListFromCache)
		userRouterGroup.POST("/get_all_users_uid", manage.GetAllUsersUid) //1
		userRouterGroup.POST("/account_check", manage.AccountCheck)       //1
		//	userRouterGroup.POST("/get_users_online_status", manage.GetUsersOnlineStatus) //1
		userRouterGroup.POST("/get_users", user.GetUsers)

		userRouterGroup.POST("/attribute_switch", user.AttributeSwitch)        //获取用户属性开关配置
		userRouterGroup.POST("/attribute_switch/set", user.AttributeSwitchSet) //用户属性开关设置
		userRouterGroup.POST("/attribute/menu", user.AttributeMenu)            //用户属性菜单

		// 用户反馈
		userRouterGroup.POST("/feedback", user.Feedback) // 用户反馈
		// 常见问题
		userRouterGroup.POST("/question", user.CommonQuestion) // 常见问题
		// 常见问题反馈
		userRouterGroup.POST("/question/feedback", user.CommonQuestionFeedback) // 常见问题反馈
	}
	//friend routing group
	friendRouterGroup := r.Group("/friend")
	{
		//	friendRouterGroup.POST("/get_friends_info", friend.GetFriendsInfo)
		friendRouterGroup.POST("/add_friend", friend.AddFriend)                              //1
		friendRouterGroup.POST("/delete_friend", friend.DeleteFriend)                        //1
		friendRouterGroup.POST("/get_friend_apply_list", friend.GetFriendApplyList)          //1
		friendRouterGroup.POST("/get_self_friend_apply_list", friend.GetSelfFriendApplyList) //1
		friendRouterGroup.POST("/get_friend_list", friend.GetFriendList)                     //1
		friendRouterGroup.POST("/add_friend_response", friend.AddFriendResponse)             //1
		friendRouterGroup.POST("/set_friend_remark", friend.SetFriendRemark)                 //1

		friendRouterGroup.POST("/add_black", friend.AddBlack)          //1
		friendRouterGroup.POST("/get_black_list", friend.GetBlacklist) //1
		friendRouterGroup.POST("/remove_black", friend.RemoveBlack)    //1

		friendRouterGroup.POST("/import_friend", friend.ImportFriend) //1
		friendRouterGroup.POST("/is_friend", friend.IsFriend)         //1
	}
	//group related routing group
	groupRouterGroup := r.Group("/group")
	{
		groupRouterGroup.POST("/create_group", group.CreateGroup)                                   //1
		groupRouterGroup.POST("/set_group_info", group.SetGroupInfo)                                //1
		groupRouterGroup.POST("/join_group", group.JoinGroup)                                       //1
		groupRouterGroup.POST("/quit_group", group.QuitGroup)                                       //1
		groupRouterGroup.POST("/group_application_response", group.ApplicationGroupResponse)        //1
		groupRouterGroup.POST("/transfer_group", group.TransferGroupOwner)                          //1
		groupRouterGroup.POST("/get_recv_group_applicationList", group.GetRecvGroupApplicationList) //1
		groupRouterGroup.POST("/get_user_req_group_applicationList", group.GetUserReqGroupApplicationList)
		groupRouterGroup.POST("/get_groups_info", group.GetGroupsInfo) //1
		groupRouterGroup.POST("/kick_group", group.KickGroupMember)    //1
		//	groupRouterGroup.POST("/get_group_member_list", group.GetGroupMemberList)        //no use
		groupRouterGroup.POST("/get_group_all_member_list", group.GetGroupAllMemberList) //1
		groupRouterGroup.POST("/get_group_members_info", group.GetGroupMembersInfo)      //1
		groupRouterGroup.POST("/invite_user_to_group", group.InviteUserToGroup)          //1
		//only for supergroup
		groupRouterGroup.POST("/invite_user_to_groups", group.InviteUserToGroups)
		groupRouterGroup.POST("/get_joined_group_list", group.GetJoinedGroupList)
		groupRouterGroup.POST("/dismiss_group", group.DismissGroup) //
		groupRouterGroup.POST("/mute_group_member", group.MuteGroupMember)
		groupRouterGroup.POST("/cancel_mute_group_member", group.CancelMuteGroupMember) //MuteGroup
		groupRouterGroup.POST("/mute_group", group.MuteGroup)
		groupRouterGroup.POST("/cancel_mute_group", group.CancelMuteGroup)
		groupRouterGroup.POST("/set_group_member_nickname", group.SetGroupMemberNickname)
		groupRouterGroup.POST("/set_group_member_info", group.SetGroupMemberInfo)
		groupRouterGroup.POST("/get_group_abstract_info", group.GetGroupAbstractInfo)
		//groupRouterGroup.POST("/get_group_all_member_list_by_split", group.GetGroupAllMemberListBySplit)
		groupRouterGroup.POST("/get_group_history_members", group.GetGroupHistoryMembers) // 获取群历史成员列表
	}
	superGroupRouterGroup := r.Group("/super_group")
	{
		superGroupRouterGroup.POST("/get_joined_group_list", group.GetJoinedSuperGroupList)
		superGroupRouterGroup.POST("/get_groups_info", group.GetSuperGroupsInfo)
	}
	//certificate
	authRouterGroup := r.Group("/auth")
	{
		authRouterGroup.POST("/user_register", apiAuth.UserRegister) //account rpc 调用
		authRouterGroup.POST("/user_token", apiAuth.UserToken)       //1
		authRouterGroup.POST("/parse_token", apiAuth.ParseToken)     //1
		authRouterGroup.POST("/force_logout", apiAuth.ForceLogout)   //1
	}
	//Third service
	thirdGroup := r.Group("/third")
	{
		thirdGroup.POST("/tencent_cloud_storage_credential", apiThird.TencentCloudStorageCredential)
		thirdGroup.POST("/ali_oss_credential", apiThird.AliOSSCredential)
		thirdGroup.POST("/minio_storage_credential", apiThird.MinioStorageCredential)
		thirdGroup.POST("/minio_upload", apiThird.MinioUploadFile)
		thirdGroup.POST("/upload_update_app", apiThird.UploadUpdateApp)
		thirdGroup.POST("/get_download_url", apiThird.GetDownloadURL)
		thirdGroup.POST("/get_rtc_invitation_info", apiThird.GetRTCInvitationInfo)
		thirdGroup.POST("/get_rtc_invitation_start_app", apiThird.GetRTCInvitationInfoStartApp)
		thirdGroup.POST("/fcm_update_token", apiThird.FcmUpdateToken)
		thirdGroup.POST("/aws_storage_credential", apiThird.AwsStorageCredential)
		thirdGroup.POST("/set_app_badge", apiThird.SetAppBadge)
	}
	//Message
	chatGroup := r.Group("/msg")
	{
		chatGroup.POST("/newest_seq", apiChat.GetSeq)
		chatGroup.POST("/send_msg", apiChat.SendMsg)
		chatGroup.POST("/pull_msg_by_seq", apiChat.PullMsgBySeqList)
		chatGroup.POST("/del_msg", apiChat.DelMsg)
		chatGroup.POST("/del_super_group_msg", apiChat.DelSuperGroupMsg)
		chatGroup.POST("/clear_msg", apiChat.ClearMsg)
		chatGroup.POST("/manage_send_msg", manage.ManagementSendMsg)
		chatGroup.POST("/batch_send_msg", manage.ManagementBatchSendMsg)
		chatGroup.POST("/check_msg_is_send_success", manage.CheckMsgIsSendSuccess)
		chatGroup.POST("/set_msg_min_seq", apiChat.SetMsgMinSeq)

		chatGroup.POST("/set_message_reaction_extensions", apiChat.SetMessageReactionExtensions)
		chatGroup.POST("/get_message_list_reaction_extensions", apiChat.GetMessageListReactionExtensions)
		chatGroup.POST("/add_message_reaction_extensions", apiChat.AddMessageReactionExtensions)
		chatGroup.POST("/delete_message_reaction_extensions", apiChat.DeleteMessageReactionExtensions)

		chatGroup.POST("/collect", apiChat.MsgCollect)          //消息收藏
		chatGroup.POST("/collect/list", apiChat.MsgCollectList) //消息收藏列表
		chatGroup.POST("/collect/del", apiChat.MsgCollectDel)   //删除消息收藏
	}
	//Conversation
	conversationGroup := r.Group("/conversation")
	{ //1
		conversationGroup.POST("/get_all_conversations", conversation.GetAllConversations)
		conversationGroup.POST("/get_conversation", conversation.GetConversation)
		conversationGroup.POST("/get_conversations", conversation.GetConversations)
		//deprecated
		conversationGroup.POST("/set_conversation", conversation.SetConversation)
		conversationGroup.POST("/batch_set_conversation", conversation.BatchSetConversations)
		//deprecated
		conversationGroup.POST("/set_recv_msg_opt", conversation.SetRecvMsgOpt)
		conversationGroup.POST("/modify_conversation_field", conversation.ModifyConversationField)
	}
	// office
	officeGroup := r.Group("/office")
	{
		officeGroup.POST("/get_user_tags", office.GetUserTags)
		officeGroup.POST("/get_user_tag_by_id", office.GetUserTagByID)
		officeGroup.POST("/create_tag", office.CreateTag)
		officeGroup.POST("/delete_tag", office.DeleteTag)
		officeGroup.POST("/set_tag", office.SetTag)
		officeGroup.POST("/send_msg_to_tag", office.SendMsg2Tag)
		officeGroup.POST("/get_send_tag_log", office.GetTagSendLogs)

		officeGroup.POST("/create_one_work_moment", office.CreateOneWorkMoment)
		officeGroup.POST("/delete_one_work_moment", office.DeleteOneWorkMoment)
		officeGroup.POST("/like_one_work_moment", office.LikeOneWorkMoment)
		officeGroup.POST("/comment_one_work_moment", office.CommentOneWorkMoment)
		officeGroup.POST("/get_work_moment_by_id", office.GetWorkMomentByID)
		officeGroup.POST("/get_user_work_moments", office.GetUserWorkMoments)
		officeGroup.POST("/get_user_friend_work_moments", office.GetUserFriendWorkMoments)
		officeGroup.POST("/set_user_work_moments_level", office.SetUserWorkMomentsLevel)
		officeGroup.POST("/delete_comment", office.DeleteComment)
	}

	organizationGroup := r.Group("/organization")
	{
		organizationGroup.POST("/create_department", organization.CreateDepartment)
		organizationGroup.POST("/update_department", organization.UpdateDepartment)
		organizationGroup.POST("/get_sub_department", organization.GetSubDepartment)
		organizationGroup.POST("/delete_department", organization.DeleteDepartment)
		organizationGroup.POST("/get_all_department", organization.GetAllDepartment)

		organizationGroup.POST("/create_organization_user", organization.CreateOrganizationUser)
		organizationGroup.POST("/update_organization_user", organization.UpdateOrganizationUser)
		organizationGroup.POST("/delete_organization_user", organization.DeleteOrganizationUser)

		organizationGroup.POST("/create_department_member", organization.CreateDepartmentMember)
		organizationGroup.POST("/get_user_in_department", organization.GetUserInDepartment)
		organizationGroup.POST("/update_user_in_department", organization.UpdateUserInDepartment)

		organizationGroup.POST("/get_department_member", organization.GetDepartmentMember)
		organizationGroup.POST("/delete_user_in_department", organization.DeleteUserInDepartment)
		organizationGroup.POST("/get_user_in_organization", organization.GetUserInOrganization)
	}

	initGroup := r.Group("/init")
	{
		initGroup.POST("/set_client_config", clientInit.SetClientInitConfig)
		initGroup.POST("/get_client_config", clientInit.GetClientInitConfig)
	}

	systemGroup := r.Group("/system")
	{
		systemGroup.POST("/wgt_version", system.WgtVersion)       //wgt版本
		systemGroup.POST("/latest_version", system.LatestVersion) //家等你app最新版本
	}

	return r
}
