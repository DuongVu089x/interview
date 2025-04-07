package port

import (
	"context"
)

// ConsumerService defines the interface for message consumption services
type ConsumerService interface {
	// Start starts all registered consumers
	Start(ctx context.Context) error
	// Stop gracefully stops all consumers
	Stop() error
}

// ConsumerRegistry defines the interface for registering message handlers
type ConsumerRegistry interface {
	// RegisterHandlers registers all message handlers for this consumer service
	RegisterHandlers() error
}
