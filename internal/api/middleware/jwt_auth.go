package middleware

import (
	rocksCache "crazy_server/pkg/common/db/rocks_cache"
	"crazy_server/pkg/common/log"
	"crazy_server/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		operationId := utils.OperationIDGenerator()
		unionId := c.Request.Header.Get("token")
		userId := rocksCache.GetUsersUnionIdFromCache(unionId)
		//ok, userId, _ := token_verify.GetUserIDFromToken(c.Request.Header.Get("unionid"), operationId)
		if userId == "" {
			log.NewError("", "GetUsersUnionIdFromCache false ", unionId)
			c.Abort()
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "token认证失败"})
			return
		} else {
			// 用户id
			c.Set("userId", userId)

			//operationID
			c.Set("operationId", operationId)
		}
	}
}
