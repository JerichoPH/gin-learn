package tools

import "sync"

type correct struct {
	msg       string
	content   interface{}
	status    int
	errorCode int
}

var responseIns *correct
var once sync.Once

func CorrectIns() *correct {
	once.Do(func() { responseIns = &correct{} })
	return responseIns
}

func (cls *correct) Get() map[string]interface{} {
	ret := map[string]interface{}{
		"msg":      cls.msg,
		"content":  cls.content,
		"status":   cls.status,
		"erroCode": cls.errorCode,
	}
	return ret
}

func (cls *correct) Set(msg string, content interface{}, status int, errorCode int) *correct {
	cls.msg = msg
	cls.content = content
	if status == 0 {
		cls.status = 200
	} else {
		cls.status = status
	}
	cls.errorCode = errorCode
	return cls
}

func (cls *correct) Ok(msg string, content interface{}) (int, map[string]interface{}) {
	if msg == "" {
		msg = "OK"
	}

	return 200, cls.Set(msg, content, 200, 0).Get()
}

func (cls *correct) Created(msg string, content interface{}) (int, map[string]interface{}) {
	if msg == "" {
		msg = "新建成功"
	}

	return 201, cls.Set(msg, content, 201, 0).Get()
}

func (cls *correct) Updated(msg string, content interface{}) (int, map[string]interface{}) {
	if msg == "" {
		msg = "编辑成功"
	}

	return 202, cls.Set(msg, content, 202, 0).Get()
}

func (cls *correct) Deleted(msg string) (int, interface{}) {
	if msg == "" {
		msg = "删除成功"
	}

	return 204, cls.Set(msg, nil, 204, 0).Get()
}
