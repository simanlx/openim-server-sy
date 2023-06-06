package agent

import (
	"crazy_server/internal/api/common"
	"crazy_server/pkg/base_info"
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/log"
	"crazy_server/pkg/grpc-etcdv3/getcdv3"
	rpc "crazy_server/pkg/proto/agent"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// 互娱商城购买咖豆下单(预提交)
func ChessShopPurchaseBean(c *gin.Context) {
	params := base_info.ChessShopPurchaseBeanReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	operationId := c.GetString("operationId")
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAgentName, operationId)
	if etcdConn == nil {
		errMsg := operationId + "getcdv3.GetDefaultConn == nil"
		log.NewError(operationId, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": errMsg})
		return
	}

	req := &rpc.ChessShopPurchaseBeanReq{
		AgentNumber:  params.AgentNumber,
		ChessOrderNo: params.ChessOrderNo,
		ChessUserId:  params.ChessUserId,
		ConfigId:     params.ConfigId,
		OperationId:  operationId,
	}

	client := rpc.NewAgentSystemServiceClient(etcdConn)
	RpcResp, err := client.ChessShopPurchaseBean(c, req)
	if err != nil {
		log.NewError(operationId, "ChessShopPurchaseBean failed ", err.Error(), req.String())
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	// handle rpc err
	if common.HandleAgentCommonRespErr(RpcResp.CommonResp, c) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": RpcResp})
	return
}

// 推广员购买咖豆
func PurchaseBean(c *gin.Context) {
	params := base_info.PurchaseBeanReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	operationId := c.GetString("operationId")
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAgentName, operationId)
	if etcdConn == nil {
		errMsg := operationId + "getcdv3.GetDefaultConn == nil"
		log.NewError(operationId, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": errMsg})
		return
	}

	req := &rpc.AgentPurchaseBeanReq{
		UserId:      c.GetString("userId"),
		ConfigId:    params.ConfigId,
		OperationId: operationId,
	}

	client := rpc.NewAgentSystemServiceClient(etcdConn)
	RpcResp, err := client.AgentPurchaseBean(c, req)
	if err != nil {
		log.NewError(operationId, "AgentPurchaseBean failed ", err.Error(), req.String())
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	// handle rpc err
	if common.HandleAgentCommonRespErr(RpcResp.CommonResp, c) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": RpcResp})
	return
}

// 推广员余额提现
func Withdraw(c *gin.Context) {
	params := base_info.WithdrawReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	operationId := c.GetString("operationId")
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAgentName, operationId)
	if etcdConn == nil {
		errMsg := operationId + "getcdv3.GetDefaultConn == nil"
		log.NewError(operationId, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": errMsg})
		return
	}

	//提现金额规则限制、最低10元
	if params.Amount < 1000 {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "最低提现金额为10元"})
		return
	}

	//提现金额规则限制、最高1万元
	if params.Amount > 10000*100 {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "单笔最高提现金额1万元"})
		return
	}

	req := &rpc.BalanceWithdrawalReq{
		UserId:          c.GetString("userId"),
		Amount:          params.Amount,
		PaymentPassword: params.PaymentPassword,
		OperationId:     operationId,
	}

	client := rpc.NewAgentSystemServiceClient(etcdConn)
	RpcResp, err := client.BalanceWithdrawal(c, req)
	if err != nil {
		log.NewError(operationId, "BalanceWithdrawal failed ", err.Error(), req.String())
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	// handle rpc err
	if common.HandleAgentCommonRespErr(RpcResp.CommonResp, c) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": RpcResp})
	return
}
