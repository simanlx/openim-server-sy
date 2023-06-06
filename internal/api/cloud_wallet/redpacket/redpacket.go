package redpacket

import (
	"context"
	"crazy_server/internal/api/common"
	"crazy_server/pkg/agora"
	"crazy_server/pkg/base_info/redpacket_struct"
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/log"
	"crazy_server/pkg/grpc-etcdv3/getcdv3"
	rpc "crazy_server/pkg/proto/cloud_wallet"
	"crazy_server/pkg/tencent_cloud"
	"crazy_server/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// 发送红包接口
func SendRedPacket(c *gin.Context) {
	params := redpacket_struct.SendRedPacket{}
	if err := c.BindJSON(&params); err != nil {
		log.Error("0", "ChargeNotify", err.Error(), params)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	//解析token、获取用户id
	userId, ok := common.ParseImToken(c, params.OperationID)
	if !ok {
		return
	}
	params.UserId = userId

	UserID := "10018"
	if params.UserId != "" {
		UserID = params.UserId
	}
	// 复制结构体
	req := &rpc.SendRedPacketReq{}
	err := utils.CopyStructFields(req, &params)
	if err != nil {
		log.NewError(params.OperationID, "CopyStructFields failed ", err.Error(), params)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	req.UserId = UserID

	//调用rpc
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCloudWalletName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	// 创建rpc连接
	client := rpc.NewCloudWalletServiceClient(etcdConn)
	RpcResp, err := client.SendRedPacket(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "调用失败 ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errCode": 200, "data": RpcResp})
	return
}

// 抢红包接口
func ClickRedPacket(c *gin.Context) {
	params := redpacket_struct.ClickRedPacketReq{}
	if err := c.BindJSON(&params); err != nil {
		log.Error("0", "ChargeNotify", err.Error(), params)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	//解析token、获取用户id
	userId, ok := common.ParseImToken(c, params.OperateID)
	if !ok {
		return
	}
	params.UserId = userId

	// 复制结构体
	req := &rpc.ClickRedPacketReq{}
	utils.CopyStructFields(req, &params)

	//调用rpc
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCloudWalletName, req.OperateID)
	if etcdConn == nil {
		errMsg := req.OperateID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperateID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	// 创建rpc连接
	client := rpc.NewCloudWalletServiceClient(etcdConn)
	RpcResp, err := client.ClickRedPacket(context.Background(), req)
	if err != nil {
		log.NewError(req.OperateID, "BindUserBankcard failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errCode": RpcResp.CommonResp.ErrCode, "errMsg": RpcResp.CommonResp.ErrMsg})
	return
}

// 确认发送红包
func SendRedPacketConfirm(c *gin.Context) {
	params := redpacket_struct.ConfirmSendRedPacketReq{}
	if err := c.BindJSON(&params); err != nil {
		log.Error("0", "ChargeNotify", err.Error(), params)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	//解析token、获取用户id
	_, ok := common.ParseImToken(c, params.OperateID)
	if !ok {
		return
	}

	// 复制结构体
	req := &rpc.SendRedPacketConfirmReq{}
	utils.CopyStructFields(req, &params)

	//调用rpc
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCloudWalletName, req.OperateID)
	if etcdConn == nil {
		errMsg := req.OperateID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperateID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	// 创建rpc连接
	client := rpc.NewCloudWalletServiceClient(etcdConn)
	RpcResp, err := client.SendRedPacketConfirm(context.Background(), req)
	if err != nil {
		log.NewError(req.OperateID, "BindUserBankcard failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errCode": RpcResp.CommonResp.ErrCode, "errMsg": RpcResp.CommonResp.ErrMsg})
	return
}

// 红包领取明细
func RedPacketReceiveDetail(c *gin.Context) {
	params := redpacket_struct.RedPacketReceiveDetailReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	//解析token、获取用户id
	userId, ok := common.ParseImToken(c, params.OperationID)
	if !ok {
		return
	}

	req := &rpc.RedPacketReceiveDetailReq{
		UserId:      userId,
		StartTime:   params.StartTime,
		EndTime:     params.EndTime,
		OperationID: params.OperationID,
	}

	//调用rpc
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCloudWalletName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	// 创建rpc连接
	client := rpc.NewCloudWalletServiceClient(etcdConn)
	RpcResp, err := client.RedPacketReceiveDetail(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "RedPacketReceiveDetail failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errCode": 200, "data": RpcResp})
	return
}

// 获取红包详情接口
func GetRedPacketInfo(c *gin.Context) {
	params := redpacket_struct.RedPacketInfoReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	//解析token、获取用户id
	userId, ok := common.ParseImToken(c, params.OperationID)
	if !ok {
		return
	}

	req := &rpc.RedPacketInfoReq{
		UserId:      userId,
		PacketId:    params.PacketId,
		OperationID: params.OperationID,
	}

	//调用rpc
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCloudWalletName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	// 创建rpc连接
	client := rpc.NewCloudWalletServiceClient(etcdConn)
	RpcResp, err := client.RedPacketInfo(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "RedPacketInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errCode": 200, "timestamp": time.Now().Unix(), "data": RpcResp})
	return
}

// 禁止用户抢红包
func BanGroupClickRedPacket(c *gin.Context) {
	params := redpacket_struct.BanRedPacketReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	// 获取用户ID
	/*	ok, UserID, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), params.OperateID)
		if !ok {
			errMsg := params.OperateID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
			log.NewError(params.OperateID, errMsg)
			c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
			return
		}*/
	UserID := "10000"
	if params.UserId == "" {
		UserID = "1000"
	} else {
		UserID = params.UserId
	}

	req := &rpc.ForbidGroupRedPacketReq{
		UserId:      UserID,
		GroupId:     params.GroupId,
		Forbid:      params.Forbid,
		OperationID: params.OperationID,
	}

	//调用rpc
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCloudWalletName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	// 创建rpc连接
	client := rpc.NewCloudWalletServiceClient(etcdConn)
	RpcResp, err := client.ForbidGroupRedPacket(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "RedPacketInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errCode": 200, "data": RpcResp})
	return
}

// 这里是获取声网（完成）
func GetAgoraToken(c *gin.Context) {
	params := redpacket_struct.AgoraTokenReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	//解析token、获取用户id
	_, ok := common.ParseImToken(c, params.OperationID)
	if !ok {
		return
	}
	//
	//var role rtctokenbuilder.Role
	//switch params.Role {
	//case 1:
	//	role = rtctokenbuilder.RolePublisher
	//case 2:
	//	role = rtctokenbuilder.RoleSubscriber
	//}

	// 生成token
	result, appid, err := agora.GenerateRtcToken("0", params.OperationID, params.Channel_name, int(params.Role))
	if err != nil {
		log.NewError(params.OperationID, "RedPacketInfo failed ", err.Error(), params)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp := redpacket_struct.AgoraTokenResp{
		Token: result,
		AppID: appid,
	}

	c.JSON(http.StatusOK, gin.H{"errCode": 200, "msg": "获取成功", "data": resp})
}

// 翻译音频 （完成）
func TranslateVideo(c *gin.Context) {
	params := redpacket_struct.TencentMsgEscapeReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	//解析token、获取用户id
	_, ok := common.ParseImToken(c, params.OperationID)
	if !ok {
		return
	}

	// 消息翻译
	result, err := tencent_cloud.GetTencentCloudTranslate(params.ContentUrl, params.OperationID)
	if err != nil {
		log.NewError(params.OperationID, "RedPacketInfo failed ", err.Error(), params)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 200, "errMsg": "获取成功", "data": result})
}

// 获取版本号
func GetVersion(c *gin.Context) {
	// param in : 版本号
	// param out : 最新版本号、下载地址、更新内容、是否强制更新

	params := redpacket_struct.GetVersionReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	rpcReq := rpc.GetVersionReq{
		Version:     params.VersionCode,
		OperationID: params.OperationID,
	}

	// etcdConn
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCloudWalletName, params.OperationID)
	if etcdConn == nil {
		errMsg := params.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	// 创建rpc连接
	client := rpc.NewCloudWalletServiceClient(etcdConn)
	RpcResp, err := client.GetVersion(context.Background(), &rpcReq)
	if err != nil {
		log.NewError(params.OperationID, "GetVersion failed ", err.Error(), params)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errCode": 200, "data": RpcResp})
}

// 红包退还
func RefoundPacket(c *gin.Context) {
	// param in : 版本号
	// param out : 最新版本号、下载地址、更新内容、是否强制更新

	params := redpacket_struct.ReFoundPacketReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	if params.Secret != redpacket_struct.ReFoundPacketSecret {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "秘钥错误"})
		return
	}

	rpcReq := rpc.RefoundPacketReq{
		IP:          c.Request.RemoteAddr,
		OperationID: params.OperationID,
	}

	// etcdConn
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCloudWalletName, params.OperationID)
	if etcdConn == nil {
		errMsg := params.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	// 创建rpc连接
	client := rpc.NewCloudWalletServiceClient(etcdConn)
	RpcResp, err := client.RefoundPacket(context.Background(), &rpcReq)
	if err != nil {
		log.NewError(params.OperationID, "GetVersion failed ", err.Error(), params)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errCode": 200, "data": RpcResp})
}

// 第三方支付
func ThirdPay(c *gin.Context) {
	params := redpacket_struct.ThirdPayReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	//解析token、获取用户id
	userID, ok := common.ParseImToken(c, params.OprationID)
	if !ok {
		// 用户token解析失败
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "非法token"})
		return
	}
	//rpc
	rpcReq := rpc.ThirdPayReq{
		OperationID:      params.OprationID,
		OrderNo:          params.OrderNo,
		Password:         params.Password,
		SendType:         params.SendType,
		BankcardProtocol: params.BankcardProtocol,
		Userid:           userID,
	}

	//etcdConn
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCloudWalletName, params.OprationID)
	if etcdConn == nil {
		errMsg := params.OprationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(params.OprationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	// 创建rpc连接
	client := rpc.NewCloudWalletServiceClient(etcdConn)
	RpcResp, err := client.ThirdPay(context.Background(), &rpcReq)
	if err != nil {
		log.NewError(params.OprationID, "GetVersion failed ", err.Error(), params)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errCode": 200, "data": RpcResp})
}

// 从新生账户提现到用户的银行卡
func ThirdWithdraw(c *gin.Context) {
	params := redpacket_struct.ThirdWithdrawReq{}
	// 创建提现订单
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	//解析token、获取用户id
	userid, ok := common.ParseImToken(c, params.OperationID)
	if !ok {
		// 用户token解析失败
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "非法token"})
		return
	}

	//rpc
	rpcReq := rpc.ThirdWithdrawalReq{
		ThirdOrderId: params.ThirdOrderId,
		Amount:       params.Amount,
		Commission:   params.Commission,
		Password:     params.Password,
		UserId:       userid,
		OperationID:  params.OperationID,
	}
	//etcdConn
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCloudWalletName, params.OperationID)
	if etcdConn == nil {
		errMsg := "getcdv3.GetDefaultConn == nil"
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	// 创建rpc连接
	client := rpc.NewCloudWalletServiceClient(etcdConn)
	RpcResp, err := client.ThirdWithdrawal(context.Background(), &rpcReq)
	if err != nil {
		log.NewError(params.OperationID, "GetVersion failed ", err.Error(), params)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 200, "data": RpcResp})
}

