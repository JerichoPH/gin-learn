package v1

import (
	"gin-learn/controllers"
	"gin-learn/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccountRouter struct {
	Router *gin.Engine
	DB     *gorm.DB
}

func (cls *AccountRouter) Load() {
	r := cls.Router.Group("/v1", middlewares.JwtCheck(cls.DB))
	{
		// 用户列表
		r.GET("/account", func(ctx *gin.Context) {
			(&controllers.AccountController{CTX: ctx, DB: cls.DB}).Index()
		})

		// 根据id获取用户详情
		r.GET("/account/:id", func(ctx *gin.Context) {
			(&controllers.AccountController{CTX: ctx, DB: cls.DB}).Show()
		})
	}
}
