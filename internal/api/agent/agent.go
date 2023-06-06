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

// 推广员申请提交
func AgentApply(c *gin.Context) {
	params := base_info.AgentApplyReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	operationId := c.GetString("operationId")
	req := &rpc.AgentApplyReq{
		UserId:      c.GetString("userId"),
		Name:        params.Name,
		Mobile:      params.Mobile,
		ChessUserId: params.ChessUserId,
		OperationId: operationId,
	}

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAgentName, operationId)
	if etcdConn == nil {
		errMsg := operationId + "getcdv3.GetDefaultConn == nil"
		log.NewError(operationId, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": errMsg})
		return
	}

	client := rpc.NewAgentSystemServiceClient(etcdConn)
	RpcResp, err := client.AgentApply(c, req)
	if err != nil {
		log.NewError(operationId, "AgentApply failed ", err.Error(), req.String())
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

// 获取当前用户的推广员信息以及绑定关系
func GetUserAgentInfo(c *gin.Context) {
	params := base_info.GetUserAgentInfoReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	operationId := c.GetString("operationId")
	req := &rpc.GetUserAgentInfoReq{
		UserId:      c.GetString("userId"),
		ChessUserId: params.ChessUserId,
		OperationId: operationId,
	}

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAgentName, operationId)
	if etcdConn == nil {
		errMsg := operationId + "getcdv3.GetDefaultConn == nil"
		log.NewError(operationId, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": errMsg})
		return
	}

	client := rpc.NewAgentSystemServiceClient(etcdConn)
	RpcResp, err := client.GetUserAgentInfo(c, req)
	if err != nil {
		log.NewError(operationId, "GetUserAgentInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": RpcResp})
	return
}

// 绑定推广员
func BindAgentNumber(c *gin.Context) {
	params := base_info.BindAgentNumberReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	operationId := c.GetString("operationId")
	req := &rpc.BindAgentNumberReq{
		UserId:        "",
		AgentNumber:   params.AgentNumber,
		ChessUserId:   params.ChessUserId,
		ChessNickname: params.ChessNickname,
		OperationId:   operationId,
	}

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAgentName, operationId)
	if etcdConn == nil {
		errMsg := operationId + "getcdv3.GetDefaultConn == nil"
		log.NewError(operationId, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": errMsg})
		return
	}

	client := rpc.NewAgentSystemServiceClient(etcdConn)
	RpcResp, err := client.BindAgentNumber(c, req)
	if err != nil {
		log.NewError(operationId, "BindAgentNumber failed ", err.Error(), req.String())
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

// 推广员主页信息
func AgentMainInfo(c *gin.Context) {
	operationId := c.GetString("operationId")
	req := &rpc.AgentMainInfoReq{
		UserId:      c.GetString("userId"),
		OperationId: operationId,
	}

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAgentName, operationId)
	if etcdConn == nil {
		errMsg := operationId + "getcdv3.GetDefaultConn == nil"
		log.NewError(operationId, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": errMsg})
		return
	}

	client := rpc.NewAgentSystemServiceClient(etcdConn)
	RpcResp, err := client.AgentMainInfo(c, req)
	if err != nil {
		log.NewError(operationId, "AgentMainInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": RpcResp})
	return
}

// 账户明细收益趋势图
func AgentAccountIncomeChart(c *gin.Context) {
	params := base_info.AgentAccountIncomeChartReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	//默认7天
	if params.DateType == 0 {
		params.DateType = 1
	}

	operationId := c.GetString("operationId")
	req := &rpc.AgentAccountIncomeChartReq{
		UserId:      c.GetString("userId"),
		DateType:    params.DateType,
		OperationId: operationId,
	}

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAgentName, operationId)
	if etcdConn == nil {
		errMsg := operationId + "getcdv3.GetDefaultConn == nil"
		log.NewError(operationId, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": errMsg})
		return
	}

	client := rpc.NewAgentSystemServiceClient(etcdConn)
	RpcResp, err := client.AgentAccountIncomeChart(c, req)
	if err != nil {
		log.NewError(operationId, "AgentAccountIncomeChart failed ", err.Error(), req.String())
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": RpcResp.IncomeChartData})
	return
}

// 账户明细详情列表
func AgentAccountRecordList(c *gin.Context) {
	params := base_info.AgentAccountRecordListReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
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

	operationId := c.GetString("operationId")
	req := &rpc.AgentAccountRecordListReq{
		UserId:       c.GetString("userId"),
		Date:         params.Date,
		BusinessType: params.BusinessType,
		Keyword:      params.Keyword,
		Page:         params.Page,
		Size:         params.Size,
		OperationId:  operationId,
	}

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAgentName, operationId)
	if etcdConn == nil {
		errMsg := operationId + "getcdv3.GetDefaultConn == nil"
		log.NewError(operationId, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": errMsg})
		return
	}

	client := rpc.NewAgentSystemServiceClient(etcdConn)
	RpcResp, err := client.AgentAccountRecordList(c, req)
	if err != nil {
		log.NewError(operationId, "AgentAccountRecordList failed ", err.Error(), req.String())
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": RpcResp})
	return
}

// 推广下属用户列表
func AgentMemberList(c *gin.Context) {
	params := base_info.AgentMemberListReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	operationId := c.GetString("operationId")
	req := &rpc.AgentMemberListReq{
		UserId:      c.GetString("userId"),
		Keyword:     params.Keyword,
		Page:        params.Page,
		Size:        params.Size,
		OrderBy:     params.OrderBy,
		OperationId: operationId,
	}

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAgentName, operationId)
	if etcdConn == nil {
		errMsg := operationId + "getcdv3.GetDefaultConn == nil"
		log.NewError(operationId, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": errMsg})
		return
	}

	client := rpc.NewAgentSystemServiceClient(etcdConn)
	RpcResp, err := client.AgentMemberList(c, req)
	if err != nil {
		log.NewError(operationId, "AgentMemberList failed ", err.Error(), req.String())
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": RpcResp})
	return
}

// 开通推广员
func OpenAgent(c *gin.Context) {
	params := base_info.OpenAgentReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	req := &rpc.OpenAgentReq{
		ApplyId: params.ApplyId,
	}

	operationId := ""
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAgentName, operationId)
	if etcdConn == nil {
		errMsg := operationId + "getcdv3.GetDefaultConn == nil"
		log.NewError(operationId, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": errMsg})
		return
	}

	client := rpc.NewAgentSystemServiceClient(etcdConn)
	RpcResp, err := client.OpenAgent(c, req)
	if err != nil {
		log.NewError(operationId, "OpenAgent failed ", err.Error(), req.String())
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
