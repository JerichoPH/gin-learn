package controllers

import (
	"gin-learn/errors"
	"gin-learn/tools"
	"log"
	"reflect"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Status struct {
	gorm.Model
	Name string `form:"name" binding:"required" gorm:"type=VARCHAR(64);unique;NOT NULL;comment '状态名称'"`
}

type StatusController struct {
	CTX *gin.Context
	DB  *gorm.DB
	Status
	Statuses []Status
}

// Store 新建
func (cls *StatusController) Store() *StatusController {
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
	cls.Status = status
	return cls
}

// DeleteById 根据id删除
func (cls *StatusController) DeleteById(id int) *StatusController {
	cls.FindOneById(id)

	log.Println(cls.Status)

	cls.DB.Delete(&cls.Status)

	return cls
}

// UpdateById 根据id编辑
func (cls *StatusController) UpdateById(id int) *StatusController {
	cls.FindOneById(id)

	var status Status
	if err := cls.CTX.ShouldBind(&status); err != nil {
		panic(err)
	}

	var repeatStatus Status
	cls.DB.Where(map[string]interface{}{"name": status.Name}).Not(map[string]interface{}{"id": id}).First(&repeatStatus)
	if !reflect.DeepEqual(repeatStatus, Status{}) {
		panic(errors.ThrowForbidden("状态名称重复"))
	}

	cls.Status.Name = status.Name
	cls.DB.Save(cls.Status)

	return cls
}

// FindMoreByQuery 根据Query读取用户列表
func (cls *StatusController) FindMoreByQuery() *StatusController {

	w := make(map[string]interface{})
	n := make(map[string]interface{})

	if name := cls.CTX.Query("name"); name != "" {
		w["name"] = name
	}

	(&tools.QueryBuilder{CTX: cls.CTX, DB: cls.DB}).Init(w, n).Find(&cls.Statuses)

	return cls
}

// FindOneById 根据编号搜索
func (cls *StatusController) FindOneById(id int) *StatusController {
	cls.DB.Where(map[string]interface{}{"id": id}).First(&cls.Status)

	if reflect.DeepEqual(cls.Status, Status{}) {
		panic(errors.ThrowEmpty(""))
	}

	return cls
}
