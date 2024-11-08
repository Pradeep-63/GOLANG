// config/smtp.go
package config

import (
	"log"

	"gopkg.in/gomail.v2"
)

func SendEmail(recipients []string, subject, body string) error {
	log.Printf("Sending email to: %v", recipients)
	m := gomail.NewMessage()
	m.SetHeader("From", AppConfig.SmtpFromEmail)
	m.SetHeader("To", recipients...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(AppConfig.SmtpServer, 587, AppConfig.SmtpEmail, AppConfig.SmtpPassword)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Println("Email sent successfully")
	return nil
}
