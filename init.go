package main

import (
	"changeDetector/dto"
	"changeDetector/task"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"os"
)

func loadJson(db *gorm.DB, table string, path string) {
	// 打开json文件
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	switch table {
	case "item":
		var items []dto.WebItem
		json.Unmarshal(byteValue, &items)
		db.Create(&items)
	case "rule":
		var rules []dto.WebRule
		json.Unmarshal(byteValue, &rules)
		db.Create(&rules)
	}
}

type ServerType struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type ConfigType struct {
	Server ServerType   `yaml:"server"`
	Email  dto.Profile  `yaml:"email"`
	Task   dto.CronTask `yaml:"task"`
	User   dto.User     `yaml:"user"`
}

var Config = ConfigType{
	Server: ServerType{Host: "0.0.0.0", Port: "8080"},
}

func loadConfig(path string) *ConfigType {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		fmt.Println(err.Error())
	}
	return &Config
}

func init() {
	loadConfig("config.yaml")
	_, err := os.Stat(task.DBPath)
	if err == nil {
		return
	}
	err = os.Mkdir("data", os.ModePerm)
	db, err := gorm.Open(sqlite.Open("data/change-detector.db"), &gorm.Config{})
	db.AutoMigrate(
		&dto.WebItem{},
		&dto.WebRule{},
		&dto.Profile{},
		&dto.Notice{},
		&dto.User{},
		&dto.CronTask{},
	)
	loadJson(db, "item", "data/items.json")
	loadJson(db, "rule", "data/rules.json")
	db.Create(&dto.Notice{
		App:   "System",
		Title: "Welcome to ChangeDetector",
		Body:  "Welcome to ChangeDetector, please configure your profile and start checking!",
		Link:  "/setting.html",
	})

	if Config.Email.SmtpUser != "" {
		db.Create(&Config.Email)
	}

	if Config.Task.CronTime != "" {
		db.Create(&Config.Task)
	} else {
		cronTask := dto.CronTask{
			TaskName: "Default Task",
			CronTime: "25 17 * * *",
		}
		db.Create(&cronTask)
		log.Printf("Create Default task : %s\n", cronTask.CronTime)
	}

	if Config.User.UserName != "" {
		db.Create(&Config.User)
	} else {
		user := dto.User{
			UserName: "Admin",
			PassWord: "Admin",
		}
		db.Create(&user)
		log.Printf("Create Default User: %s Password: %s\n", user.UserName, user.PassWord)
	}
}
