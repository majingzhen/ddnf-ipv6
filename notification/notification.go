package notification

import (
	"fmt"
	"net/smtp"

	"ddns-ipv6/config"
)

// SendNotification 发送邮件通知
func SendNotification(emailCfg config.Email, subject, body string) error {
	auth := smtp.PlainAuth("", emailCfg.Username, emailCfg.Password, emailCfg.SMTPServer)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", emailCfg.Username, emailCfg.Recipient, subject, body)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", emailCfg.SMTPServer, emailCfg.SMTPPort),
		auth,
		emailCfg.Username,
		[]string{emailCfg.Recipient},
		[]byte(msg),
	)

	return err
}
