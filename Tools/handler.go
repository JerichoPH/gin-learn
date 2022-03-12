package Tools

import (
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
			default:
				// 意外错误
				c.JSON(ResponseINS().ErrAccident(reco))
			}

			c.Abort()
		}
	}()

	c.Next()
}

// recover错误，转string
func errorToString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	default:
		return r.(string)
	}
}
