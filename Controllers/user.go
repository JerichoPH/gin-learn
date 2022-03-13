package Controllers

import (
	"gin-learn/Errors"
	"reflect"
	"time"

	"github.com/dgrijalva/jwt-go"
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

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
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
func (cls *UserController) Login() string {
	var user User
	cls.DB.Where(&User{Username: cls.userFormLogin.Username}).First(&user)

	if reflect.DeepEqual(user, User{}) {
		panic(Errors.ThrowEmpty("用户不存在"))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cls.userFormLogin.Password)); err != nil {
		panic(Errors.ThrowUnAuthorization("账号或密码错误"))
	}

	cls.user = user

	token, err := GenerateToken(cls.user.Username, cls.user.Password)
	if err != nil {
		// 生成jwt错误
		panic(err)
	}

	return token
}

var jwtSecret = []byte("gin-learn") // 加密密钥

// 根据用户的用户名和密码产生token
func GenerateToken(username, password string) (string, error) {
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

// 根据传入的token值获取到Claims对象信息，（进而获取其中的用户名和密码）
func ParseToken(token string) (*Claims, error) {

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
