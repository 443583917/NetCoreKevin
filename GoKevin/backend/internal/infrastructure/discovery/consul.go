package discovery

import (
	"context"
	"fmt"
	"log"
)

// ConsulDiscovery implements ServiceDiscovery using Consul
type ConsulDiscovery struct {
	config Config
}

// NewConsulDiscovery creates a new Consul discovery
func NewConsulDiscovery(config Config) *ConsulDiscovery {
	return &ConsulDiscovery{config: config}
}

// Register registers a service
func (d *ConsulDiscovery) Register(ctx context.Context, service *Service) error {
	// In production, call Consul API
	log.Printf("[Consul] Register service: %s (%s:%d)", service.Name, service.Address, service.Port)
	return nil
}

// Deregister deregisters a service
func (d *ConsulDiscovery) Deregister(ctx context.Context, serviceID string) error {
	log.Printf("[Consul] Deregister service: %s", serviceID)
	return nil
}

// Discover discovers service instances
func (d *ConsulDiscovery) Discover(ctx context.Context, serviceName string) ([]*Service, error) {
	// Mock: return a single service instance
	return []*Service{
		{
			ID:      serviceName + "-1",
			Name:    serviceName,
			Address: "127.0.0.1",
			Port:    8080,
		},
	}, nil
}

// HealthCheck checks service health
func (d *ConsulDiscovery) HealthCheck(ctx context.Context, serviceID string) (bool, error) {
	// In production, check Consul health endpoint
	return true, nil
}

// String returns a string representation
func (d *ConsulDiscovery) String() string {
	return fmt.Sprintf("ConsulDiscovery{address: %s:%d}", d.config.Address, d.config.Port)
}
