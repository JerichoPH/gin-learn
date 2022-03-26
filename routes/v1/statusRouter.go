package v1

import (
	"gin-learn/controllers"
	"gin-learn/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StatusRouter struct {
	Router *gin.Engine
	DB     *gorm.DB
}

func (cls *StatusRouter) Load() {
	r := cls.Router.Group("/v1", middlewares.JwtCheck(cls.DB))
	{
		// 状态列表
		r.GET("/status", func(ctx *gin.Context) {
			(&controllers.StatusController{CTX: ctx, DB: cls.DB}).Index()
		})

		// 状态详情
		r.GET("/status/:id", func(ctx *gin.Context) {
			(&controllers.StatusController{CTX: ctx, DB: cls.DB}).Show()
		})

		// 新建状态
		r.POST("/status", func(ctx *gin.Context) {
			(&controllers.StatusController{CTX: ctx, DB: cls.DB}).Store()
		})

		// 编辑状态
		r.PUT("/status/:id", func(ctx *gin.Context) {
			(&controllers.StatusController{CTX: ctx, DB: cls.DB}).Update()
		})

		// 删除状态
		r.DELETE("/status/:id", func(ctx *gin.Context) {
			(&controllers.StatusController{CTX: ctx, DB: cls.DB}).Destroy()
		})
	}
}
