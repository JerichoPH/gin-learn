package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/mysql"

	"gin-learn/Controllers"
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

	router.Use(Tools.Recover) // 异常处理

	v1 := router.Group("/v1")
	{

		v1.POST("/authorization/register", func(ctx *gin.Context) {
			userController := &Controllers.UserController{CTX: *ctx, DB: *db}
			userController.BindFormRegister().Register()
			ctx.JSON(Tools.ResponseINS().Ok("注册成功", nil))
		})

		v1.POST("/authorization/login", func(ctx *gin.Context) {
			userController := &Controllers.UserController{CTX: *ctx, DB: *db}
			user := userController.BindFormLogin().Login()
			ctx.JSON(Tools.ResponseINS().Ok("登陆成功", user))
		})

		v1.GET("/user", func(ctx *gin.Context) {
			userController := &Controllers.UserController{CTX: *ctx, DB: *db}
			users := userController.GetUsers()
			ctx.JSON(Tools.ResponseINS().Ok("", gin.H{"users": users}))
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
