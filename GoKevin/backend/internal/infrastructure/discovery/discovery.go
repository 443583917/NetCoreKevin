package discovery

import "context"

// Service represents a service instance
type Service struct {
	ID      string
	Name    string
	Address string
	Port    int
	Tags    []string
	Meta    map[string]string
}

// ServiceDiscovery is the interface for service discovery
type ServiceDiscovery interface {
	// Register registers a service
	Register(ctx context.Context, service *Service) error

	// Deregister deregisters a service
	Deregister(ctx context.Context, serviceID string) error

	// Discover discovers service instances
	Discover(ctx context.Context, serviceName string) ([]*Service, error)

	// HealthCheck checks service health
	HealthCheck(ctx context.Context, serviceID string) (bool, error)
}

// Config represents service discovery configuration
type Config struct {
	Provider string // consul
	Address  string
	Port     int
}
