package middleware

import (
	"crazy_server/pkg/common/log"
	"crazy_server/pkg/common/token_verify"
	"crazy_server/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, userID, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), "")

		// log.NewInfo("0", utils.GetSelfFuncName(), "userID: ", userID)
		c.Set("userID", userID)
		if !ok {
			log.NewError("", "GetUserIDFromToken false ", c.Request.Header.Get("token"))
			c.Abort()
			c.JSON(http.StatusOK, gin.H{"errCode": 403, "errMsg": errInfo})
			return
		} else {
			//if !utils.IsContain(userID, config.Config.Manager.AppManagerUid) {
			//	c.Abort()
			//	c.JSON(http.StatusOK, gin.H{"errCode": 400, "errMsg": "user is not admin"})
			//	return
			//}
			log.NewInfo("0", utils.GetSelfFuncName(), "failed: ", errInfo)
		}
	}
}
