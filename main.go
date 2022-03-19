package main

import (
	"fmt"
	"gin-learn/controllers"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/mysql"

	"gin-learn/errors"
	"gin-learn/middlewares"
	"gin-learn/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {

	dsn := "root:root@tcp(127.0.0.1:3306)/gin_learn?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{CreateBatchSize: 1000})

	errAutoMigrate := db.AutoMigrate(
		&models.Account{},
		&models.Status{},
	)
	if errAutoMigrate != nil {
		fmt.Println("自动迁移错误：", errAutoMigrate)
		return nil
	}
	// Set table options
	errSetTable := db.Set("gorm:table_options", "ENGINE=Distributed(cluster, default, hits)").AutoMigrate(&models.Account{})
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
		(&controllers.AuthorizationController{CTX: ctx, DB: db}).PostRegister()
	})

	// 登录
	router.POST("/v1/authorization/login", func(ctx *gin.Context) {
		(&controllers.AuthorizationController{CTX: ctx, DB: db}).PostLogin()
	})

	v1 := router.Group("/v1", middlewares.JwtCheck())
	{
		// 用户列表
		v1.GET("/account", func(ctx *gin.Context) {
			(&controllers.AccountController{CTX: ctx, DB: db}).Index()
		})

		// 根据id获取用户详情
		v1.GET("/account/:id", func(ctx *gin.Context) {
			(&controllers.AccountController{CTX: ctx, DB: db}).Show()
		})

		// 状态列表
		v1.GET("/status", func(ctx *gin.Context) {
			(&controllers.StatusController{CTX: ctx, DB: db}).Index()
		})

		// 状态详情
		v1.GET("/status/:id", func(ctx *gin.Context) {
			(&controllers.StatusController{CTX: ctx, DB: db}).Show()
		})

		// 新建状态
		v1.POST("/status", func(ctx *gin.Context) {
			(&controllers.StatusController{CTX: ctx, DB: db}).Store()
		})

		// 编辑状态
		v1.PUT("/status/:id", func(ctx *gin.Context) {
			(&controllers.StatusController{CTX: ctx, DB: db}).Update()
		})

		// 删除状态
		v1.DELETE("/status/:id", func(ctx *gin.Context) {
			(&controllers.StatusController{CTX: ctx, DB: db}).Destroy()
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
