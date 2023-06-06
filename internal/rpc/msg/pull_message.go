package msg

import (
	"context"
	"crazy_server/pkg/utils"
	go_redis "github.com/go-redis/redis/v8"

	commonDB "crazy_server/pkg/common/db"
	"crazy_server/pkg/common/log"
	crazy_server_sdk "crazy_server/pkg/proto/sdk_ws"

	promePkg "crazy_server/pkg/common/prometheus"
)

func (rpc *rpcChat) GetMaxAndMinSeq(_ context.Context, in *crazy_server_sdk.GetMaxAndMinSeqReq) (*crazy_server_sdk.GetMaxAndMinSeqResp, error) {
	log.NewInfo(in.OperationID, "rpc getMaxAndMinSeq is arriving", in.String())
	resp := new(crazy_server_sdk.GetMaxAndMinSeqResp)
	m := make(map[string]*crazy_server_sdk.MaxAndMinSeq)
	var maxSeq, minSeq uint64
	var err1, err2 error
	maxSeq, err1 = commonDB.DB.GetUserMaxSeq(in.UserID)
	minSeq, err2 = commonDB.DB.GetUserMinSeq(in.UserID)
	if (err1 != nil && err1 != go_redis.Nil) || (err2 != nil && err2 != go_redis.Nil) {
		log.NewError(in.OperationID, "getMaxSeq from redis error", in.String())
		if err1 != nil {
			log.NewError(in.OperationID, utils.GetSelfFuncName(), err1.Error())
		}
		if err2 != nil {
			log.NewError(in.OperationID, utils.GetSelfFuncName(), err2.Error())
		}
		resp.ErrCode = 200
		resp.ErrMsg = "redis get err"
		return resp, nil
	}
	resp.MaxSeq = uint32(maxSeq)
	resp.MinSeq = uint32(minSeq)
	for _, groupID := range in.GroupIDList {
		x := new(crazy_server_sdk.MaxAndMinSeq)
		maxSeq, _ := commonDB.DB.GetGroupMaxSeq(groupID)
		minSeq, _ := commonDB.DB.GetGroupUserMinSeq(groupID, in.UserID)
		x.MaxSeq = uint32(maxSeq)
		x.MinSeq = uint32(minSeq)
		m[groupID] = x
	}
	resp.GroupMaxAndMinSeq = m
	return resp, nil
}

func (rpc *rpcChat) PullMessageBySeqList(_ context.Context, in *crazy_server_sdk.PullMessageBySeqListReq) (*crazy_server_sdk.PullMessageBySeqListResp, error) {
	log.NewInfo(in.OperationID, "rpc PullMessageBySeqList is arriving", in.String())
	resp := new(crazy_server_sdk.PullMessageBySeqListResp)
	m := make(map[string]*crazy_server_sdk.MsgDataList)
	redisMsgList, failedSeqList, err := commonDB.DB.GetMessageListBySeq(in.UserID, in.SeqList, in.OperationID)
	if err != nil {
		if err != go_redis.Nil {
			promePkg.PromeAdd(promePkg.MsgPullFromRedisFailedCounter, len(failedSeqList))
			log.Error(in.OperationID, "get message from redis exception", err.Error(), failedSeqList)
		} else {
			log.Debug(in.OperationID, "get message from redis is nil", failedSeqList)
		}
		msgList, err1 := commonDB.DB.GetMsgBySeqListMongo2(in.UserID, failedSeqList, in.OperationID)
		if err1 != nil {
			promePkg.PromeAdd(promePkg.MsgPullFromMongoFailedCounter, len(failedSeqList))
			log.Error(in.OperationID, "PullMessageBySeqList data error", in.String(), err1.Error())
			resp.ErrCode = 201
			resp.ErrMsg = err1.Error()
			return resp, nil
		} else {
			promePkg.PromeAdd(promePkg.MsgPullFromMongoSuccessCounter, len(msgList))
			redisMsgList = append(redisMsgList, msgList...)
			resp.List = redisMsgList
		}
	} else {
		promePkg.PromeAdd(promePkg.MsgPullFromRedisSuccessCounter, len(redisMsgList))
		resp.List = redisMsgList
	}

	for k, v := range in.GroupSeqList {
		x := new(crazy_server_sdk.MsgDataList)
		redisMsgList, failedSeqList, err := commonDB.DB.GetMessageListBySeq(k, v.SeqList, in.OperationID)
		if err != nil {
			if err != go_redis.Nil {
				promePkg.PromeAdd(promePkg.MsgPullFromRedisFailedCounter, len(failedSeqList))
				log.Error(in.OperationID, "get message from redis exception", err.Error(), failedSeqList)
			} else {
				log.Debug(in.OperationID, "get message from redis is nil", failedSeqList)
			}
			msgList, _, err1 := commonDB.DB.GetSuperGroupMsgBySeqListMongo(k, failedSeqList, in.OperationID)
			if err1 != nil {
				promePkg.PromeAdd(promePkg.MsgPullFromMongoFailedCounter, len(failedSeqList))
				log.Error(in.OperationID, "PullMessageBySeqList data error", in.String(), err1.Error())
				resp.ErrCode = 201
				resp.ErrMsg = err1.Error()
				return resp, nil
			} else {
				promePkg.PromeAdd(promePkg.MsgPullFromMongoSuccessCounter, len(msgList))
				redisMsgList = append(redisMsgList, msgList...)
				x.MsgDataList = redisMsgList
				m[k] = x
			}
		} else {
			promePkg.PromeAdd(promePkg.MsgPullFromRedisSuccessCounter, len(redisMsgList))
			x.MsgDataList = redisMsgList
			m[k] = x
		}
	}
	resp.GroupMsgDataList = m
	return resp, nil
}

type MsgFormats []*crazy_server_sdk.MsgData

// Implement the sort.Interface interface to get the number of elements method
func (s MsgFormats) Len() int {
	return len(s)
}

//Implement the sort.Interface interface comparison element method
func (s MsgFormats) Less(i, j int) bool {
	return s[i].SendTime < s[j].SendTime
}

//Implement the sort.Interface interface exchange element method
func (s MsgFormats) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
