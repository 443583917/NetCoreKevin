package sms

import "context"

// SMS is the interface for SMS service
type SMS interface {
	// Send sends an SMS message
	Send(ctx context.Context, phone string, content string) error

	// SendWithTemplate sends an SMS using a template
	SendWithTemplate(ctx context.Context, phone string, templateID string, params map[string]string) error

	// SendBatch sends SMS to multiple phones
	SendBatch(ctx context.Context, phones []string, content string) error
}

// Config represents SMS configuration
type Config struct {
	Provider     string // aliyun, tencent
	AccessKey    string
	AccessSecret string
	SignName     string
	Region       string
}
