package email

import "context"

// Email represents an email message
type Email struct {
	To          []string
	CC          []string
	BCC         []string
	Subject     string
	Body        string
	IsHTML      bool
	Attachments []Attachment
}

// Attachment represents an email attachment
type Attachment struct {
	Name    string
	Content []byte
}

// EmailService is the interface for email service
type EmailService interface {
	// Send sends an email
	Send(ctx context.Context, email *Email) error

	// SendWithTemplate sends an email using a template
	SendWithTemplate(ctx context.Context, to []string, templateID string, data interface{}) error
}

// Config represents email configuration
type Config struct {
	SMTPHost    string
	SMTPPort    int
	Username    string
	Password    string
	FromAddress string
	FromName    string
	UseTLS      bool
}
