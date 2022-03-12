package Tools

import (
	"gin-learn/Errors"
	"log"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Recover(c *gin.Context) {
	defer func() {
		if reco := recover(); reco != nil {
			//打印错误堆栈信息
			log.Printf("panic: %v\n", reco)
			debug.PrintStack()

			// 判断错误类型
			switch reco.(type) {
			case validator.ValidationErrors:
				// 表单验证错误
				c.JSON(ResponseINS().ErrValidate("", errorToString(reco)))
			case Errors.ForbiddenError:
				// 禁止操作
				c.JSON(ResponseINS().ErrForbidden(errorToString(reco)))
			case Errors.EmptyError:
				// 空数据
				c.JSON(ResponseINS().ErrEmpty(errorToString(reco)))
			case Errors.UnAuthorizationError:
				// 未授权
				c.JSON(ResponseINS().ErrUnAuthorization())
			case Errors.UnLoginError:
				// 未登录
				c.JSON(ResponseINS().ErrUnLogin())
			default:
				// 其他错误
				c.JSON(ResponseINS().ErrAccident(errorToString(reco), reco))
			}

			c.Abort()
		}
	}()

	c.Next()
}

// recover错误，转string
func errorToString(reco interface{}) string {
	switch v := reco.(type) {
	case error:
		return v.Error()
	default:
		return reco.(string)
	}
}
