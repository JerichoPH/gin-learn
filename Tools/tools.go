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
func IsEmpty(ins interface{}, class interface{}, name string) (isEmpty bool) {
	isEmpty = reflect.DeepEqual(ins, class)

	if name != "" {
		if isEmpty {
			panic(errors.ThrowEmpty(fmt.Sprintf("%v不存在", name)))
		}
	}

	return isEmpty
}

// IsRepeat 判断是否重复
func IsRepeat(ins interface{}, class interface{}, name string) (isRepeat bool) {
	isRepeat = !reflect.DeepEqual(ins, class)

	if name != "" {
		if isRepeat {
			panic(errors.ThrowForbidden(fmt.Sprintf("%v重复", name)))
		}
	}

	return
}
