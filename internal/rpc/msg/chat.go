package msg

import (
	"context"
	"crazy_server/pkg/common/constant"
	"crazy_server/pkg/common/db"
	"crazy_server/pkg/common/log"
	"crazy_server/pkg/common/token_verify"
	pbChat "crazy_server/pkg/proto/msg"
	"crazy_server/pkg/utils"
)

func (rpc *rpcChat) ClearMsg(_ context.Context, req *pbChat.ClearMsgReq) (*pbChat.ClearMsgResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc req: ", req.String())
	if req.OpUserID != req.UserID && !token_verify.IsManagerUserID(req.UserID) {
		errMsg := "No permission" + req.OpUserID + req.UserID
		log.Error(req.OperationID, errMsg)
		return &pbChat.ClearMsgResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}, nil
	}
	log.Debug(req.OperationID, "CleanUpOneUserAllMsgFromRedis args", req.UserID)
	err := db.DB.CleanUpOneUserAllMsgFromRedis(req.UserID, req.OperationID)
	if err != nil {
		errMsg := "CleanUpOneUserAllMsgFromRedis failed " + err.Error() + req.OperationID + req.UserID
		log.Error(req.OperationID, errMsg)
		return &pbChat.ClearMsgResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}
	log.Debug(req.OperationID, "CleanUpUserMsgFromMongo args", req.UserID)
	err = db.DB.CleanUpUserMsgFromMongo(req.UserID, req.OperationID)
	if err != nil {
		errMsg := "CleanUpUserMsgFromMongo failed " + err.Error() + req.OperationID + req.UserID
		log.Error(req.OperationID, errMsg)
		return &pbChat.ClearMsgResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}

	resp := pbChat.ClearMsgResp{ErrCode: 0}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return &resp, nil
}

func (rpc *rpcChat) SetMsgMinSeq(_ context.Context, req *pbChat.SetMsgMinSeqReq) (*pbChat.SetMsgMinSeqResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc req: ", req.String())
	if req.OpUserID != req.UserID && !token_verify.IsManagerUserID(req.UserID) {
		errMsg := "No permission" + req.OpUserID + req.UserID
		log.Error(req.OperationID, errMsg)
		return &pbChat.SetMsgMinSeqResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}, nil
	}
	if req.GroupID == "" {
		err := db.DB.SetUserMinSeq(req.UserID, req.MinSeq)
		if err != nil {
			errMsg := "SetUserMinSeq failed " + err.Error() + req.OperationID + req.UserID + utils.Uint32ToString(req.MinSeq)
			log.Error(req.OperationID, errMsg)
			return &pbChat.SetMsgMinSeqResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
		}
		return &pbChat.SetMsgMinSeqResp{}, nil
	}
	err := db.DB.SetGroupUserMinSeq(req.GroupID, req.UserID, uint64(req.MinSeq))
	if err != nil {
		errMsg := "SetGroupUserMinSeq failed " + err.Error() + req.OperationID + req.GroupID + req.UserID + utils.Uint32ToString(req.MinSeq)
		log.Error(req.OperationID, errMsg)
		return &pbChat.SetMsgMinSeqResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}
	return &pbChat.SetMsgMinSeqResp{}, nil
}
