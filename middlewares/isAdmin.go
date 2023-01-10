package middlewares

import (
	"net/http"
	"web_app/models"

	"github.com/gin-gonic/gin"
)

func IsAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*models.CustomClaims)
		if currentUser.AuthorityId != 3 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "对不起，只有管理员用户可以访问此页面",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
