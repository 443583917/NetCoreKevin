package sms

import (
	"context"
	"fmt"
	"log"
)

// TencentSMS implements SMS interface for Tencent Cloud
type TencentSMS struct {
	config Config
}

// NewTencentSMS creates a new Tencent SMS service
func NewTencentSMS(config Config) *TencentSMS {
	return &TencentSMS{config: config}
}

// Send sends an SMS
func (s *TencentSMS) Send(ctx context.Context, phone string, content string) error {
	log.Printf("[TencentSMS] Send to %s: %s", phone, content)
	return nil
}

// SendWithTemplate sends SMS with template
func (s *TencentSMS) SendWithTemplate(ctx context.Context, phone string, templateID string, params map[string]string) error {
	log.Printf("[TencentSMS] Send to %s with template %s: %v", phone, templateID, params)
	return nil
}

// SendBatch sends SMS to multiple phones
func (s *TencentSMS) SendBatch(ctx context.Context, phones []string, content string) error {
	for _, phone := range phones {
		if err := s.Send(ctx, phone, content); err != nil {
			return fmt.Errorf("send to %s: %w", phone, err)
		}
	}
	return nil
}
