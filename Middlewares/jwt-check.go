package Middlewares

import (
	"gin-learn/Controllers"
	"gin-learn/Errors"
	"time"

	"github.com/gin-gonic/gin"
)

func JwtCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header["Token"][0]

		ok := false

		if token == "" {
			panic(Errors.ThrowUnAuthorization("令牌不存在"))
		} else {
			claims, err := Controllers.ParseToken(token)
			if err != nil {
				panic(Errors.ThrowUnAuthorization("令牌解析失败"))
			} else if time.Now().Unix() > claims.ExpiresAt {
				panic(Errors.ThrowUnAuthorization("令牌过期"))
			}
		}

		ok = true

		if !ok {
			ctx.Abort()
		}

		ctx.Next()
	}
}
