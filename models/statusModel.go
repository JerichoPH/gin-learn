package models

import (
	"gin-learn/errors"
	"gin-learn/tools"
	"gorm.io/gorm/clause"
	"reflect"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Status struct {
	gorm.Model
	UniqueCode string    `gorm:"type:VARCHAR(64);unique;NOT NULL;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;COMMENT:'状态代码';"`
	Name       string    `form:"name" binding:"required" gorm:"type:VARCHAR(64);unique;NOT NULL;COMMENT:'状态名称';"`
	Accounts   []Account `gorm:"foreignKey:StatusUniqueCode;references:UniqueCode;"`
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
	tools.IsRepeat(repeatStatus, Status{}, "状态名称")

	cls.DB.Omit(clause.Associations).Create(status)
	return status
}

// DeleteById 根据id删除
func (cls *StatusModel) DeleteById(id int) *StatusModel {
	cls.DB.Delete(&Status{}, id)

	return cls
}

// UpdateById 根据id编辑
func (cls *StatusModel) UpdateById(id int) Status {
	status := cls.FindOneById(id, "Accounts", "Accounts.Status")

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
	cls.DB.Omit(clause.Associations).Save(&status)

	return status
}

// FindOneById 根据编号搜索
func (cls *StatusModel) FindOneById(id int, preloads ...string) Status {
	var status Status
	query := cls.DB.Preload(clause.Associations).Where(map[string]interface{}{"id": id})
	if preloads != nil {
		for _, preload := range preloads {
			query.Preload(preload)
		}
	}
	query.First(&status)

	tools.IsEmpty(status, Status{}, "状态")

	return status
}

// FindManyByQuery 根据Query读取用户列表
func (cls *StatusModel) FindManyByQuery(preloads ...string) []Status {
	var statuses []Status
	w := make(map[string]interface{})
	n := make(map[string]interface{})

	if name := cls.CTX.Query("name"); name != "" {
		w["name"] = name
	}

	query := (&tools.QueryBuilder{CTX: cls.CTX, DB: cls.DB}).Init(w, n)
	if preloads != nil {
		for _, preload := range preloads {
			query.Preload(preload)
		}
	}
	if name := cls.CTX.Query("name"); name != "" {
		query.Where("`name` like '%?%'", name)
	}
	query.Find(&statuses)

	return statuses
}
