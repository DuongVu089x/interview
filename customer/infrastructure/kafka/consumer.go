package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	port "github.com/DuongVu089x/interview/customer/application/port"
	domain "github.com/DuongVu089x/interview/customer/domain"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Consumer implements the MessageConsumer interface
type Consumer struct {
	consumer *kafka.Consumer
	handlers map[string]func(domain.Message) error
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(config ConsumerConfig) (*Consumer, error) {
	c, err := kafka.NewConsumer(config.NewConfigMap())
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	return &Consumer{
		consumer: c,
		handlers: make(map[string]func(domain.Message) error),
	}, nil
}

func (c *Consumer) RegisterHandler(topic string, handler func(domain.Message) error) error {
	if _, exists := c.handlers[topic]; exists {
		return fmt.Errorf("handler already registered for topic: %s", topic)
	}

	c.handlers[topic] = handler

	// Subscribe to this topic
	topics := c.getSubscribedTopics()
	return c.consumer.SubscribeTopics(topics, nil)
}

func (c *Consumer) getSubscribedTopics() []string {
	topics := make([]string, 0, len(c.handlers))
	for topic := range c.handlers {
		topics = append(topics, topic)
	}
	return topics
}

func (c *Consumer) Start(ctx context.Context) error {
	if len(c.handlers) == 0 {
		return fmt.Errorf("no handlers registered, cannot start consumer")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := c.consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Timeout is not an error, just continue
				if err.(kafka.Error).Code() == kafka.ErrTimedOut {
					continue
				}
				return err
			}

			// Convert Kafka message to domain message
			value := domain.MessageValue{}
			err = json.Unmarshal(msg.Value, &value)
			if err != nil {
				return fmt.Errorf("failed to unmarshal message value: %w", err)
			}
			domainMsg := domain.Message{
				Key:       string(msg.Key),
				Value:     value,
				Topic:     *msg.TopicPartition.Topic,
				Partition: int(msg.TopicPartition.Partition),
				Offset:    int64(msg.TopicPartition.Offset),
			}

			// Find and execute the appropriate handler
			if handler, ok := c.handlers[domainMsg.Topic]; ok {
				if err := handler(domainMsg); err != nil {
					// Log error but continue processing
					fmt.Printf("Error processing message from topic %s: %v\n",
						domainMsg.Topic, err)
				}
			} else {
				fmt.Printf("No handler registered for topic: %s\n", domainMsg.Topic)
			}
		}
	}
}

// Close shuts down the consumer
func (c *Consumer) Close() error {
	return c.consumer.Close()
}

// Ensure Consumer implements MessageConsumer interface
var _ port.MessageConsumer = (*Consumer)(nil)
