package task

import (
	"changeDetector/dto"
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/smtp"
	"time"
)

func SendEmail(account, password, server, port string) {
	db, err := gorm.Open(sqlite.Open(DBPath), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return
	}
	var notices []dto.Notice
	db.Where("mark=?", 0).Find(&notices)
	if len(notices) == 0 {
		log.Println("No new notice.")
		return
	}

	//设置头信息头部分
	subject := fmt.Sprintf("Subject: %s\n", "应用更新提醒")
	sender := fmt.Sprintf("From: %s\n", account)
	receiver := fmt.Sprintf("To: %s\n", account)
	header := subject + sender + receiver + "MIME-Version: 1.0\nContent-Type: text/html; charset=UTF-8\n\n"

	mailContent := header + generateHtml(notices)

	addr := fmt.Sprintf("%s:%s", server, port)
	from := account
	to := []string{account}
	auth := LoginAuth(account, password)
	//log.Println(addr, auth, from, to, mailContent)
	err = smtp.SendMail(addr, auth, from, to, []byte(mailContent))
	if err != nil {
		log.Fatal(err)
		return
	}
	//log.Println(mailContent)
	db.Where("mark=?", 0).Updates(dto.Notice{Mark: 1})
	log.Println("Send mail success!")
}

func generateHtml(notices []dto.Notice) string {

	var content string
	for _, notice := range notices {
		content += fmt.Sprintf(`
		<li><div>
			<h4>%s：%s</h4>
			<pre>%s</pre>
			<a href="%s">点击查看</a>
		</div></li>`, notice.App, notice.Title, notice.Body, notice.Link)
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
<title>更新检测</title>
</head>
<body>
	<h3>应用更新：</h3>
	<ul>
		%s
	</ul>
	<h3>检查时间：%s</h3>
</body>
</html>`, content, time.Now().Format("2006-01-02 15:04:05"))
}

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("unknown fromServer")
		}
	}
	return nil, nil
}
