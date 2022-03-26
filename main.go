package main

import (
	"fmt"
	v1 "gin-learn/routes/v1"
	"gorm.io/driver/mysql"
	"log"
	"net/http"
	"time"

	"gin-learn/errors"
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

func initServer(router *gin.Engine) {
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

func main() {
	db := initDB()
	router := gin.Default()

	router.Use(errors.RecoverHandler)             // 异常处理
	(&v1.RoutesV1{Router: router, DB: db}).Load() // 加载v1路由

	initServer(router) // 启动服务
}
