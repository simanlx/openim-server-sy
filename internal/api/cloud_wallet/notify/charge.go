package notify

import (
	"context"
	"crazy_server/pkg/base_info/notify"
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/log"
	"crazy_server/pkg/grpc-etcdv3/getcdv3"
	rpc "crazy_server/pkg/proto/cloud_wallet"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// 充值回调
func ChargeNotify(c *gin.Context) {
	params := notify.ChargeNotifyReq{}
	if err := c.ShouldBind(&params); err != nil {
		log.Error("0", "ChargeNotify", err.Error(), params)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	req := &rpc.ChargeNotifyReq{
		MerOrderId:     params.MerOrderId,
		ResultCode:     params.ResultCode,
		ErrorCode:      params.ErrorCode,
		ErrorMsg:       params.ErrorMsg,
		NcountOrderId:  params.NcountOrderId,
		TranAmount:     params.TranAmount,
		SubmitTime:     params.SubmitTime,
		TranFinishTime: params.TranFinishTime,
		FeeAmount:      params.FeeAmount,
	}

	//调用rpc
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCloudWalletName, "0000")
	if etcdConn == nil {
		errMsg := "0000" + "getcdv3.GetDefaultConn == nil"
		log.NewError("0000", errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	// 创建rpc连接
	client := rpc.NewCloudWalletServiceClient(etcdConn)
	RpcResp, err := client.ChargeNotify(context.Background(), req)
	if err != nil {
		log.NewError("0000", "ChargeNotify failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 200, "data": RpcResp})
	return
}

// 提现回调
func WithDrawNotify(c *gin.Context) {
	params := notify.WithdrawNotifyReq{}
	if err := c.ShouldBind(&params); err != nil {
		log.Error("0", "WithDrawNotify", err.Error(), params)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	req := &rpc.DrawNotifyReq{
		MerOrderId:     params.MerOrderId,
		ResultCode:     params.ResultCode,
		ErrorCode:      params.ErrorCode,
		ErrorMsg:       params.ErrorMsg,
		NcountOrderId:  params.NcountOrderId,
		TranFinishDate: params.TranFinishDate,
		ServiceAmount:  params.ServiceAmount,
		PayAcctAmount:  params.PayAcctAmount,
	}

	fmt.Println("--req DrawNotifyReq", req, "---", params)

	//调用rpc
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCloudWalletName, "0000")
	if etcdConn == nil {
		errMsg := "0000" + "getcdv3.GetDefaultConn == nil"
		log.NewError("0000", errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	// 创建rpc连接
	client := rpc.NewCloudWalletServiceClient(etcdConn)
	RpcResp, err := client.WithDrawNotify(context.Background(), req)
	if err != nil {
		log.NewError("0000", "WithDrawNotify failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 200, "data": RpcResp})
	return
}
