package port

import domain "github.com/DuongVu089x/interview/order/domain"

// MessageProducer defines the interface for sending messages
type MessageProducer interface {
	Publish(message domain.Message) error
	Close() error
}
