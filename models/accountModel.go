package models

import (
	"database/sql"
	"gin-learn/errors"
	"gin-learn/tools"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
)

// Account 用户表
type Account struct {
	gorm.Model
	Username         string       `form:"username" gorm:"type:VARCHAR(64);NOT NULL;unique_index:users__username__uidx;COMMENT:'用户名';"`
	Password         string       `form:"password" gorm:"type:VARCHAR(128);NOT NULL;comment:'密码';"`
	Nickname         string       `form:"nickname" gorm:"type:VARCHAR(64);NOT NULL;DEFAULT '';index:users__nickname__idx;COMMENT:'昵称';"`
	ActivatedAt      sql.NullTime `gorm:"type:DATETIME;COMMENT:'激活时间';"`
	StatusUniqueCode string       `gorm:"type:VARCHAR(64);NOT NULL;DEFAULT '';COMMENT:'状态代码';"`
	Status           Status       `gorm:"references:UniqueCode;foreignKey:StatusUniqueCode;"`
}

// AccountModel 用户模型
type AccountModel struct {
	CTX *gin.Context
	DB  *gorm.DB
}

// ScopeActivated 已经激活的
func (cls *AccountModel) ScopeActivated(db *gorm.DB) *gorm.DB {
	return db.Not(map[string]interface{}{"activated_at": nil})
}

// ScopeNotActivated 未激活的
func (cls *AccountModel) ScopeNotActivated(db *gorm.DB) *gorm.DB {
	return db.Where(map[string]interface{}{"activated_at": nil})
}

// FindOneById 根据id获取用户
func (cls *AccountModel) FindOneById(id int) Account {
	var account Account
	cls.DB.Where(map[string]interface{}{"id": id}).Find(account)
	return account
}

// FindOneByUsername 根据用户名读取用户
func (cls *AccountModel) FindOneByUsername(username string) Account {
	if username == "" {
		panic(errors.ThrowForbidden("用户名不能为空"))
	}

	var account Account
	cls.DB.Scopes(cls.ScopeActivated).Where(map[string]interface{}{"username": username}).Find(&account)
	return account
}

// FindMoreByQuery 根据Query读取用户列表
func (cls *AccountModel) FindMoreByQuery() []Account {
	var (
		account  Account
		accounts []Account
	)

	if err := cls.CTX.ShouldBindQuery(&account); err != nil {
		panic(err)
	}

	// 搜索条件
	w := make(map[string]interface{})
	n := make(map[string]interface{})
	if username := cls.CTX.Query("username"); username != "" {
		w["username"] = username
	}
	if nickname := cls.CTX.Query("nickname"); nickname != "" {
		w["nickname"] = nickname
	}
	if activatedAt := cls.CTX.Query("activated_at"); activatedAt != "" {
		switch activatedAt {
		case "1":
			n["activated_at"] = nil
		case "0":
			w["activated_at"] = nil
		}
	}

	query := (&tools.QueryBuilder{CTX: cls.CTX, DB: cls.DB}).Init(w, n)
	if activatedAtBetween := cls.CTX.Query("activated_at_between"); activatedAtBetween != "" {
		query.Where("activated_at BETWEEN ? AND ?", strings.Split(activatedAtBetween, "~"))
	}

	query.Find(accounts)

	return accounts
}
