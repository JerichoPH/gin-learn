package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/mysql"

	"gin-learn/controllers"
	"gin-learn/errors"
	"gin-learn/middlewares"
	"gin-learn/tools"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	dsn := "root:root@tcp(127.0.0.1:3306)/gin_learn?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{CreateBatchSize: 1000})

	errAutoMigrate := db.AutoMigrate(&controllers.Account{})
	if errAutoMigrate != nil {
		fmt.Println("自动迁移错误：", errAutoMigrate)
		return nil
	}
	// Set table options
	errSetTable := db.Set("gorm:table_options", "ENGINE=Distributed(cluster, default, hits)").AutoMigrate(&controllers.Account{})
	if errSetTable != nil {
		fmt.Println("增加表错误：", errSetTable)
		return nil
	}
	return db
}

func main() {
	db := initDB()
	router := gin.Default()

	router.Use(errors.RecoverHandler) // 异常处理

	// 注册
	router.POST("/v1/authorization/register", func(ctx *gin.Context) {
		userController := &controllers.AccountController{CTX: *ctx, DB: *db}
		userController.BindFormRegister().Register()
		ctx.JSON(tools.CorrectIns().Ok("注册成功", nil))
	})

	// 登录
	router.POST("/v1/authorization/login", func(ctx *gin.Context) {
		userController := &controllers.AccountController{CTX: *ctx, DB: *db}
		token := userController.BindFormLogin().Login()
		ctx.JSON(tools.CorrectIns().Ok("登陆成功", gin.H{"token": token}))
	})

	v1 := router.Group("/v1", middlewares.JwtCheck(db))
	{
		// 用户列表
		v1.GET("/account", func(ctx *gin.Context) {
			accountController := &controllers.AccountController{CTX: *ctx, DB: *db}
			accounts := accountController.FindMoreByQuery().Accounts
			ctx.JSON(tools.CorrectIns().Ok("", gin.H{"accounts": accounts}))
		})

		// 根据id获取用户详情
		v1.GET("/account/:id", func(ctx *gin.Context) {
			accountController := &controllers.AccountController{CTX: *ctx, DB: *db}
			if accountController.FindById().IsEmpty() {
				panic(errors.ThrowEmpty("用户不存在"))
			}
			user := accountController.Account
			ctx.JSON(tools.CorrectIns().Ok("", gin.H{"accounts": user}))
		})

		// 状态列表
		v1.GET("/status", func(ctx *gin.Context) {
			statusController := &controllers.StatusController{CTX: *ctx, DB: *db}
			statuses := statusController.FindMoreByQuery().Statuses
			ctx.JSON(tools.CorrectIns().Ok("", gin.H{"statuses": statuses}))
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
