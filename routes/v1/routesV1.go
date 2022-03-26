package v1

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoutesV1 struct {
	Router *gin.Engine
	DB     *gorm.DB
}

func (cls *RoutesV1) Load() {
	(&AuthorizationRouter{Router: cls.Router, DB: cls.DB}).Load() // 授权管理路由
	(&AccountRouter{Router: cls.Router, DB: cls.DB}).Load()       // 用户管理路由
	(&StatusRouter{Router: cls.Router, DB: cls.DB}).Load()        // 状态管理路由
}
