package kafka

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	port "github.com/DuongVu089x/interview/customer/application/port"
	domain "github.com/DuongVu089x/interview/customer/domain"
)

const (
	// Header keys
	RetryCountHeader    = "retry-count"
	OriginalTopicHeader = "original-topic"
	OriginalKeyHeader   = "original-key"
	ErrorMessageHeader  = "error-message"
)

// RetryableConsumer wraps the base consumer with retry functionality
type RetryableConsumer struct {
	baseConsumer  *Consumer
	retryProducer *Producer
	retryConfig   RetryConfig
	handlers      map[string]func(domain.Message) error
}

// NewRetryableConsumer creates a new consumer with retry functionality
func NewRetryableConsumer(consumerConfig ConsumerConfig, producerConfig ProducerConfig, retryConfig RetryConfig) (*RetryableConsumer, error) {
	consumer, err := NewConsumer(consumerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	producer, err := NewProducer(producerConfig)
	if err != nil {
		consumer.Close()
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	// Set default values if not provided
	if retryConfig.RetryTopicSuffix == "" {
		retryConfig.RetryTopicSuffix = "-retry"
	}

	if retryConfig.DLQTopicSuffix == "" {
		retryConfig.DLQTopicSuffix = "-dlq"
	}

	if retryConfig.MaxRetryAttempts == 0 {
		retryConfig.MaxRetryAttempts = 3
	}

	if retryConfig.RetryBackoffInitial == 0 {
		retryConfig.RetryBackoffInitial = 1 * time.Second
	}

	if retryConfig.RetryBackoffMax == 0 {
		retryConfig.RetryBackoffMax = 1 * time.Minute
	}

	if retryConfig.RetryBackoffFactor == 0 {
		retryConfig.RetryBackoffFactor = 2 // Default exponential backoff factor
	}

	if retryConfig.BackoffStrategy == nil {
		// BackoffStrategy calculates the delay before retrying a failed message
		// attempt: The number of retry attempts so far
		// Returns: How long to wait before the next retry
		// The delay increases exponentially (RetryBackoffFactor) up to RetryBackoffMax
		retryConfig.BackoffStrategy = func(attempt int) time.Duration {
			// Calculate backoff duration by multiplying initial duration by attempt number and backoff factor
			backoff := retryConfig.RetryBackoffInitial * time.Duration(float64(attempt)*retryConfig.RetryBackoffFactor)

			// Cap the backoff at the configured maximum
			if backoff > retryConfig.RetryBackoffMax {
				return retryConfig.RetryBackoffMax
			}
			return backoff
		}
	}

	retryConfig.CloneProducerConfig = ProducerConfig{
		BootstrapServers: producerConfig.BootstrapServers,
		SecurityProtocol: producerConfig.SecurityProtocol,
		// DefaultTopic:     producerConfig.DefaultTopic,
		// Topics:           producerConfig.Topics,
	}

	return &RetryableConsumer{
		baseConsumer:  consumer,
		retryProducer: producer,
		retryConfig:   retryConfig,
		handlers:      make(map[string]func(domain.Message) error),
	}, nil
}

// RegisterHandler registers a handler for a topic with retry functionality
func (rc *RetryableConsumer) RegisterHandler(topic string, handler func(msg domain.Message) error) error {

	// Create retry and DLQ topics if they don't exist
	retryTopic := topic + rc.retryConfig.RetryTopicSuffix
	dlqTopic := topic + rc.retryConfig.DLQTopicSuffix

	// Ensure retry topics are created
	retryTopicConfig := TopicConfig{
		Name:              retryTopic,
		NumPartitions:     3,
		ReplicationFactor: 3,
	}
	dlqTopicConfig := TopicConfig{
		Name:              dlqTopic,
		NumPartitions:     3,
		ReplicationFactor: 3,
	}

	// Create the topics (ignoring already exists errors)
	CreateTopics(ProducerConfig{
		BootstrapServers: rc.retryConfig.CloneProducerConfig.BootstrapServers,
		SecurityProtocol: rc.retryConfig.CloneProducerConfig.SecurityProtocol,
		Topics:           []TopicConfig{retryTopicConfig, dlqTopicConfig},
	})

	// Register the handler with retry functionality
	err := rc.safeRegisterHandler(topic, func(msg domain.Message) error {
		return rc.processWithRetry(msg, handler, topic)
	})
	if err != nil {
		return err
	}

	// Register the handler with retry topic
	err = rc.safeRegisterHandler(retryTopic, func(msg domain.Message) error {
		return rc.processRetry(msg, handler, topic)
	})
	if err != nil {
		return err
	}

	// err = rc.baseConsumer.Subscribe()
	// if err != nil {
	// 	return err
	// }

	return nil
}

// Safe handler registration with retries
func (rc *RetryableConsumer) safeRegisterHandler(topic string, handler func(msg domain.Message) error) error {
	maxAttempts := 5
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := rc.baseConsumer.RegisterHandler(topic, handler)
		if err == nil {
			log.Printf("Successfully registered handler for %s", topic)
			return nil
		}

		lastErr = err
		log.Printf("Failed to register handler for %s (attempt %d): %v",
			topic, attempt, err)

		// Wait before retrying
		wait := time.Duration(math.Pow(2, float64(attempt-1))) * time.Second
		log.Printf("Waiting %v before retry...", wait)
		time.Sleep(wait)
	}

	return fmt.Errorf("failed to register handler for %s after %d attempts: %w",
		topic, maxAttempts, lastErr)
}

// processWithRetry handles processing with retry logic for normal messages
func (rc *RetryableConsumer) processWithRetry(msg domain.Message, handler func(msg domain.Message) error, originTopic string) error {
	// msg.Headers = append(msg.Headers, domain.Header{Key: OriginalTopicHeader, Value: originTopic})

	meta := domain.MetaData{}
	meta.OriginalTopic = originTopic
	meta.OriginalKey = msg.Key

	msg.Value.Meta = &meta

	err := handler(msg)
	if err != nil {
		// First failure, send to retry topic
		return rc.sendToRetryTopic(msg, err, 1, originTopic)
	}
	return nil
}

// processRetry handles processing of messages from retry topics
func (rc *RetryableConsumer) processRetry(msg domain.Message, handler func(msg domain.Message) error, originTopic string) error {
	// Extract retry count
	retryCount := 1
	// for _, header := range msg.Value.Meta {
	// if header.Key == RetryCountHeader {
	if msg.Value.Meta.RetryCount != 0 {
		retryCount = msg.Value.Meta.RetryCount
	}
	// }

	// Try processing again
	err := handler(msg)
	if err != nil {
		// Increment retry count
		retryCount++

		// Check if max retries exceeded
		if retryCount > rc.retryConfig.MaxRetryAttempts {
			return rc.sendToDLQ(msg, err, retryCount, originTopic)
		}

		// Send to retry topic again with increased count
		return rc.sendToRetryTopic(msg, err, retryCount, originTopic)
	}

	return nil
}

// sendToRetryTopic sends a failed message to the retry topic
func (rc *RetryableConsumer) sendToRetryTopic(msg domain.Message, processingErr error, retryCount int, originTopic string) error {
	retryTopic := originTopic + rc.retryConfig.RetryTopicSuffix

	// Add or update headers
	// var headers []kafka.Header
	meta := msg.Value.Meta
	meta.RetryCount = retryCount
	meta.OriginalTopic = originTopic
	meta.OriginalKey = msg.Key
	meta.ErrorMessage = processingErr.Error()

	msg.Value.Meta = meta

	// Create retry message
	retryMsg := domain.Message{
		Key:       msg.Key,
		Value:     msg.Value,
		Topic:     retryTopic,
		Partition: 0, // Let Kafka decide
	}

	// Calculate backoff delay
	delay := rc.retryConfig.BackoffStrategy(retryCount)

	// Schedule retry after backoff
	go func() {
		time.Sleep(delay)
		if err := rc.retryProducer.Publish(retryMsg); err != nil {
			log.Printf("Failed to send message to retry topic %s: %v", retryTopic, err)
		}
	}()

	return nil
}

// sendToDLQ sends a failed message to the DLQ topic
func (rc *RetryableConsumer) sendToDLQ(msg domain.Message, processingErr error, retryCount int, originTopic string) error {
	dlqTopic := originTopic + rc.retryConfig.DLQTopicSuffix

	if msg.Value.Meta.RetryCount != 0 {
		msg.Value.Meta.RetryCount = retryCount
	}

	// Add final error message
	msg.Value.Meta.ErrorMessage = processingErr.Error()

	// Create DLQ message
	dlqMsg := domain.Message{
		Key:       msg.Key,
		Value:     msg.Value,
		Topic:     dlqTopic,
		Partition: 0, // Let Kafka decide
	}

	// Send to DLQ immediately
	return rc.retryProducer.Publish(dlqMsg)
}

// Subscribe subscribes to all registered topics
func (rc *RetryableConsumer) Subscribe() error {
	if rc.baseConsumer == nil {
		return fmt.Errorf("consumer not initialized")
	}
	return rc.baseConsumer.Subscribe()
}

// Start starts consuming messages with retry functionality
func (rc *RetryableConsumer) Start(ctx context.Context) error {
	return rc.baseConsumer.Start(ctx)
}

// Close closes both consumer and producer
func (rc *RetryableConsumer) Close() error {
	err := rc.baseConsumer.Close()
	if err != nil {
		return err
	}
	err = rc.retryProducer.Close()
	if err != nil {
		return err
	}
	return nil
}

// Ensure Consumer implements MessageConsumer interface
var _ port.MessageConsumer = (*RetryableConsumer)(nil)
