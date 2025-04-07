package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/DuongVu089x/interview/customer/component/appctx"
)

// ConsumerService manages the lifecycle of message consumers
type ConsumerService struct {
	appCtx    appctx.AppContext
	consumers []Consumer
}

// Consumer interface defines what a consumer should implement
type Consumer interface {
	Setup() error
}

// NewConsumerService creates a new consumer service
func NewConsumerService(appCtx appctx.AppContext, consumers ...Consumer) *ConsumerService {
	return &ConsumerService{
		appCtx:    appCtx,
		consumers: consumers,
	}
}

// Start initializes and starts all registered consumers
func (s *ConsumerService) SetupConsumers(ctx context.Context) error {
	// Setup all consumers
	for _, consumer := range s.consumers {
		if err := consumer.Setup(); err != nil {
			return fmt.Errorf("failed to setup consumer: %w", err)
		}
	}

	log.Println("All consumers have been started successfully")

	kafkaConsumer := s.appCtx.GetKafkaConsumer()
	if kafkaConsumer != nil {
		err := kafkaConsumer.Start(context.Background())
		if err != nil {
			return fmt.Errorf("error starting kafka consumer: %v", err)
		}
	}
	return nil
}
