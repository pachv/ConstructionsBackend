package services

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"path/filepath"
	"strings"
	"time"
)

type MailSendingService struct {
	From string

	SMTPHost string
	SMTPPort int

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
	htmlBody, err := renderHTMLTemplate(templatePath, data)
	if err != nil {
		return fmt.Errorf("render template: %w", err)
	}

	msg := buildHTMLMessage(s.From, to, subject, htmlBody)

	if s.SMTPPort == 465 {
		return s.sendTLS(to, msg)
	}
	return s.sendSTARTTLS(to, msg)
}

func renderHTMLTemplate(path string, data any) (string, error) {
	abs := path
	if !filepath.IsAbs(path) {
		abs = filepath.Clean(path)
	}

	tpl, err := template.ParseFiles(abs)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func buildHTMLMessage(from string, to []string, subject string, htmlBody string) []byte {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("From: %s\r\n", from))
	b.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(to, ", ")))
	b.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Date: " + time.Now().Format(time.RFC1123Z) + "\r\n")
	b.WriteString(`Content-Type: text/html; charset="UTF-8"` + "\r\n")
	b.WriteString("\r\n")
	b.WriteString(htmlBody)
	return []byte(b.String())
}

func (s *MailSendingService) sendTLS(to []string, msg []byte) error {
	addr := fmt.Sprintf("%s:%d", s.SMTPHost, s.SMTPPort)

	tlsConfig := &tls.Config{
		ServerName: s.SMTPHost,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("tls dial: %w", err)
	}
	defer conn.Close()

	c, err := smtp.NewClient(conn, s.SMTPHost)
	if err != nil {
		return fmt.Errorf("smtp client: %w", err)
	}
	defer func() { _ = c.Quit() }()

	auth := smtp.PlainAuth("", s.From, s.AppPassword, s.SMTPHost)
	if err := c.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}

	if err := c.Mail(s.From); err != nil {
		return fmt.Errorf("mail from: %w", err)
	}
	for _, rcpt := range to {
		if err := c.Rcpt(rcpt); err != nil {
			return fmt.Errorf("rcpt %s: %w", rcpt, err)
		}
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("data: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		_ = w.Close()
		return fmt.Errorf("write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("close data: %w", err)
	}

	return nil
}

func (s *MailSendingService) sendSTARTTLS(to []string, msg []byte) error {
	addr := fmt.Sprintf("%s:%d", s.SMTPHost, s.SMTPPort)

	c, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	defer c.Close()

	if ok, _ := c.Extension("STARTTLS"); ok {
		if err := c.StartTLS(&tls.Config{ServerName: s.SMTPHost}); err != nil {
			return fmt.Errorf("starttls: %w", err)
		}
	}

	auth := smtp.PlainAuth("", s.From, s.AppPassword, s.SMTPHost)
	if ok, _ := c.Extension("AUTH"); ok {
		if err := c.Auth(auth); err != nil {
			return fmt.Errorf("auth: %w", err)
		}
	}

	if err := c.Mail(s.From); err != nil {
		return fmt.Errorf("mail from: %w", err)
	}
	for _, rcpt := range to {
		if err := c.Rcpt(rcpt); err != nil {
			return fmt.Errorf("rcpt %s: %w", rcpt, err)
		}
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("data: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		_ = w.Close()
		return fmt.Errorf("write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("close data: %w", err)
	}

	return nil
}
