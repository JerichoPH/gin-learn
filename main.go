package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"net/http"
	"time"

	"gin-learn/Controllers"

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
	initDB() // 初始化数据库

	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		a := 3 << 2
		fmt.Println(a)
		//userController := &Controllers.UserController{}
		//userController.SetUser("zhangsan", "zhangsan123", "张三")
		//
		//user := userController.GetUser()
		//user1 := db.Omit("UserName", "Password", "CreatedAt").Create(&user)
		//user2 := db.Select("UserName", "Password", "CreatedAt").Create(&user)
		//
		//fmt.Printf("%v\r\n", *user1)
		//fmt.Printf("%v\r\n", *user2)
	})

	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	serverErr := server.ListenAndServe()
	if serverErr != nil {
		return
	}
}
