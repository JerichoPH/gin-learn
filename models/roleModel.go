package models

import (
	"gin-learn/errors"
	"gin-learn/tools"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
)

type Role struct {
	gorm.Model
	Name     string     `gorm:"type:VARCHAR(64);unique;NOT NULL;COMMENT:'角色名称';"`
	Accounts []*Account `gorm:"many2many:role_accounts;"`
}

type RoleModel struct {
	CTX *gin.Context
	DB  *gorm.DB
}

// Store 新建
func (cls *RoleModel) Store() Role {
	var role Role
	if err := cls.CTX.ShouldBind(&role); err != nil {
		panic(err)
	}

	var repeatRole Role
	cls.DB.Where(map[string]interface{}{"name": role.Name}).First(&repeatRole)
	if !reflect.DeepEqual(repeatRole, Role{}) {
		panic(errors.ThrowForbidden("角色名称重复"))
	}

	cls.DB.Omit(clause.Associations).Create(&role)

	return role
}

// DeleteOneById 根据id删除
func (cls *RoleModel) DeleteOneById(id int) *RoleModel {
	role := cls.FindOneById(id)
	cls.DB.Omit(clause.Associations).Delete(&role)

	return cls
}

// UpdateOneById 根据id编辑
func (cls *RoleModel) UpdateOneById(id int) Role {
	role := cls.FindOneById(id)

	var roleForm Role
	if err := cls.CTX.ShouldBind(&roleForm); err != nil {
		panic(err)
	}

	var repeatRole Role
	cls.DB.Where(map[string]interface{}{"name": roleForm.Name}).Not(map[string]interface{}{"id": id}).First(&repeatRole)
	if !reflect.DeepEqual(repeatRole, Role{}) {
		panic(errors.ThrowForbidden("角色名称重复"))
	}

	role.Name = roleForm.Name
	cls.DB.Omit(clause.Associations).Save(&role)

	return role
}

// FindOneById 根据编号查询
func (cls *RoleModel) FindOneById(id int, preloads ...string) Role {
	var role Role
	query := cls.DB.Where(map[string]interface{}{"id": id})
	if preloads != nil {
		for _, preload := range preloads {
			query.Preload(preload)
		}
	}
	query.Preload(clause.Associations).First(&role)

	tools.IsEmpty(role, Role{}, "角色")

	return role
}

// FindManyByQuery 根据query参数获取列表
func (cls *RoleModel) FindManyByQuery(preloads ...string) []Role {
	var roles []Role

	// 搜索条件
	w := make(map[string]interface{})
	n := make(map[string]interface{})

	query := (&tools.QueryBuilder{CTX: cls.CTX, DB: cls.DB}).Init(w, n)
	if preloads != nil {
		for _, preload := range preloads {
			query.Preload(preload)
		}
	}
	if name := cls.CTX.Query("name"); name != "" {
		query.Where("name LIKE '%?%'", name)
	}
	query.Preload(clause.Associations).Find(&roles)

	return roles
}
