package main

import (
	"changeDetector/dto"
	"changeDetector/route"
	"changeDetector/task"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func getAuthUser() dto.User {
	authUser := dto.User{}
	db, _ := gorm.Open(sqlite.Open(task.DBPath), &gorm.Config{})
	db.First(&authUser)
	return authUser
}

func main() {
	//test()
	task.UpdateTask()
	router := gin.Default()
	route.PublicRoute(router)
	route.AuthorizedRoute(router, getAuthUser())
	// 启动服务
	log.Printf("Server is running on port %s", Config.Server.Port)
	err := router.Run(Config.Server.Host + ":" + Config.Server.Port)
	if err != nil {
		panic(err)
	}
}
