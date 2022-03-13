package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/mysql"

	"gin-learn/Controllers"
	"gin-learn/Errors"
	"gin-learn/Middlewares"
	"gin-learn/Tools"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	dsn := "root:root@tcp(127.0.0.1:3306)/gin_learn?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{CreateBatchSize: 1000})

	errAutoMigrate := db.AutoMigrate(&Controllers.User{})
	if errAutoMigrate != nil {
		fmt.Println("自动迁移错误：", errAutoMigrate)
		return nil
	}
	// Set table options
	errSetTable := db.Set("gorm:table_options", "ENGINE=Distributed(cluster, default, hits)").AutoMigrate(&Controllers.User{})
	if errSetTable != nil {
		fmt.Println("增加表错误：", errSetTable)
		return nil
	}
	return db
}

func main() {
	db := initDB()
	router := gin.Default()

	router.Use(Errors.RecoverHandler) // 异常处理

	// 权鉴
	v1Authorization := router.Group("/v1/authorization")
	{
		// 注册
		v1Authorization.POST("/register", func(ctx *gin.Context) {
			userController := &Controllers.UserController{CTX: *ctx, DB: *db}
			userController.BindFormRegister().Register()
			ctx.JSON(Tools.ResponseINS().Ok("注册成功", nil))
		})

		// 登录
		v1Authorization.POST("/login", func(ctx *gin.Context) {
			userController := &Controllers.UserController{CTX: *ctx, DB: *db}
			token := userController.BindFormLogin().Login()
			ctx.JSON(Tools.ResponseINS().Ok("登陆成功", gin.H{"token": token}))
		})
	}

	// 用户
	v1User := router.Group("/v1/user", Middlewares.JwtCheck())
	{
		// 列表
		v1User.GET("/", func(ctx *gin.Context) {
			userController := &Controllers.UserController{CTX: *ctx, DB: *db}
			users := userController.GetUsers()
			ctx.JSON(Tools.ResponseINS().Ok("", gin.H{"users": users}))
		})

		// 根据id获取用户详情
		v1User.GET("/:id", func(ctx *gin.Context) {
			userController := &Controllers.UserController{CTX: *ctx, DB: *db}
			user := userController.FindById().GetUser()
			ctx.JSON(Tools.ResponseINS().Ok("", gin.H{"user": user}))
		})

	}

	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	serverErr := server.ListenAndServe()
	if serverErr != nil {
		log.Println("服务器启动错误：", serverErr)
	}
}
