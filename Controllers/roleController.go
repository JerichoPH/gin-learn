package controllers

import (
	"gin-learn/errors"
	"gin-learn/models"
	"gin-learn/tools"
	gin "github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleController struct {
	CTX *gin.Context
	DB  *gorm.DB
}

// Index 列表
func (cls *RoleController) Index() {
	roleModel := models.RoleModel{CTX: cls.CTX, DB: cls.DB}
	roles := roleModel.FindManyByQuery()
	cls.CTX.JSON(tools.CorrectIns().Ok("", gin.H{"roles": roles}))
}

// Show 详情
func (cls *RoleController) Show() {
	id := tools.GetID(cls.CTX.Param("id"))

	roleModel := models.RoleModel{CTX: cls.CTX, DB: cls.DB}
	role := roleModel.FindOneById(id)

	cls.CTX.JSON(tools.CorrectIns().Ok("", gin.H{"role": role}))
}

// Store 新建
func (cls *RoleController) Store() {
	role := (&models.RoleModel{CTX: cls.CTX, DB: cls.DB}).Store()
	cls.CTX.JSON(tools.CorrectIns().Created("", gin.H{"role": role}))
}

// Update 编辑
func (cls *RoleController) Update() {
	id := tools.GetID(cls.CTX.Param("id"))

	role := (&models.RoleModel{CTX: cls.CTX, DB: cls.DB}).UpdateOneById(id)
	cls.CTX.JSON(tools.CorrectIns().Updated("", gin.H{"role": role}))
}

// Destroy 删除
func (cls *RoleController) Destroy() {
	id := tools.GetID(cls.CTX.Param("id"))

	(&models.RoleModel{CTX: cls.CTX, DB: cls.DB}).DeleteOneById(id)
	cls.CTX.JSON(tools.CorrectIns().Deleted(""))
}

// BindAccounts 绑定用户
func (cls *RoleController) PostBindAccounts() {
	roleID := cls.CTX.Param("id")
	if accountIDs, err := cls.CTX.GetQueryArray("account_ids"); !err {
		panic(errors.ThrowForbidden("获取参数失败"))
	} else {
		(&models.RoleAccountModel{CTX:cls.CTX,DB:cls.DB}).StoreBatch(roleID,accountIDS)
	}
}
