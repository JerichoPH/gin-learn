package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleAccount struct {
	RoleID    uint
	AccountID uint
}

type RoleAccountModel struct {
	CTX *gin.Context
	DB  *gorm.DB
}

// Store 新建
func (cls *RoleAccountModel) Store(roleID, accountID uint) (roleAccount RoleAccount) {
	cls.DB.Where(map[string]interface{}{"role_id": roleID, "account_id": accountID}).FirstOrCreate(&roleAccount)
	return
}

func (cls *RoleAccountModel) StoreBatch(roleID uint, accountIDs []uint) []RoleAccount {
	//roleAccounts := make([]RoleAccount, 200)
	var roleAccounts []RoleAccount

	for _, accountID := range accountIDs {
		roleAccount := &RoleAccount{RoleID: roleID, AccountID: accountID}
		roleAccounts2 := append(roleAccounts, roleAccount)

	}
}
