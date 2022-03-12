package Controllers

import (
	"gin-learn/Errors"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id        int64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string `form:"username" uri:"username" binding:"required" gorm:"type:VARCHAR(64);NOT NULL;unique_index:users__username__uidx;"`
	Password  string `form:"password" binding:"required" gorm:"type:VARCHAR(128);NOT NULL;"`
	Nickname  string `form:"nickname" uri:"nickname" binding:"required" gorm:"type:VARCHAR(64);NOT NULL;DEFAULT '';index:users__nickname__idx;"`
}

type UserFormRegister struct {
	Username      string `form:"username" binding:"required"`
	Password      string `form:"password" binding:"required"`
	PasswordCheck string `form:"password_check" binding:"required"`
	Nickname      string `form:"nickname" binding:"required"`
}

type UserFormLogin struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type UserController struct {
	CTX              gin.Context
	DB               gorm.DB
	user             User
	users            []User
	userFormRegister UserFormRegister
	userFormLogin    UserFormLogin
}

// 获取用户数据
func (cls *UserController) GetUser() User {
	return cls.user
}

// 获取列表
func (cls *UserController) GetUsers() []User {
	cls.DB.Find(&cls.users)

	return cls.users
}

// 绑定注册表单
func (cls *UserController) BindFormRegister() *UserController {
	var userRegister UserFormRegister
	if err := cls.CTX.ShouldBind(&userRegister); err != nil {
		panic(err)
	}

	if userRegister.Password != userRegister.PasswordCheck {
		panic(Errors.ThrowForbidden("两次密码输入不一致"))
	}

	cls.userFormRegister = userRegister
	return cls
}

// 注册
func (cls *UserController) Register() *UserController {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(cls.userFormRegister.Password), 14)

	cls.user = User{
		Username: cls.userFormRegister.Username,
		Password: string(bytes),
		Nickname: cls.userFormRegister.Nickname,
	}

	cls.DB.Create(&cls.user)

	return cls
}

// 绑定登录表单
func (cls *UserController) BindFormLogin() *UserController {
	var userLogin UserFormLogin
	if err := cls.CTX.ShouldBind(&userLogin); err != nil {
		panic(err)
	}

	cls.userFormLogin = userLogin

	return cls
}

// 登录
func (cls *UserController) Login() *UserController {
	var user User
	cls.DB.Where(&User{Username: cls.userFormLogin.Username}).First(&user)

	if reflect.DeepEqual(user, User{}) {
		panic(Errors.ThrowEmpty("用户不存在"))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cls.userFormLogin.Password)); err != nil {
		panic(Errors.ThrowUnAuthorization("账号或密码错误"))
	}

	cls.user = user

	return cls
}
