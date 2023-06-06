package msg

import (
	"context"
	"crazy_server/internal/api/common"
	api "crazy_server/pkg/base_info"
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/log"
	"crazy_server/pkg/grpc-etcdv3/getcdv3"
	pbChat "crazy_server/pkg/proto/msg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// 消息收藏列表
func MsgCollectList(c *gin.Context) {
	params := api.MsgCollectListReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	//解析token、获取用户id
	userId, ok := common.ParseImToken(c, params.OperationID)
	if !ok {
		return
	}

	req := &pbChat.MsgCollectListReq{
		MsgType:     params.MsgType,
		Keyword:     params.Keyword,
		Page:        params.Page,
		Size:        params.Size,
		OperationID: params.OperationID,
		UserId:      userId,
	}
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, params.OperationID)
	if etcdConn == nil {
		errMsg := params.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	msgClient := pbChat.NewMsgClient(etcdConn)
	RpcResp, err := msgClient.MsgCollectList(context.Background(), req)
	if err != nil {
		log.NewError(params.OperationID, "MsgCollectList failed ", err.Error(), req.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errCode": 200, "data": RpcResp})
	return

}

// 消息收藏
func MsgCollect(c *gin.Context) {
	params := api.MsgCollectReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	//解析token、获取用户id
	userId, ok := common.ParseImToken(c, params.OperationID)
	if !ok {
		return
	}

	req := &pbChat.MsgCollectReq{
		MsgType:     params.MsgType,
		Content:     params.Content,
		UserId:      userId,
		OperationID: params.OperationID,
	}
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, params.OperationID)
	if etcdConn == nil {
		errMsg := params.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	msgClient := pbChat.NewMsgClient(etcdConn)
	RpcResp, err := msgClient.MsgCollect(context.Background(), req)
	if err != nil {
		log.NewError(params.OperationID, "MsgCollect failed ", err.Error(), req.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errCode": 200, "data": RpcResp})
	return
}

// 删除消息收藏
func MsgCollectDel(c *gin.Context) {
	params := api.DelMsgCollectReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	//解析token、获取用户id
	userId, ok := common.ParseImToken(c, params.OperationID)
	if !ok {
		return
	}

	req := &pbChat.MsgCollectDelReq{
		CollectId:   params.CollectId,
		UserId:      userId,
		OperationID: params.OperationID,
	}
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, params.OperationID)
	if etcdConn == nil {
		errMsg := params.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	msgClient := pbChat.NewMsgClient(etcdConn)
	RpcResp, err := msgClient.MsgCollectDel(context.Background(), req)
	if err != nil {
		log.NewError(params.OperationID, "MsgCollectDel failed ", err.Error(), req.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errCode": 200, "data": RpcResp})
	return
}
