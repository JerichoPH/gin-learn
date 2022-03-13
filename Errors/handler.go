package Errors

import (
	"fmt"
	"gin-learn/Tools"
	"log"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func RecoverHandler(c *gin.Context) {
	defer func() {
		if reco := recover(); reco != nil {
			//打印错误堆栈信息
			log.Printf("panic: %v\n", reco)

			// 判断错误类型
			switch fmt.Sprintf("%T", reco) {
			case "*validator.ValidationErrors":
				// 表单验证错误
				c.JSON(Tools.ResponseINS().ErrValidate("", errorToString(reco)))
			case "*Errors.ForbiddenError":
				// 禁止操作
				c.JSON(Tools.ResponseINS().ErrForbidden(errorToString(reco)))
			case "*Errors.EmptyError":
				// 空数据
				c.JSON(Tools.ResponseINS().ErrEmpty(errorToString(reco)))
			case "*Errors.UnAuthorizationError":
				// 未授权
				c.JSON(Tools.ResponseINS().ErrUnAuthorization(errorToString(reco)))
			case "*Errors.UnLoginError":
				// 未登录
				c.JSON(Tools.ResponseINS().ErrUnLogin())
			default:
				// 其他错误
				c.JSON(Tools.ResponseINS().ErrAccident(errorToString(reco), reco))
				debug.PrintStack() // 打印堆栈信息
			}

			c.Abort()
		}
	}()

	c.Next()
}

// recover错误，转string
func errorToString(reco interface{}) string {
	switch errorType := reco.(type) {
	case error:
		return errorType.Error()
	default:
		return reco.(string)
	}
}
