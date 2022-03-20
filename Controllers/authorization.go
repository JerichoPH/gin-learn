package controllers

import (
	"gin-learn/errors"
	"gin-learn/models"
	"gin-learn/tools"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"reflect"
	"time"
)

var jwtSecret = []byte("gin-learn") // 加密密钥

// AuthorizationController 验签控制器
type AuthorizationController struct {
	CTX *gin.Context
	DB  *gorm.DB
}

// Claims Jwt 表单
type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// RegisterForm 注册表单
type RegisterForm struct {
	Username      string `form:"username" binding:"required"`
	Password      string `form:"password" binding:"required"`
	PasswordCheck string `form:"password_check" binding:"required"`
	Nickname      string `form:"nickname" binding:"required"`
}

// LoginForm 登录表单
type LoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// PostRegister 注册
func (cls *AuthorizationController) PostRegister() {
	// 表单验证
	var registerForm RegisterForm
	if err := cls.CTX.ShouldBind(&registerForm); err != nil {
		panic(err)
	}
	if registerForm.Password != registerForm.PasswordCheck {
		panic(errors.ThrowForbidden("两次密码输入不一致"))
	}

	// 检查重复项
	accountModel := &models.AccountModel{CTX: cls.CTX, DB: cls.DB}
	repeatAccount := accountModel.FindOneByUsername(registerForm.Username)
	if !reflect.DeepEqual(repeatAccount, models.Account{}) {
		panic(errors.ThrowForbidden("用户名被占用"))
	}

	// 密码加密
	bytes, _ := bcrypt.GenerateFromPassword([]byte(registerForm.Password), 14)

	// 保存新用户
	account := models.Account{
		Username: registerForm.Username,
		Password: string(bytes),
		Nickname: registerForm.Nickname,
	}
	cls.DB.Create(&account)

	cls.CTX.JSON(tools.CorrectIns().Ok("注册成功", nil))
}

// PostLogin 登录
func (cls *AuthorizationController) PostLogin() {
	// 表单验证
	var loginForm LoginForm
	if err := cls.CTX.ShouldBind(&loginForm); err != nil {
		panic(err)
	}

	// 检查用户是否存在
	accountModel := &models.AccountModel{CTX: cls.CTX, DB: cls.DB}
	account := accountModel.FindOneByUsername(loginForm.Username)
	if reflect.DeepEqual(account, models.Account{}) {
		panic(errors.ThrowEmpty("用户不存在"))
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(loginForm.Password)); err != nil {
		panic(errors.ThrowUnAuthorization("账号或密码错误"))
	}

	// 生成Jwt
	token, err := cls.GenerateJwt(account.Username, account.Password)
	if err != nil {
		// 生成jwt错误
		panic(err)
	}

	cls.CTX.Set("__currentAccount", account) // 保存用户数据到本次请求

	cls.CTX.JSON(tools.CorrectIns().Ok("登陆成功", gin.H{"token": token}))
}

// GenerateJwt 生成Jwt
func (cls *AuthorizationController) GenerateJwt(username, password string) (string, error) {
	// 设置token有效时间
	nowTime := time.Now()
	expireTime := nowTime.Add(168 * time.Hour)

	claims := Claims{
		Username: username,
		Password: password,
		StandardClaims: jwt.StandardClaims{
			// 过期时间
			ExpiresAt: expireTime.Unix(),
			// 指定token发行人
			Issuer: "gin-learn",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//该方法内部生成签名字符串，再用于获取完整、已签名的token
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// ParseJwt 根据传入的token值获取到Claims对象信息，（进而获取其中的用户名和密码）
func (cls *AuthorizationController) ParseJwt(token string) (*Claims, error) {

	//用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		// 要传入指针，项目中结构体都是用指针传递，节省空间。
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
