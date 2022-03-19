package models

import (
	"gin-learn/errors"
	"gin-learn/tools"
	"reflect"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Status struct {
	gorm.Model
	Name string `form:"name" binding:"required" gorm:"type=VARCHAR(64);unique;NOT NULL;comment '状态名称'"`
}

type StatusModel struct {
	CTX *gin.Context
	DB  *gorm.DB
}

// Store 新建
func (cls *StatusModel) Store() Status {
	var status Status
	if err := cls.CTX.ShouldBind(&status); err != nil {
		panic(err)
	}

	var repeatStatus Status
	cls.DB.Where(map[string]interface{}{"name": status.Name}).First(&repeatStatus)
	if !reflect.DeepEqual(repeatStatus, Status{}) {
		panic(errors.ThrowForbidden("状态名称重复"))
	}

	cls.DB.Create(status)
	return status
}

// DeleteById 根据id删除
func (cls *StatusModel) DeleteById(id int) *StatusModel {
	cls.DB.Delete(&Status{}, id)

	return cls
}

// UpdateById 根据id编辑
func (cls *StatusModel) UpdateById(id int) *StatusModel {
	status := cls.FindOneById(id)

	var statusForm Status
	if err := cls.CTX.ShouldBind(&statusForm); err != nil {
		panic(err)
	}

	var repeatStatus Status
	cls.DB.Where(map[string]interface{}{"name": statusForm.Name}).Not(map[string]interface{}{"id": id}).First(&repeatStatus)
	if !reflect.DeepEqual(repeatStatus, Status{}) {
		panic(errors.ThrowForbidden("状态名称重复"))
	}

	status.Name = statusForm.Name
	cls.DB.Save(status)

	return cls
}

// FindMoreByQuery 根据Query读取用户列表
func (cls *StatusModel) FindMoreByQuery() []Status {

	var statuses []Status
	w := make(map[string]interface{})
	n := make(map[string]interface{})

	if name := cls.CTX.Query("name"); name != "" {
		w["name"] = name
	}

	(&tools.QueryBuilder{CTX: cls.CTX, DB: cls.DB}).Init(w, n).Find(&statuses)

	return statuses
}

// FindOneById 根据编号搜索
func (cls *StatusModel) FindOneById(id int) Status {
	var status Status
	cls.DB.Where(map[string]interface{}{"id": id}).First(&status)

	if reflect.DeepEqual(status, Status{}) {
		panic(errors.ThrowEmpty(""))
	}

	return status
}
