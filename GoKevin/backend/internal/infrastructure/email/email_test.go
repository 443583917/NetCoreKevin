package email

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSMTPEmail_New(t *testing.T) {
	email := NewSMTPEmail(Config{
		SMTPHost:    "smtp.example.com",
		SMTPPort:    587,
		Username:    "test@example.com",
		Password:    "password",
		FromAddress: "test@example.com",
		FromName:    "Test",
	})

	assert.NotNil(t, email)
}

func TestEmail_Struct(t *testing.T) {
	email := &Email{
		To:      []string{"user@example.com"},
		Subject: "Test",
		Body:    "Hello",
	}

	assert.Len(t, email.To, 1)
	assert.Equal(t, "Test", email.Subject)
}
