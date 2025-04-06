package port

import (
	"context"

	domain "github.com/DuongVu089x/interview/customer/domain"
)

// MessageConsumer defines the interface for receiving messages
type MessageConsumer interface {
	RegisterHandler(topic string, handler func(domain.Message) error) error
	Start(ctx context.Context) error
	Close() error
}
