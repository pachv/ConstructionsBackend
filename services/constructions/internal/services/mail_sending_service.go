package services

import (
	"bytes"
	"fmt"
	"html/template"

	gomail "gopkg.in/gomail.v2"
)

type MailSendingService struct {
	From        string
	SMTPHost    string
	SMTPPort    int
	AppPassword string
}

func NewMailSendingService(from, smtpHost, appPassword string, smtpPort int) *MailSendingService {
	return &MailSendingService{
		From:        from,
		SMTPHost:    smtpHost,
		SMTPPort:    smtpPort,
		AppPassword: appPassword,
	}
}

func (s *MailSendingService) SendHTMLFromTemplate(to []string, subject, templatePath string, data any) error {

	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", buf.String())

	d := gomail.NewDialer(s.SMTPHost, s.SMTPPort, s.From, s.AppPassword)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("send mail: %w", err)
	}

	return nil
}
