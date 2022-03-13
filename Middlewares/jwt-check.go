package Middlewares

import (
	"gin-learn/Controllers"
	"gin-learn/Errors"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func JwtCheck(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokens := ctx.Request.Header["Token"]

		if len(tokens) == 0 {
			panic(Errors.ThrowUnAuthorization("令牌不存在"))
		}
		token := tokens[0]

		ok := false

		if token == "" {
			panic(Errors.ThrowUnAuthorization("令牌不存在"))
		} else {
			claims, err := Controllers.ParseToken(token)

			// 判断令牌是否有效
			if err != nil {
				panic(Errors.ThrowUnAuthorization("令牌解析失败"))
			} else if time.Now().Unix() > claims.ExpiresAt {
				panic(Errors.ThrowUnAuthorization("令牌过期"))
			}

			// 判断用户是否存在
			accountController := Controllers.AccountController{DB: *db}
			if accountController.FindByUsername(claims.Username).IsEmpty() {
				panic(Errors.ThrowUnAuthorization("用户不存在"))
			}

			ctx.Set("__currentAccount", accountController.GetAccount())
		}

		ok = true

		if !ok {
			ctx.Abort()
		}

		ctx.Next()
	}
}
