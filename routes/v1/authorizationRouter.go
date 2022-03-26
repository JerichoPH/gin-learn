package v1

import (
	"gin-learn/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthorizationRouter struct {
	Router *gin.Engine
	DB     *gorm.DB
}

func (cls *AuthorizationRouter) Load() {
	// 注册
	cls.Router.POST("/v1/authorization/register", func(ctx *gin.Context) {
		(&controllers.AuthorizationController{CTX: ctx, DB: cls.DB}).PostRegister()
	})

	// 登录
	cls.Router.POST("/v1/authorization/login", func(ctx *gin.Context) {
		(&controllers.AuthorizationController{CTX: ctx, DB: cls.DB}).PostLogin()
	})
}
