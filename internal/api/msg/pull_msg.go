package msg

import (
	"context"
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/log"
	"crazy_server/pkg/common/token_verify"
	"crazy_server/pkg/grpc-etcdv3/getcdv3"
	"crazy_server/pkg/proto/msg"
	crazy_server_sdk "crazy_server/pkg/proto/sdk_ws"
	"crazy_server/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type paramsUserPullMsg struct {
	ReqIdentifier *int   `json:"reqIdentifier" binding:"required"`
	SendID        string `json:"sendID" binding:"required"`
	OperationID   string `json:"operationID" binding:"required"`
	Data          struct {
		SeqBegin *int64 `json:"seqBegin" binding:"required"`
		SeqEnd   *int64 `json:"seqEnd" binding:"required"`
	}
}

type paramsUserPullMsgBySeqList struct {
	ReqIdentifier int      `json:"reqIdentifier" binding:"required"`
	SendID        string   `json:"sendID" binding:"required"`
	OperationID   string   `json:"operationID" binding:"required"`
	SeqList       []uint32 `json:"seqList"`
}

func PullMsgBySeqList(c *gin.Context) {
	params := paramsUserPullMsgBySeqList{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	token := c.Request.Header.Get("token")
	if ok, err := token_verify.VerifyToken(token, params.SendID); !ok {
		if err != nil {
			log.NewError(params.OperationID, utils.GetSelfFuncName(), err.Error(), token, params.SendID)
		}
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "token validate err"})
		return
	}
	pbData := crazy_server_sdk.PullMessageBySeqListReq{}
	pbData.UserID = params.SendID
	pbData.OperationID = params.OperationID
	pbData.SeqList = params.SeqList

	grpcConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, pbData.OperationID)
	if grpcConn == nil {
		errMsg := pbData.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(pbData.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	msgClient := msg.NewMsgClient(grpcConn)
	reply, err := msgClient.PullMessageBySeqList(context.Background(), &pbData)
	if err != nil {
		log.Error(pbData.OperationID, "PullMessageBySeqList error", err.Error())
		return
	}
	log.NewInfo(pbData.OperationID, "rpc call success to PullMessageBySeqList", reply.String(), len(reply.List))
	c.JSON(http.StatusOK, gin.H{
		"errCode":       reply.ErrCode,
		"errMsg":        reply.ErrMsg,
		"reqIdentifier": params.ReqIdentifier,
		"data":          reply.List,
	})
}
