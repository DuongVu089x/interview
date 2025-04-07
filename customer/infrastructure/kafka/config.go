package kafka

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// TopicConfig holds configuration for a Kafka topic
type TopicConfig struct {
	Name              string
	NumPartitions     int
	ReplicationFactor int
}

// Config holds Kafka configuration
type ProducerConfig struct {
	BootstrapServers string
	SecurityProtocol string
	DefaultTopic     string
	Topics           []TopicConfig // Add topic configurations
}

// RetryConfig holds configuration for retry behavior
type RetryConfig struct {
	// Retry topics configuration
	RetryTopicSuffix    string        // Suffix for retry topic (e.g., "-retry")
	DLQTopicSuffix      string        // Suffix for dead letter queue topic (e.g., "-dlq")
	MaxRetryAttempts    int           // Maximum number of retry attempts before sending to DLQ
	RetryBackoffInitial time.Duration // Initial backoff duration
	RetryBackoffMax     time.Duration // Maximum backoff duration
	RetryBackoffFactor  float64       // Backoff multiplier between retries

	// Optional function to customize backoff time based on retry count
	BackoffStrategy func(attempt int) time.Duration

	// Clone producer config
	CloneProducerConfig ProducerConfig
}

// ConsumerConfig holds Kafka consumer configuration
type ConsumerConfig struct {
	BootstrapServers            string
	SecurityProtocol            string
	GroupID                     string
	AutoOffsetReset             string
	SessionTimeoutMs            int
	HeartbeatIntervalMs         int
	MaxPollIntervalMs           int
	PartitionAssignmentStrategy string
	EnableAutoCommit            bool
	AutoCommitIntervalMs        int
}

// NewConfigMap converts our config to Kafka's ConfigMap
func (c *ProducerConfig) NewConfigMap() *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers": c.BootstrapServers,
		"security.protocol": c.SecurityProtocol,
		"acks":              "all",
	}
}

// NewConsumerConfigMap converts our config to Kafka's ConfigMap
func (c *ConsumerConfig) NewConfigMap() *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers":             c.BootstrapServers,
		"security.protocol":             c.SecurityProtocol,
		"group.id":                      c.GroupID,
		"auto.offset.reset":             c.AutoOffsetReset,
		"enable.auto.commit":            c.EnableAutoCommit,
		"session.timeout.ms":            c.SessionTimeoutMs,
		"heartbeat.interval.ms":         c.HeartbeatIntervalMs,
		"max.poll.interval.ms":          c.MaxPollIntervalMs,
		"partition.assignment.strategy": c.PartitionAssignmentStrategy,
		"auto.commit.interval.ms":       c.AutoCommitIntervalMs,
	}
}
