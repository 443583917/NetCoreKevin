package email

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

// SMTPEmail implements EmailService using SMTP
type SMTPEmail struct {
	config Config
}

// NewSMTPEmail creates a new SMTP email service
func NewSMTPEmail(config Config) *SMTPEmail {
	return &SMTPEmail{config: config}
}

// Send sends an email
func (s *SMTPEmail) Send(ctx context.Context, email *Email) error {
	addr := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)

	// Build message
	msg := s.buildMessage(email)

	// Authentication
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.SMTPHost)

	// Collect all recipients
	to := append(email.To, email.CC...)
	to = append(to, email.BCC...)

	err := smtp.SendMail(addr, auth, s.config.FromAddress, to, []byte(msg))
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	log.Printf("[Email] Sent to %v: %s", email.To, email.Subject)
	return nil
}

// SendWithTemplate sends email with template
func (s *SMTPEmail) SendWithTemplate(ctx context.Context, to []string, templateID string, data interface{}) error {
	// In production, load template and render
	email := &Email{
		To:      to,
		Subject: "Template Email",
		Body:    fmt.Sprintf("Template: %s, Data: %v", templateID, data),
		IsHTML:  false,
	}

	return s.Send(ctx, email)
}

// buildMessage builds the email message
func (s *SMTPEmail) buildMessage(email *Email) string {
	var msg strings.Builder

	msg.WriteString(fmt.Sprintf("From: %s <%s>\r\n", s.config.FromName, s.config.FromAddress))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(email.To, ",")))

	if len(email.CC) > 0 {
		msg.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(email.CC, ",")))
	}

	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))

	if email.IsHTML {
		msg.WriteString("MIME-Version: 1.0\r\n")
		msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	} else {
		msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	}

	msg.WriteString("\r\n")
	msg.WriteString(email.Body)

	return msg.String()
}
