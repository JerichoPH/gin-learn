package middlewares

import (
	"gin-learn/controllers"
	"gin-learn/errors"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func JwtCheck(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokens := ctx.Request.Header["Token"]

		if len(tokens) == 0 {
			panic(errors.ThrowUnAuthorization("令牌不存在"))
		}
		token := tokens[0]

		ok := false

		if token == "" {
			panic(errors.ThrowUnAuthorization("令牌不存在"))
		} else {
			claims, err := controllers.ParseToken(token)

			// 判断令牌是否有效
			if err != nil {
				panic(errors.ThrowUnAuthorization("令牌解析失败"))
			} else if time.Now().Unix() > claims.ExpiresAt {
				panic(errors.ThrowUnAuthorization("令牌过期"))
			}

			// 判断用户是否存在
			accountController := controllers.AccountController{DB: *db}
			if accountController.FindByUsername(claims.Username).IsEmpty() {
				panic(errors.ThrowUnAuthorization("令牌解析失败：用户不存在"))
			}

			ctx.Set("__currentAccount", accountController.Account)
		}

		ok = true

		if !ok {
			ctx.Abort()
		}

		ctx.Next()
	}
}
