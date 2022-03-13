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

type Account struct {
	Id        int64 `gorm:"primary_key" uri:"id" binding:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string `form:"username" gorm:"type:VARCHAR(64);NOT NULL;unique_index:users__username__uidx;"`
	Password  string `form:"password" gorm:"type:VARCHAR(128);NOT NULL;"`
	Nickname  string `form:"nickname" gorm:"type:VARCHAR(64);NOT NULL;DEFAULT '';index:users__nickname__idx;"`
}

type AccountFormRegister struct {
	Username      string `form:"username" binding:"required"`
	Password      string `form:"password" binding:"required"`
	PasswordCheck string `form:"password_check" binding:"required"`
	Nickname      string `form:"nickname" binding:"required"`
}

type AccountFormLogin struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

type AccountController struct {
	CTX                 gin.Context
	DB                  gorm.DB
	account             Account
	accounts            []Account
	accountFormRegister AccountFormRegister
	accountFormLogin    AccountFormLogin
}

func (cls *AccountController) IsEmpty() bool {
	return reflect.DeepEqual(cls.account, Account{})
}

// 根据id获取用户
func (cls *AccountController) FindById() *AccountController {
	id := cls.CTX.Param("id")

	if id == "" {
		panic(Errors.ThrowForbidden("id 不能为空"))
	}

	cls.DB.Where(map[string]interface{}{"id": id}).Find(&cls.account)
	return cls
}

// 根据用户名获取用户
func (cls *AccountController) FindByUsername(username string) *AccountController {
	if username == "" {
		panic(Errors.ThrowForbidden("用户名不能为空"))
	}

	cls.DB.Where(map[string]interface{}{"username": username}).Find(&cls.account)
	return cls
}

// 获取用户数据
func (cls *AccountController) GetAccount() Account {
	return cls.account
}

// 获取列表
func (cls *AccountController) GetAccounts() []Account {
	cls.DB.Find(&cls.accounts)

	return cls.accounts
}

// 绑定注册表单
func (cls *AccountController) BindFormRegister() *AccountController {
	var accountRegister AccountFormRegister
	if err := cls.CTX.ShouldBind(&accountRegister); err != nil {
		panic(err)
	}

	if accountRegister.Password != accountRegister.PasswordCheck {
		panic(Errors.ThrowForbidden("两次密码输入不一致"))
	}

	cls.accountFormRegister = accountRegister
	return cls
}

// 注册
func (cls *AccountController) Register() *AccountController {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(cls.accountFormRegister.Password), 14)

	cls.account = Account{
		Username: cls.accountFormRegister.Username,
		Password: string(bytes),
		Nickname: cls.accountFormRegister.Nickname,
	}

	cls.DB.Create(&cls.account)

	return cls
}

// 绑定登录表单
func (cls *AccountController) BindFormLogin() *AccountController {
	var accountLogin AccountFormLogin
	if err := cls.CTX.ShouldBind(&accountLogin); err != nil {
		panic(err)
	}

	cls.accountFormLogin = accountLogin

	return cls
}

// 登录
func (cls *AccountController) Login() string {
	var account Account
	cls.DB.Where(&Account{Username: cls.accountFormLogin.Username}).First(&account)

	if reflect.DeepEqual(account, Account{}) {
		panic(Errors.ThrowEmpty("用户不存在"))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(cls.accountFormLogin.Password)); err != nil {
		panic(Errors.ThrowUnAuthorization("账号或密码错误"))
	}

	cls.account = account

	token, err := GenerateToken(cls.account.Username, cls.account.Password)
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
