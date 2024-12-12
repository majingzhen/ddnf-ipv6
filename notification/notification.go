package notification

import (
	"fmt"
	"net/smtp"

	"ddns-ipv6/config"
)

// SendNotification 发送邮件通知
func SendNotification(cfg config.EmailConfig, subject, body string) error {
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.SMTPServer)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", cfg.Username, cfg.ToEmail, subject, body)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", cfg.SMTPServer, cfg.SMTPPort),
		auth,
		cfg.Username,
		[]string{cfg.ToEmail},
		[]byte(msg),
	)

	return err
}
