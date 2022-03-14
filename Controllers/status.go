package controllers

import (
	"gin-learn/errors"
	"gin-learn/tools"
	"reflect"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Status struct {
	gorm.Model
	Name string `form:"name" gorm:"type=VARCHAR(64);unique;NOT NULL;comment '状态名称'"`
}

type StatusController struct {
	CTX gin.Context
	DB  gorm.DB
	Status
	Statuses []Status
}

// 绑定新建表单
func (cls *StatusController) BindFormStore() *StatusController {
	var status Status
	if err := cls.CTX.ShouldBind(&status); err != nil {
		panic(err)
	}

	cls.Status = status
	return cls
}

// 新建
func (cls *StatusController) Store() *StatusController {
	var repeatStauts Status
	cls.DB.Where(map[string]interface{}{"name": cls.Status.Name}).First(&repeatStauts)

	if !reflect.DeepEqual(repeatStauts, Status{}) {
		panic(errors.ThrowForbidden("状态名称重复"))
	}

	cls.DB.Create(&cls.Status)

	return cls
}

// 根据Query读取用户列表
func (cls *StatusController) FindMoreByQuery() *StatusController {

	w := make(map[string]interface{})
	n := make(map[string]interface{})

	if name := cls.CTX.Query("name"); name != "" {
		w["name"] = name
	}

	(&tools.QueryBuilder{CTX: cls.CTX, DB: cls.DB}).Init(w, n).Find(&cls.Statuses)

	return cls
}
