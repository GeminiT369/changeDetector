package task

import (
	"changeDetector/dto"
	"crypto/md5"
	"encoding/hex"
	"github.com/PuerkitoBio/goquery"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
)

func UpdateTask() {
	// Query CronTime
	db, err := gorm.Open(sqlite.Open(DBPath), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	var task dto.CronTask
	db.First(&task)

	//	reset cron task
	cronObj.Remove(cronId)
	cronId, _ = cronObj.AddFunc(task.CronTime, doCheck)
	cronObj.Start()
	go doCheck()
}

func doCheck() {
	db, err := gorm.Open(sqlite.Open(DBPath), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	var rules []dto.WebRule
	var items []dto.WebItem
	db.Find(&rules)
	db.Find(&items)
	ruleMap := make(map[string]dto.WebRule)
	for _, rule := range rules {
		ruleMap[rule.Name] = rule
	}
	for _, item := range items {
		checkItem(&item, ruleMap[item.Type], db)
	}
	var profile dto.Profile
	db.First(&profile)
	if profile.ID != 0 {
		SendEmail(profile.SmtpUser, profile.SmtpPswd, profile.SmtpServer, profile.SmtpPort)
	}
}

func checkItem(item *dto.WebItem, rule dto.WebRule, db *gorm.DB) {
	resp, err := http.Get(item.Url)
	//resp.Close = true
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("Status code error: %d %s", resp.StatusCode, resp.Status)
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	md5Text := doc.Find(rule.Md5).Text()
	md5String := md5V(md5Text)
	if item.Md5 == "" {
		item.Md5 = md5String
		db.Select("md5").Updates(item)
		log.Printf("init md5 for: %s", item.Name)
		return
	}
	if item.Md5 == md5String {
		log.Printf("md5 not change: %s", item.Name)
		return
	}
	item.Md5 = md5String
	db.Select("md5").Updates(item)
	log.Printf("md5 change: %s", item.Name)
	url := ""
	if rule.Link != "" {
		url, _ = doc.Find(rule.Link).Attr("href")
	}
	if url == "" {
		url = item.Url
	}
	if url[0] == '/' {
		url = item.Url[:strings.Index(item.Url[8:], "/")+8] + url
	}
	body := doc.Find(rule.Body).Text()
	body = MulRn.ReplaceAllString(body, "\n")
	body = MulSp.ReplaceAllString(body, " ")
	notice := dto.Notice{
		App:   item.Name,
		Title: doc.Find(rule.Title).Text(),
		Body:  body,
		Link:  url,
	}
	db.Create(&notice)
	//fmt.Printf("{\n\tApp: %s\n\tTitle: %s\n\tBody: %s\n\tLink: %s\n}", notice.App, notice.Title, notice.Body, notice.Link)
}

func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
