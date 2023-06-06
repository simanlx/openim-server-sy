package logic

import (
	tpns "crazy_server/internal/push/sdk/tpns-server-sdk-go/go"
	"crazy_server/internal/push/sdk/tpns-server-sdk-go/go/auth"
	"crazy_server/internal/push/sdk/tpns-server-sdk-go/go/common"
	"crazy_server/internal/push/sdk/tpns-server-sdk-go/go/req"
	"crazy_server/pkg/common/config"
)

var badgeType = -2
var iosAcceptId = auth.Auther{AccessID: config.Config.Push.Tpns.Ios.AccessID, SecretKey: config.Config.Push.Tpns.Ios.SecretKey}

func IOSAccountListPush(accounts []string, title, content, jsonCustomContent string) {
	var iosMessage = tpns.Message{
		Title:   title,
		Content: content,
		IOS: &tpns.IOSParams{
			Aps: &tpns.Aps{
				BadgeType: &badgeType,
				Sound:     "default",
				Category:  "INVITE_CATEGORY",
			},
			CustomContent: jsonCustomContent,
			//CustomContent: `"{"key\":\"value\"}"`,
		},
	}
	pushReq, reqBody, err := req.NewListAccountPush(accounts, iosMessage)
	if err != nil {
		return
	}
	iosAcceptId.Auth(pushReq, auth.UseSignAuthored, iosAcceptId, reqBody)
	common.PushAndGetResult(pushReq)
}
