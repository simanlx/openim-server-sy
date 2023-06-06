package account

import (
	"context"
	"crazy_server/internal/api/common"
	"crazy_server/pkg/base_info/account"
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/log"
	"crazy_server/pkg/grpc-etcdv3/getcdv3"
	rpc "crazy_server/pkg/proto/cloud_wallet"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// 充值
func ChargeAccount(c *gin.Context) {
	params := account.UserRechargeReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	//校验金额
	if params.Amount%100 != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "充值金额错误"})
		return
	}

	//解析token、获取用户id
	userId, ok := common.ParseImToken(c, params.OperationID)
	if !ok {
		return
	}

	req := &rpc.UserRechargeReq{
		UserId:        userId,
		BindCardAgrNo: params.BindCardAgrNo,
		Amount:        params.Amount,
		OperationID:   params.OperationID,
	}

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCloudWalletName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	client := rpc.NewCloudWalletServiceClient(etcdConn)
	RpcResp, err := client.UserRecharge(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "UserRecharge failed ", err.Error(), req.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	// handle rpc err
	if common.HandleCommonRespErr(RpcResp.CommonResp, c) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"errCode": 200, "data": RpcResp})
	return
}

// 账户充值code确定
func ChargeAccountConfirm(c *gin.Context) {
	params := account.UserRechargeConfirmReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	//解析token、获取用户id
	userId, ok := common.ParseImToken(c, params.OperationID)
	if !ok {
		return
	}

	req := &rpc.UserRechargeConfirmReq{
		UserId:      userId,
		MerOrderId:  params.OrderNo,
		SmsCode:     params.Code,
		OperationID: params.OperationID,
	}

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCloudWalletName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	client := rpc.NewCloudWalletServiceClient(etcdConn)
	RpcResp, err := client.UserRechargeConfirm(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "UserRechargeConfirm failed ", err.Error(), req.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	// handle rpc err
	if common.HandleCommonRespErr(RpcResp.CommonResp, c) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"errCode": 200, "data": RpcResp})
	return
}

// 账户提现
func DrawAccount(c *gin.Context) {
	params := account.DrawAccountReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	//支付密码
	if len(params.PaymentPassword) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "请输入支付密码"})
		return
	}

	//提现金额限制
	if params.Amount < 1000 {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "提现金额最少10元"})
		return
	}

	//解析token、获取用户id
	userId, ok := common.ParseImToken(c, params.OperationID)
	if !ok {
		return
	}

	req := &rpc.DrawAccountReq{
		UserId:          userId,
		BindCardAgrNo:   params.BindCardAgrNo,
		Amount:          params.Amount,
		PaymentPassword: params.PaymentPassword,
		OperationID:     params.OperationID,
	}

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCloudWalletName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	client := rpc.NewCloudWalletServiceClient(etcdConn)
	RpcResp, err := client.UserWithdrawal(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "UserWithdrawal failed ", err.Error(), req.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	// handle rpc err
	if common.HandleCommonRespErr(RpcResp.CommonResp, c) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"errCode": 200, "data": RpcResp})
	return
}
