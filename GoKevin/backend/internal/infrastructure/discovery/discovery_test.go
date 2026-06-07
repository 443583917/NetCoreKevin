package discovery

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsulDiscovery_Register(t *testing.T) {
	discovery := NewConsulDiscovery(Config{
		Provider: "consul",
		Address:  "127.0.0.1",
		Port:     8500,
	})

	service := &Service{
		ID:      "test-1",
		Name:    "test-service",
		Address: "127.0.0.1",
		Port:    8080,
	}

	err := discovery.Register(context.Background(), service)
	assert.NoError(t, err)
}

func TestConsulDiscovery_Deregister(t *testing.T) {
	discovery := NewConsulDiscovery(Config{
		Provider: "consul",
		Address:  "127.0.0.1",
		Port:     8500,
	})

	err := discovery.Deregister(context.Background(), "test-1")
	assert.NoError(t, err)
}

func TestConsulDiscovery_Discover(t *testing.T) {
	discovery := NewConsulDiscovery(Config{
		Provider: "consul",
		Address:  "127.0.0.1",
		Port:     8500,
	})

	services, err := discovery.Discover(context.Background(), "test-service")
	assert.NoError(t, err)
	assert.Len(t, services, 1)
	assert.Equal(t, "test-service", services[0].Name)
}

func TestConsulDiscovery_HealthCheck(t *testing.T) {
	discovery := NewConsulDiscovery(Config{
		Provider: "consul",
		Address:  "127.0.0.1",
		Port:     8500,
	})

	healthy, err := discovery.HealthCheck(context.Background(), "test-1")
	assert.NoError(t, err)
	assert.True(t, healthy)
}

func TestConsulDiscovery_String(t *testing.T) {
	discovery := NewConsulDiscovery(Config{
		Provider: "consul",
		Address:  "127.0.0.1",
		Port:     8500,
	})

	str := discovery.String()
	assert.Contains(t, str, "127.0.0.1")
	assert.Contains(t, str, "8500")
}

func TestServiceDiscovery_Interface(t *testing.T) {
	var sd ServiceDiscovery = NewConsulDiscovery(Config{
		Provider: "consul",
		Address:  "127.0.0.1",
		Port:     8500,
	})

	assert.NotNil(t, sd)
}
