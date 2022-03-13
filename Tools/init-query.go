package Tools

import (
	"gin-learn/Errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type QueryBuilder struct {
	CTX     gin.Context
	DB      gorm.DB
	Queries map[string]interface{}
}

func (cls *QueryBuilder) Init() *gorm.DB {
	tx := cls.DB.Where(cls.Queries)

	// 排序
	if order := cls.CTX.Query("order"); order != "" {
		tx.Order(order)
	}

	// offset
	if offset := cls.CTX.Query("offset"); offset != "" {
		if offset, err := strconv.Atoi(offset); err != nil {
			panic(Errors.ThrowForbidden("offset 参数只能填写数字"))
		} else {
			tx.Offset(offset)
		}
	}

	// limit
	if limit := cls.CTX.Query("limit"); limit != "" {
		if limit, err := strconv.Atoi(limit); err != nil {
			panic(Errors.ThrowForbidden("limit 参数只能填写数字"))
		} else {
			tx.Limit(limit)
		}
	}

	return tx
}
