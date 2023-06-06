package friend

import (
	cbApi "crazy_server/pkg/call_back_struct"
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/constant"
	"crazy_server/pkg/common/http"
	"crazy_server/pkg/common/log"
	pbFriend "crazy_server/pkg/proto/friend"
	//"crazy_server/pkg/proto/msg"
	"crazy_server/pkg/utils"
	http2 "net/http"
)

func callbackBeforeAddFriend(req *pbFriend.AddFriendReq) cbApi.CommonCallbackResp {
	callbackResp := cbApi.CommonCallbackResp{OperationID: req.CommID.OperationID}
	if !config.Config.Callback.CallbackBeforeAddFriend.Enable {
		return callbackResp
	}
	log.NewDebug(req.CommID.OperationID, utils.GetSelfFuncName(), req.String())
	commonCallbackReq := &cbApi.CallbackBeforeAddFriendReq{
		CallbackCommand: constant.CallbackBeforeAddFriendCommand,
		FromUserID:      req.CommID.FromUserID,
		ToUserID:        req.CommID.ToUserID,
		ReqMsg:          req.ReqMsg,
		OperationID:     req.CommID.OperationID,
	}
	resp := &cbApi.CallbackBeforeAddFriendResp{
		CommonCallbackResp: &callbackResp,
	}
	//utils.CopyStructFields(req, msg.MsgData)
	defer log.NewDebug(req.CommID.OperationID, utils.GetSelfFuncName(), commonCallbackReq, *resp)
	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackBeforeAddFriendCommand, commonCallbackReq, resp, config.Config.Callback.CallbackBeforeAddFriend.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
		if !config.Config.Callback.CallbackBeforeAddFriend.CallbackFailedContinue {
			callbackResp.ActionCode = constant.ActionForbidden
			return callbackResp
		} else {
			callbackResp.ActionCode = constant.ActionAllow
			return callbackResp
		}
	}
	return callbackResp
}
