package Controllers

import "time"

type User struct {
	Id int64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username string `gorm:"type:VARCHAR(64);NOT NULL;unique_index:users__username__uidx;"`
	Password string `gorm:"type:VARCHAR(128);NOT NULL;"`
	Nickname string `gorm:"type:VARCHAR(64);NOT NULL;DEFAULT '';index:users__nickname__idx;"`
}

type UserController struct {
	user User
}

func (cls *UserController) SetUser(username, password, nickname string) *UserController {
	user := User{
		Username: username,
		Password: password,
		Nickname: nickname,
	}
	cls.user = user
	return cls
}

func (cls *UserController) GetUser() User {
	return cls.user
}

func (cls *UserController) Login() bool {
	return cls.user.Username == "zhangsan" && cls.user.Password == "zhangsan123"
}
