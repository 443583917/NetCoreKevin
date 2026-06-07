package sms

import (
	"context"
	"fmt"
	"log"
)

// AliyunSMS implements SMS interface for Aliyun
type AliyunSMS struct {
	config Config
}

// NewAliyunSMS creates a new Aliyun SMS service
func NewAliyunSMS(config Config) *AliyunSMS {
	return &AliyunSMS{config: config}
}

// Send sends an SMS
func (s *AliyunSMS) Send(ctx context.Context, phone string, content string) error {
	// In production, call Aliyun SMS API
	log.Printf("[AliyunSMS] Send to %s: %s", phone, content)
	return nil
}

// SendWithTemplate sends SMS with template
func (s *AliyunSMS) SendWithTemplate(ctx context.Context, phone string, templateID string, params map[string]string) error {
	log.Printf("[AliyunSMS] Send to %s with template %s: %v", phone, templateID, params)
	return nil
}

// SendBatch sends SMS to multiple phones
func (s *AliyunSMS) SendBatch(ctx context.Context, phones []string, content string) error {
	for _, phone := range phones {
		if err := s.Send(ctx, phone, content); err != nil {
			return fmt.Errorf("send to %s: %w", phone, err)
		}
	}
	return nil
}
