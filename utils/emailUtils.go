package utils

import (
	"fmt"
	"net/smtp"
	"testbot/conf"
)

var smtpServer string
var senderEmail string
var auth smtp.Auth
var mailType string

func init() {
	smtpServer = conf.Config.Email.SmtpServer
	senderEmail = conf.Config.Email.SenderEmail
	senderPassword := conf.Config.Email.SenderPassword
	auth = smtp.PlainAuth("",
		senderEmail,
		senderPassword,
		smtpServer)
	mailType = "Content-Type: text/html; charset=UTF-8"

}
func SendEmail(toEmail []string, subject string, body string) {
	s := fmt.Sprintf("To:%s\r\n"+
		"From:%s <%s>\r\n"+
		"Subject:%s\r\n"+
		"%s\r\n\r\n%s",
		toEmail[0], "miniBot", senderEmail, subject, mailType, body)
	msg := []byte(s)
	Info(s)

	err := smtp.SendMail(
		smtpServer+":587",
		auth, senderEmail,
		toEmail,
		msg)
	if err != nil {
		Error(err)
	}

}
