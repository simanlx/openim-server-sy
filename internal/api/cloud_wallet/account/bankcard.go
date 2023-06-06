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

// 绑定银行卡
func BindUserBankCard(c *gin.Context) {
	params := account.BindUserBankCardReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	//解析token、获取用户id
	userId, ok := common.ParseImToken(c, params.OperationID)
	if !ok {
		return
	}

	req := &rpc.BindUserBankcardReq{
		UserId:            userId,
		CardOwner:         params.CardOwner,
		BankCardNumber:    params.BankCard,
		Mobile:            params.Mobile,
		CardAvailableDate: params.CardAvailableDate,
		Cvv2:              params.Cvv2,
		OperationID:       params.OperationID,
	}

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCloudWalletName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	client := rpc.NewCloudWalletServiceClient(etcdConn)
	RpcResp, err := client.BindUserBankcard(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "BindUserBankcard failed ", err.Error(), req.String())
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

// 绑定银行卡code确认
func BindUserBankcardConfirm(c *gin.Context) {
	params := account.BindUserBankcardConfirmReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	//解析token、获取用户id
	userId, ok := common.ParseImToken(c, params.OperationID)
	if !ok {
		return
	}

	req := &rpc.BindUserBankcardConfirmReq{
		UserId:      userId,
		BankCardId:  params.BankCardId,
		SmsCode:     params.Code,
		MerUserIp:   c.ClientIP(),
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
	RpcResp, err := client.BindUserBankcardConfirm(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "BindUserBankcardConfirm failed ", err.Error(), req.String())
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

// 解绑银行卡
func UnBindUserBankcard(c *gin.Context) {
	params := account.UnBindUserBankcardReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	//解析token、获取用户id
	userId, ok := common.ParseImToken(c, params.OperationID)
	if !ok {
		return
	}

	req := &rpc.UnBindingUserBankcardReq{
		UserId:        userId,
		BindCardAgrNo: params.BindCardAgrNo,
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
	RpcResp, err := client.UnBindingUserBankcard(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "UnBindingUserBankcard failed ", err.Error(), req.String())
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
