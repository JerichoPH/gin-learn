package errors

import "sync"

type inCorrect struct {
	msg       string
	content   interface{}
	status    int
	errorCode int
}

var responseIns *inCorrect
var once sync.Once

func Ins() *inCorrect {
	once.Do(func() { responseIns = &inCorrect{} })
	return responseIns
}

func (cls *inCorrect) Get() map[string]interface{} {
	ret := map[string]interface{}{
		"msg":      cls.msg,
		"content":  cls.content,
		"status":   cls.status,
		"erroCode": cls.errorCode,
	}
	return ret
}

func (cls *inCorrect) Set(msg string, content interface{}, status int, errorCode int) *inCorrect {
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

func (cls *inCorrect) UnAuthorization(msg string) (int, interface{}) {
	if msg == "" {
		msg = "未授权"
	}
	return 406, cls.Set(msg, nil, 406, 1).Get()
}

func (cls *inCorrect) ErrUnLogin() (int, map[string]interface{}) {
	return 401, cls.Set("未登录", nil, 401, 2).Get()
}

func (cls *inCorrect) Forbidden(msg string) (int, interface{}) {
	if msg == "" {
		msg = "禁止操作"
	}

	return 403, cls.Set(msg, nil, 403, 3).Get()
}

func (cls *inCorrect) Empty(msg string) (int, interface{}) {
	if msg == "" {
		msg = "数不存在"
	}

	return 404, cls.Set(msg, nil, 404, 4).Get()
}

func (cls *inCorrect) Validate(msg string, content interface{}) (int, map[string]interface{}) {
	if msg == "" {
		msg = "表单验证错误"
	}

	return 421, cls.Set(msg, content, 421, 5).Get()
}

func (cls *inCorrect) Accident(msg string, err interface{}) (int, map[string]interface{}) {
	if msg == "" {
		msg = "意外错误"
	}
	return 500, cls.Set(msg, err, 500, 6).Get()
}
