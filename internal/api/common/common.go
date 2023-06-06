package common

import (
	"crazy_server/pkg/common/log"
	"crazy_server/pkg/common/token_verify"
	"crazy_server/pkg/proto/agent"
	"crazy_server/pkg/proto/cloud_wallet"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 处理错误返回
func HandleCommonRespErr(commResp *cloud_wallet.CommonResp, c *gin.Context) bool {
	if commResp != nil && commResp.ErrCode != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": commResp.ErrCode, "errMsg": commResp.ErrMsg})
		return true
	}
	return false
}

// 解析im token、获取用户id
func ParseImToken(c *gin.Context, operationID string) (string, bool) {
	ok, userId, errMsg := token_verify.GetUserIDFromToken(c.Request.Header.Get("im-token"), operationID)
	if !ok {
		log.NewError(operationID, "im-token parseToken err :", errMsg)
		c.JSON(http.StatusForbidden, gin.H{"errCode": 403, "errMsg": "token授权认证失败"})
		return "", false
	}
	return userId, true
}

// 推广系统处理错误返回
func HandleAgentCommonRespErr(commResp *agent.CommonResp, c *gin.Context) bool {
	if commResp != nil && commResp.Code != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": commResp.Code, "msg": commResp.Msg})
		return true
	}
	return false
}
