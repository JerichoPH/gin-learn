package Controllers

import (
	"time"

	"github.com/gin-gonic/gin"
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

type UserRegister struct {
	Username      string `form:"username" binding:"reuiqred"`
	Password      string `from:"password" binding:"required"`
	PasswordCheck string `form:"password_check" binding:"required eqfield=Password"`
	Nickname      string `form:"nickname" binding:"required"`
}

type UserController struct {
	user         User
	userRegister UserRegister
}

func (cls *UserController) BindFormReigster(ctx *gin.Context) *UserController {
	var userRegister UserRegister
	if err := ctx.ShouldBind(&userRegister); err != nil {
		panic(err)
	}

	cls.userRegister = userRegister
	return cls
}

func (cls *UserController) Register(db *gorm.DB) *UserController {
	cls.user.Username = cls.userRegister.Username
	cls.user.Password = cls.userRegister.Password
	cls.user.Nickname = cls.userRegister.Nickname

	db.Create(&cls.user)

	return cls
}

func (cls *UserController) BindForm(ctx *gin.Context) *UserController {
	var user User
	if err := ctx.ShouldBind(&user); err != nil {
		panic(err)
	}
	cls.user = user
	return cls
}

func (cls *UserController) Store(db *gorm.DB) *UserController {
	db.Create(&cls.user)
	return cls
}

func (cls *UserController) GetUser() User {
	return cls.user
}
