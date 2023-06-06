package register

import (
	"context"
	"crazy_server/pkg/common/config"
	imdb "crazy_server/pkg/common/db/mysql_model/im_mysql_model"
	"crazy_server/pkg/common/log"
	"crazy_server/pkg/grpc-etcdv3/getcdv3"
	pbFriend "crazy_server/pkg/proto/friend"
	"crazy_server/pkg/utils"
	"strings"
)

var ChImportFriend chan *pbFriend.ImportFriendReq

func init() {
	ChImportFriend = make(chan *pbFriend.ImportFriendReq, 1000)
}

func ImportFriendRoutine() {
	for {
		req := <-ChImportFriend
		go func() {
			friendUserIDList, err := imdb.GetRegisterAddFriendList(0, 0)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), req, err.Error())
				return
			}
			log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "ImportFriendRoutine IDList", friendUserIDList)
			if len(friendUserIDList) == 0 {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "len==0")
				return
			}
			req.FriendUserIDList = friendUserIDList
			etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName, req.OperationID)
			if etcdConn == nil {
				errMsg := req.OperationID + "getcdv3.GetConn == nil"
				log.NewError(req.OperationID, errMsg)
				return
			}
			client := pbFriend.NewFriendClient(etcdConn)
			rpcResp, err := client.ImportFriend(context.Background(), req)
			if err != nil {
				log.NewError(req.OperationID, "ImportFriend failed ", err.Error(), req.String())
				return
			}
			if rpcResp.CommonResp.ErrCode != 0 {
				log.NewError(req.OperationID, "ImportFriend failed ", rpcResp)
			}
		}()
	}
}
