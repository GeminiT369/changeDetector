package route

import (
	"changeDetector/dto"
	"changeDetector/middleware"
	"changeDetector/task"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
)

var noticeCount int64

func AuthorizedRoute(router *gin.Engine, authUser dto.User) {
	auth := router.Group("/")
	auth.Use(middleware.JWTAuthMiddleware(authUser))

	db, err := gorm.Open(sqlite.Open(task.DBPath), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.Model(&dto.Notice{}).Count(&noticeCount)

	auth.GET("/api/web/item", func(c *gin.Context) {
		webList := []dto.WebItem{}
		db.Find(&webList)
		c.IndentedJSON(http.StatusOK, webList)
	})
	auth.GET("/api/web/rule", func(c *gin.Context) {
		ruleList := []dto.WebRule{}
		db.Find(&ruleList)
		c.IndentedJSON(http.StatusOK, ruleList)
	})

	auth.GET("/api/web/profile", func(c *gin.Context) {
		profile := dto.Profile{}
		db.Find(&profile)
		c.IndentedJSON(http.StatusOK, profile)
	})

	auth.POST("/api/web/profile", func(c *gin.Context) {
		profile := dto.Profile{}
		c.BindJSON(&profile)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid Data"})
			return
		}
		if profile.ID == 0 {
			db.Create(&profile)
			log.Printf("Create profile: %s", profile.SmtpUser)
		} else {
			db.Model(&profile).Where("id=?", profile.ID).Updates(profile)
			log.Printf("Update profile: %s", profile.SmtpUser)
		}
		c.IndentedJSON(http.StatusOK, profile)
	})

	auth.GET("/api/web/task", func(c *gin.Context) {
		task := dto.CronTask{}
		db.Find(&task)
		c.IndentedJSON(http.StatusOK, task)
	})

	auth.POST("/api/web/task", func(c *gin.Context) {
		cronTask := dto.CronTask{}
		err := c.BindJSON(&cronTask)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid Data"})
			return
		}
		db.Updates(cronTask)
		log.Printf("Update task: %s", cronTask.TaskName)
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Update Cron Task Success"})
		task.UpdateTask()
	})

	auth.GET("/api/web/user", func(c *gin.Context) {
		user := dto.User{}
		db.Find(&user)
		c.IndentedJSON(http.StatusOK, user)
	})

	auth.POST("/api/web/user", func(c *gin.Context) {
		user := dto.User{}
		db.First(&user)
		passwordOld, _ := c.GetPostForm("password_old")
		if user.PassWord != passwordOld {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Password is not correct"})
			return
		}
		user.UserName, _ = c.GetPostForm("username")
		user.PassWord, _ = c.GetPostForm("password_new")
		db.Updates(user)
		middleware.AuthUser = user
		log.Printf("Update user: %s", user.UserName)
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Update user success"})
	})

	auth.GET("/api/web/notice", func(c *gin.Context) {
		var notices []dto.Notice
		db.Find(&notices)
		c.IndentedJSON(http.StatusOK, notices)
		//db.Delete(&notices)
	})

	auth.GET("/api/web/notice_count", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, noticeCount)
	})

	auth.POST("/api/web/notice", func(c *gin.Context) {
		//action := c.PostForm("action")
		action, status := c.GetPostForm("action")
		if !status {
			log.Printf("Not get action")
			return
		}
		log.Printf("action: %s", action)
		if action != "marked" {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Invalid Action"})
			return
		}
		var notices []dto.Notice
		db.Find(&notices)
		db.Delete(&notices)
		noticeCount = 0
		log.Printf("Deleted all notices")
		c.JSON(http.StatusOK, gin.H{"message": "Deleted all notices"})
	})

	auth.GET("/api/web/get", func(c *gin.Context) {
		url := c.Query("url")
		log.Printf("url: %s", url)
		resp, err := http.Get(url)
		defer resp.Body.Close()
		if err != nil {
			log.Printf("error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "Request Failed"})
			return
		}
		body, _ := ioutil.ReadAll(resp.Body)
		c.String(http.StatusOK, string(body))
	})
	crudRoute(auth, db)
}

func crudRoute(auth *gin.RouterGroup, db *gorm.DB) {
	itemCRUD(auth, db)
	ruleCRUD(auth, db)
}