// 创建第三方支付订单
func CreateThirdPayOrder(c *gin.Context) {
	params := redpacket_struct.CreateThirdPayOrder{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	operation := utils.OperationIDGenerator()
	// 这个接口不需要进行token
	//rpc
	rpcReq := rpc.CreateThirdPayOrderReq{
		MerchantId:  params.MerchantID,
		MerOrderId:  params.MerOrderID,
		NotifyUrl:   params.NotifyURL,
		Amount:      params.Amount,
		Remark:      params.Remark,
		OperationID: operation,
	}

	//etcdConn
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCloudWalletName, operation)
	if etcdConn == nil {
		errMsg := "getcdv3.GetDefaultConn == nil"
		log.NewError(operation, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	// 创建rpc连接
	client := rpc.NewCloudWalletServiceClient(etcdConn)
	RpcResp, err := client.CreateThirdPayOrder(context.Background(), &rpcReq)
	if err != nil {
		log.NewError(operation, "GetVersion failed ", err.Error(), params)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 200, "data": RpcResp})
}

// 查询第三方订单
func GetThirdPayOrder(c *gin.Context) {
	params := redpacket_struct.GetThirdPayOrder{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	rpcReq := rpc.GetThirdPayOrderInfoReq{
		OrderNo:     params.OrderNO,
		OperationID: params.OperationID,
	}

	//etcdConn
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCloudWalletName, params.OperationID)
	if etcdConn == nil {
		errMsg := "getcdv3.GetDefaultConn == nil"
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	// 创建rpc连接
	client := rpc.NewCloudWalletServiceClient(etcdConn)
	RpcResp, err := client.GetThirdPayOrderInfo(context.Background(), &rpcReq)
	if err != nil {
		log.NewError(params.OperationID, "GetVersion failed ", err.Error(), params)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errCode": 200, "data": RpcResp})
}

// 第三方回调接口，通过业务的标识的透传来区分业务类型
func ThirdPayCallback(c *gin.Context) {
	params := redpacket_struct.ThirdPayCallback{}
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	operationID := utils.OperationIDGenerator()
	// 获取http query参数，查看是否存在query参数
	bussinessType := c.GetInt("bussinessType")
	if bussinessType == 0 {
		bussinessType = 100
	}
	// rpc
	rpcReq := rpc.PayCallbackReq{
		MerOrderId:     params.MerOrderId,
		ResultCode:     params.ResultCode,
		ErrorCode:      params.ErrorCode,
		ErrorMsg:       params.ErrorMsg,
		NcountOrderId:  params.NcountOrderId,
		TranAmount:     params.TranAmount,
		SubmitTime:     params.SubmitTime,
		TranFinishTime: params.TranFinishTime,
		FeeAmount:      params.FeeAmount,
		BusinessType:   rpc.PayCallbackBusinessType(bussinessType),
	}
	//etcdConn
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCloudWalletName, operationID)
	if etcdConn == nil {
		errMsg := "getcdv3.GetDefaultConn == nil"
		log.NewError(operationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	// 创建rpc连接
	client := rpc.NewCloudWalletServiceClient(etcdConn)
	RpcResp, err := client.PayCallback(context.Background(), &rpcReq)
	if err != nil {
		log.NewError(operationID, "GetVersion failed ", err.Error(), params)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 200, "data": RpcResp})
}

// 确认支付接口 ,通过typeID来进行业务区别
func PayConfirm(c *gin.Context) {
	params := redpacket_struct.ThirdPayConfirm{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	//解析token、获取用户id
	_, ok := common.ParseImToken(c, params.OperationID)
	if !ok {
		// 用户token解析失败
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "非法token"})
		return
	}

	//rpc
	rpcReq := rpc.PayConfirmReq{
		OrderNo:     params.OrderNo,
		Code:        params.Code,
		OperationID: params.OperationID,
	}

	//etcdConn
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCloudWalletName, params.OperationID)
	if etcdConn == nil {
		errMsg := params.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	// 创建rpc连接
	client := rpc.NewCloudWalletServiceClient(etcdConn)
	RpcResp, err := client.PayConfirm(context.Background(), &rpcReq)
	if err != nil {
		log.NewError(params.OperationID, "GetVersion failed ", err.Error(), params)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 200, "data": RpcResp})
}
