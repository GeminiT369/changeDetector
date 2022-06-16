package route

import (
	"changeDetector/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PublicRoute(router *gin.Engine) {
	router.Static("/static", "static")
	router.StaticFile("/favicon.ico", "static/img/favicon.ico")
	router.GET("/", func(c *gin.Context) {
		c.File("www/index.html")
	})
	router.GET("/rule.html", func(c *gin.Context) {
		c.File("www/rule.html")
	})
	router.GET("/setting.html", func(c *gin.Context) {
		c.File("www/setting.html")
	})
	router.GET("/notice.html", func(c *gin.Context) {
		c.File("www/notice.html")
	})

	router.GET("/login", func(c *gin.Context) {
		c.File("www/login.html")
	})
	router.POST("/api/web/auth", middleware.AuthHandler)

	router.GET("/api/web/checkAuth", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Authorized"})
	})
}
