package sms

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAliyunSMS_Send(t *testing.T) {
	sms := NewAliyunSMS(Config{
		Provider: "aliyun",
		SignName: "TestSign",
	})

	err := sms.Send(context.Background(), "13800138000", "Test message")
	assert.NoError(t, err)
}

func TestTencentSMS_Send(t *testing.T) {
	sms := NewTencentSMS(Config{
		Provider: "tencent",
		SignName: "TestSign",
	})

	err := sms.Send(context.Background(), "13800138000", "Test message")
	assert.NoError(t, err)
}
