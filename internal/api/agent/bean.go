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
	"time"
)

// 获取推广员游戏商城咖豆配置
func AgentGameShopBeanConfig(c *gin.Context) {
	params := base_info.AgentGameShopBeanConfigReq{}
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

	req := &rpc.AgentGameBeanShopConfigReq{
		UserId:      "",
		AgentNumber: params.AgentNumber,
		OperationId: operationId,
	}

	client := rpc.NewAgentSystemServiceClient(etcdConn)
	RpcResp, err := client.AgentGameBeanShopConfig(c, req)
	if err != nil {
		log.NewError(operationId, "AgentGameBeanShopConfig failed ", err.Error(), req.String())
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": RpcResp.BeanShopConfig})
	return
}

// 获取平台咖豆商城配置
func PlatformBeanShopConfig(c *gin.Context) {
	operationId := c.GetString("operationId")
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAgentName, operationId)
	if etcdConn == nil {
		errMsg := operationId + "getcdv3.GetDefaultConn == nil"
		log.NewError(operationId, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": errMsg})
		return
	}

	req := &rpc.PlatformBeanShopConfigReq{
		UserId:      "",
		OperationId: operationId,
	}

	client := rpc.NewAgentSystemServiceClient(etcdConn)
	RpcResp, err := client.PlatformBeanShopConfig(c, req)
	if err != nil {
		log.NewError(operationId, "AgentAccountIncomeChart failed ", err.Error(), req.String())
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": RpcResp.BeanShopConfig})
	return
}

// 推广员自定义咖豆商城配置
func AgentDiyBeanShopConfig(c *gin.Context) {
	operationId := c.GetString("operationId")
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAgentName, operationId)
	if etcdConn == nil {
		errMsg := operationId + "getcdv3.GetDefaultConn == nil"
		log.NewError(operationId, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": errMsg})
		return
	}

	req := &rpc.AgentDiyBeanShopConfigReq{
		UserId:      c.GetString("userId"),
		OperationId: operationId,
	}

	client := rpc.NewAgentSystemServiceClient(etcdConn)
	RpcResp, err := client.AgentDiyBeanShopConfig(c, req)
	if err != nil {
		log.NewError(operationId, "AgentDiyBeanShopConfig failed ", err.Error(), req.String())
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": RpcResp})
	return
}

// 咖豆账户明细详情列表
func AgentBeanAccountRecordList(c *gin.Context) {
	params := base_info.AgentBeanAccountRecordListReq{}
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

	//默认当天
	if params.Date == "" {
		params.Date = time.Now().Format("2006-01-02")
	}

	if params.Page == 0 {
		params.Page = 1
	}

	if params.Size == 0 || params.Size > 100 {
		params.Size = 20
	}

	req := &rpc.AgentBeanAccountRecordListReq{
		UserId:       c.GetString("userId"),
		Date:         params.Date,
		Page:         params.Page,
		Size:         params.Size,
		BusinessType: params.BusinessType,
		Keyword:      params.Keyword,
		OperationId:  operationId,
	}

	client := rpc.NewAgentSystemServiceClient(etcdConn)
	RpcResp, err := client.AgentBeanAccountRecordList(c, req)
	if err != nil {
		log.NewError(operationId, "AgentBeanAccountRecordList failed ", err.Error(), req.String())
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": RpcResp})
	return
}

// 咖豆管理上下架
func AgentBeanShopUpStatus(c *gin.Context) {
	params := base_info.AgentBeanShopUpStatusReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	if params.IsAll == 0 && params.ConfigId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "配置参数错误"})
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

	req := &rpc.AgentBeanShopUpStatusReq{
		UserId:      c.GetString("userId"),
		Status:      params.Status,
		ConfigId:    params.ConfigId,
		IsAll:       params.IsAll,
		OperationId: operationId,
	}

	client := rpc.NewAgentSystemServiceClient(etcdConn)
	RpcResp, err := client.AgentBeanShopUpStatus(c, req)
	if err != nil {
		log.NewError(operationId, "AgentBeanShopUpStatus failed ", err.Error(), req.String())
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

// 咖豆管理(新增、编辑)
func AgentBeanShopUpdate(c *gin.Context) {
	params := base_info.AgentBeanShopUpdateReq{}
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

	beanShopConfig := make([]*rpc.BeanShopConfig, 0)
	configLen := len(params.BeanShopConfig)
	if configLen > 0 {
		if configLen > 8 {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "咖豆配置不能超过8个"})
			return
		}

		for _, v := range params.BeanShopConfig {
			beanShopConfig = append(beanShopConfig, &rpc.BeanShopConfig{
				BeanNumber:     v.BeanNumber,
				GiveBeanNumber: v.GiveBeanNumber,
				Amount:         v.Amount,
				Status:         v.Status,
			})
		}
	}

	req := &rpc.AgentBeanShopUpdateReq{
		UserId:         c.GetString("userId"),
		BeanShopConfig: beanShopConfig,
		OperationId:    operationId,
	}

	client := rpc.NewAgentSystemServiceClient(etcdConn)
	RpcResp, err := client.AgentBeanShopUpdate(c, req)
	if err != nil {
		log.NewError(operationId, "AgentBeanShopUpdate failed ", err.Error(), req.String())
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

// 赠送下属成员咖豆
func AgentGiveMemberBean(c *gin.Context) {
	params := base_info.AgentGiveMemberBeanReq{}
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

	req := &rpc.AgentGiveMemberBeanReq{
		UserId:      c.GetString("userId"),
		ChessUserId: params.ChessUserId,
		BeanNumber:  params.BeanNumber,
		OperationId: operationId,
	}

	client := rpc.NewAgentSystemServiceClient(etcdConn)
	RpcResp, err := client.AgentGiveMemberBean(c, req)
	if err != nil {
		log.NewError(operationId, "AgentGiveMemberBean failed ", err.Error(), req.String())
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
