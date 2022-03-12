package Tools

import "sync"

type Map map[string]interface{}

type response struct {
	msg       string
	content   interface{}
	status    int
	errorCode int
}

var responseIns *response
var once sync.Once

func ResponseINS() *response {
	once.Do(func() { responseIns = &response{} })
	return responseIns
}

func (cls *response) Get() Map {
	ret := Map{
		"msg":      cls.msg,
		"content":  cls.content,
		"status":   cls.status,
		"erroCode": cls.errorCode,
	}
	return ret
}

func (cls *response) Set(msg string, content interface{}, status int, errorCode int) *response {
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

func (cls *response) Ok(msg string, content interface{}) (int, Map) {
	if msg == "" {
		msg = "OK"
	}

	return 200, cls.Set(msg, content, 200, 0).Get()
}

func (cls *response) OkCreated(msg string, content interface{}) (int, Map) {
	if msg == "" {
		msg = "新建成功"
	}

	return 201, cls.Set(msg, content, 201, 0).Get()
}

func (cls *response) OkUpdated(msg string, content interface{}) (int, Map) {
	if msg == "" {
		msg = "编辑成功"
	}

	return 202, cls.Set(msg, content, 202, 0).Get()
}

func (cls *response) OkDeleted(msg string) (int, interface{}) {
	if msg == "" {
		msg = "删除成功"
	}

	return 204, cls.Set(msg, nil, 204, 0).Get()
}

func (cls *response) ErrUnAuthorization() (int, interface{}) {
	return 406, cls.Set("未授权", nil, 406, 1).Get()
}

func (cls *response) ErrUnLogin() (int, Map) {
	return 401, cls.Set("未登录", nil, 401, 2).Get()
}

func (cls *response) ErrForbidden(msg string) (int, interface{}) {
	if msg == "" {
		msg = "禁止操作"
	}

	return 403, cls.Set(msg, nil, 403, 3).Get()
}

func (cls *response) ErrEmpty(msg string) (int, interface{}) {
	if msg == "" {
		msg = "数不存在"
	}

	return 404, cls.Set(msg, nil, 404, 4).Get()
}

func (cls *response) ErrValidate(msg string, content interface{}) (int, Map) {
	if msg == "" {
		msg = "表单验证错误"
	}

	return 421, cls.Set(msg, content, 421, 5).Get()
}

func (cls *response) ErrAccident(err interface{}) (int, Map) {
	return 500, cls.Set("意外错误", err, 500, 6).Get()
}
