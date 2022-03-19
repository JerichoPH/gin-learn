package controllers

import (
	"gin-learn/errors"
	"gin-learn/models"
	"gin-learn/tools"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"reflect"
	"strconv"
)

// AccountController 用户控制器
type AccountController struct {
	CTX *gin.Context
	DB  *gorm.DB
}

// Index 列表
func (cls *AccountController) Index() {
	accountModel := &models.AccountModel{CTX: cls.CTX, DB: cls.DB}
	accounts := accountModel.FindMoreByQuery()
	cls.CTX.JSON(tools.CorrectIns().Ok("", gin.H{"accounts": accounts}))
}

// Show 详情
func (cls *AccountController) Show() {
	id, err := strconv.Atoi(cls.CTX.Param("id"))
	if err != nil {
		panic(errors.ThrowForbidden("id必须是数字"))
	}

	accountModel := &models.AccountModel{CTX: cls.CTX, DB: cls.DB}
	account := accountModel.FindOneById(id)
	if reflect.DeepEqual(account, models.Account{}) {
		panic(errors.ThrowEmpty("用户不存在"))
	}
	cls.CTX.JSON(tools.CorrectIns().Ok("", gin.H{"accounts": account}))
}
