package port

import (
	"context"

	domain "github.com/DuongVu089x/interview/customer/domain"
)

// MessageConsumer defines the interface for message consumption
type MessageConsumer interface {

	// RegisterHandler registers a handler for a specific topic
	RegisterHandler(topic string, handler func(domain.Message) error) error

	// Subscribe subscribes to all registered topics
	Subscribe() error

	// Start starts consuming messages
	Start(ctx context.Context) error

	// Close closes the consumer
	Close() error
}
