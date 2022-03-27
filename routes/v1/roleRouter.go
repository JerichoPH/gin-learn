package v1

import (
	"gin-learn/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleRouter struct {
	Router *gin.Engine
	DB     *gorm.DB
}

func (cls *RoleRouter) Load() {
	r := cls.Router.Group("/v1")
	{
		// 列表
		r.GET("/role", func(ctx *gin.Context) {
			(&controllers.RoleController{CTX: ctx, DB: cls.DB}).Index()
		})

		// 详情
		r.GET("/role/:id", func(ctx *gin.Context) {
			(&controllers.RoleController{CTX: ctx, DB: cls.DB}).Show()
		})

		// 新建
		r.POST("/role", func(ctx *gin.Context) {
			(&controllers.RoleController{CTX: ctx, DB: cls.DB}).Store()
		})

		// 编辑
		r.PUT("/role/:id", func(ctx *gin.Context) {
			(&controllers.RoleController{CTX: ctx, DB: cls.DB}).Update()
		})

		// 删除
		r.DELETE("/role/:id", func(ctx *gin.Context) {
			(&controllers.RoleController{CTX: ctx, DB: cls.DB}).Destroy()
		})

		// 绑定用户
		r.POST("/role/:id/bindAccounts", func(ctx *gin.Context) {
			(&controllers.RoleController{CTX: ctx, DB: cls.DB}).PostBindAccounts()
		})
	}
}
