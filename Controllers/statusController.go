package controllers

import (
	"gin-learn/errors"
	"gin-learn/models"
	"gin-learn/tools"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

// StatusController 状态控制器
type StatusController struct {
	CTX *gin.Context
	DB  *gorm.DB
}

// Index 列表
func (cls *StatusController) Index() {
	statusModel := &models.StatusModel{CTX: cls.CTX, DB: cls.DB}
	statuses := statusModel.FindManyByQuery("Accounts", "Accounts.Status")
	cls.CTX.JSON(tools.CorrectIns().Ok("", gin.H{"statuses": statuses}))
}

// Show 详情
func (cls *StatusController) Show() {
	statusModel := &models.StatusModel{CTX: cls.CTX, DB: cls.DB}

	id, err := strconv.Atoi(cls.CTX.Param("id"))
	if err != nil {
		panic(errors.ThrowForbidden("id必须是数字"))
	}

	status := statusModel.FindOneById(id, "Accounts", "Accounts.Status")
	cls.CTX.JSON(tools.CorrectIns().Ok("", gin.H{"status": status}))
}

// Store 新建
func (cls *StatusController) Store() {
	statusModel := &models.StatusModel{CTX: cls.CTX, DB: cls.DB}
	status := statusModel.Store()
	cls.CTX.JSON(tools.CorrectIns().Created("", gin.H{"status": status}))
}

// Update 更新
func (cls *StatusController) Update() {
	statusModel := &models.StatusModel{CTX: cls.CTX, DB: cls.DB}

	id, err := strconv.Atoi(cls.CTX.Param("id"))
	if err != nil {
		panic(errors.ThrowForbidden("id必须是数字"))
	}

	status := statusModel.UpdateById(id)
	cls.CTX.JSON(tools.CorrectIns().Updated("", gin.H{"status": status}))
}

// Destroy 删除
func (cls *StatusController) Destroy() {
	statusModel := &models.StatusModel{CTX: cls.CTX, DB: cls.DB}

	id, err := strconv.Atoi(cls.CTX.Param("id"))
	if err != nil {
		panic(errors.ThrowForbidden("id必须是数字"))
	}

	statusModel.DeleteById(id)
	cls.CTX.JSON(tools.CorrectIns().Deleted(""))
}
