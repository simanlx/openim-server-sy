package logic

import (
	cbApi "crazy_server/pkg/call_back_struct"
	"crazy_server/pkg/common/callback"
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/constant"
	"crazy_server/pkg/common/http"
	"crazy_server/pkg/common/log"
	pbChat "crazy_server/pkg/proto/msg"
	"crazy_server/pkg/utils"
	http2 "net/http"
)

func callbackAfterConsumeGroupMsg(msg []*pbChat.MsgDataToMQ, triggerID string) cbApi.CommonCallbackResp {
	callbackResp := cbApi.CommonCallbackResp{OperationID: triggerID}
	if !config.Config.Callback.CallbackAfterConsumeGroupMsg.Enable {
		return callbackResp
	}
	for _, v := range msg {
		if v.MsgData.SessionType == constant.SuperGroupChatType || v.MsgData.SessionType == constant.GroupChatType {
			commonCallbackReq := copyCallbackCommonReqStruct(v)
			commonCallbackReq.CallbackCommand = constant.CallbackAfterConsumeGroupMsgCommand
			req := cbApi.CallbackAfterConsumeGroupMsgReq{
				CommonCallbackReq: commonCallbackReq,
				GroupID:           v.MsgData.GroupID,
			}
			resp := &cbApi.CallbackAfterConsumeGroupMsgResp{CommonCallbackResp: &callbackResp}
			defer log.NewDebug(triggerID, utils.GetSelfFuncName(), req, *resp)
			if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackAfterConsumeGroupMsgCommand, req, resp, config.Config.Callback.CallbackAfterConsumeGroupMsg.CallbackTimeOut); err != nil {
				callbackResp.ErrCode = http2.StatusInternalServerError
				callbackResp.ErrMsg = err.Error()
				return callbackResp
			}
		}
	}

	log.NewDebug(triggerID, utils.GetSelfFuncName(), msg)

	return callbackResp
}
func copyCallbackCommonReqStruct(msg *pbChat.MsgDataToMQ) cbApi.CommonCallbackReq {
	req := cbApi.CommonCallbackReq{
		SendID:           msg.MsgData.SendID,
		ServerMsgID:      msg.MsgData.ServerMsgID,
		ClientMsgID:      msg.MsgData.ClientMsgID,
		OperationID:      msg.OperationID,
		SenderPlatformID: msg.MsgData.SenderPlatformID,
		SenderNickname:   msg.MsgData.SenderNickname,
		SessionType:      msg.MsgData.SessionType,
		MsgFrom:          msg.MsgData.MsgFrom,
		ContentType:      msg.MsgData.ContentType,
		Status:           msg.MsgData.Status,
		CreateTime:       msg.MsgData.CreateTime,
		AtUserIDList:     msg.MsgData.AtUserIDList,
		SenderFaceURL:    msg.MsgData.SenderFaceURL,
		Content:          callback.GetContent(msg.MsgData),
		Seq:              msg.MsgData.Seq,
		Ex:               msg.MsgData.Ex,
	}
	return req
}
