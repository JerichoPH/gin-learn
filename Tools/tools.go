package tools

import (
	"fmt"
	"gin-learn/errors"
	"reflect"
	"strconv"
)

// GetID 获取ID
func GetID(v string) int {
	id, err := strconv.Atoi(v)
	if err != nil {
		panic(errors.ThrowForbidden("id必须是数字"))
	}

	return id
}

// IsEmpty 判断是否为空
func IsEmpty(ins interface{}, class interface{}, name string) {
	if reflect.DeepEqual(ins, class) {
		panic(errors.ThrowEmpty(fmt.Sprintf("%v不存在", name)))
	}
}

// IsRepeat 判断是否重复
func IsRepeat(ins interface{}, class interface{}, name string) {
	if !reflect.DeepEqual(ins, class) {
		panic(errors.ThrowForbidden(fmt.Sprintf("%v重复", name)))
	}
}
