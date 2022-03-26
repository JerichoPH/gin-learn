package middlewares

import (
	"gin-learn/controllers"
	"gin-learn/errors"
	"gin-learn/models"
	"gorm.io/gorm"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
)

func JwtCheck(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var account models.Account

		tokens := ctx.Request.Header["Token"]

		if len(tokens) == 0 {
			panic(errors.ThrowUnAuthorization("令牌不存在"))
		}
		token := tokens[0]

		ok := false

		if token == "" {
			panic(errors.ThrowUnAuthorization("令牌不存在"))
		} else {
			claims, err := (&controllers.AuthorizationController{}).ParseJwt(token)

			// 判断令牌是否有效
			if err != nil {
				panic(errors.ThrowUnAuthorization("令牌解析失败"))
			} else if time.Now().Unix() > claims.ExpiresAt {
				panic(errors.ThrowUnAuthorization("令牌过期"))
			}

			// 判断用户是否存在
			if reflect.DeepEqual(claims, controllers.Claims{}) {
				panic(errors.ThrowUnAuthorization("令牌解析失败：用户不存在"))
			}

			// 获取用户信息
			account = (&models.AccountModel{CTX: ctx, DB: db}).FindOneByUsername(claims.Username)
			if reflect.DeepEqual(account, models.Account{}) {
				panic(errors.ThrowEmpty("用户不存在"))
			}
		}

		ctx.Set("__currentAccount", account)
		ok = true
		if !ok {
			ctx.Abort()
		}

		ctx.Next()
	}
}
