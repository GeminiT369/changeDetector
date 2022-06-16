package route

import (
	"changeDetector/dto"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func itemCRUD(router *gin.RouterGroup, db *gorm.DB) {
	router.POST("/addItem", func(c *gin.Context) {
		item := dto.WebItem{}
		err := c.BindJSON(&item)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid Data"})
			return
		}
		db.Create(&item)
		log.Printf("Add item: %v", item)
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	router.POST("/updItem", func(c *gin.Context) {
		item := dto.WebItem{}
		err := c.BindJSON(&item)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid Data"})
			return
		}
		db.Model(&item).Where("id=?", item.ID).Updates(item)
		log.Printf("Update item: %v", item)
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	router.POST("/delItem", func(c *gin.Context) {
		item := dto.WebItem{}
		err := c.BindJSON(&item)
		if err != nil || item.ID == 0 {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid ID"})
			return
		}
		db.Delete(&item)
		log.Printf("delete item: %d", item.ID)
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
}

func ruleCRUD(router *gin.RouterGroup, db *gorm.DB) {
	router.POST("/addRule", func(c *gin.Context) {
		rule := dto.WebRule{}
		err := c.BindJSON(&rule)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid Data"})
			return
		}
		db.Create(&rule)
		log.Printf("Add rule: %v", &rule)
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	router.POST("/updRule", func(c *gin.Context) {
		rule := dto.WebRule{}
		err := c.BindJSON(&rule)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid Data"})
			return
		}
		db.Model(&rule).Where("id=?", rule.ID).Updates(rule)
		log.Printf("Update rule: %v", &rule)
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	router.POST("/delRule", func(c *gin.Context) {
		rule := dto.WebRule{}
		err := c.BindJSON(&rule)
		if err != nil || rule.ID == 0 {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid ID"})
			return
		}
		var webItem dto.WebItem
		db.Where("type = ?", rule.Name).First(&webItem)
		log.Printf("%v", webItem)
		if webItem.ID != 0 {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": "Some items are using this rule"})
			return
		}
		db.Delete(&rule)
		log.Printf("delete rule: %d", rule.ID)
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
}
