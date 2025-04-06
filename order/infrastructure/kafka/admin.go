package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// CreateTopics creates Kafka topics with the specified configurations
func CreateTopics(config ProducerConfig) error {
	// Skip if no topics are configured
	if len(config.Topics) == 0 {
		return nil
	}

	// Create admin client
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": config.BootstrapServers,
	})
	if err != nil {
		return fmt.Errorf("failed to create admin client: %w", err)
	}
	defer adminClient.Close()

	// Convert our topic configs to Kafka's format
	topics := make([]kafka.TopicSpecification, len(config.Topics))
	for i, topic := range config.Topics {
		topics[i] = kafka.TopicSpecification{
			Topic:             topic.Name,
			NumPartitions:     topic.NumPartitions,
			ReplicationFactor: topic.ReplicationFactor,
		}
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create the topics
	results, err := adminClient.CreateTopics(ctx, topics)
	if err != nil {
		return fmt.Errorf("failed to create topics: %w", err)
	}

	// Check individual topic creation results
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError {
			// Topic might already exist - this is usually not fatal
			fmt.Printf("Warning: Topic '%s': %s\n", result.Topic, result.Error.String())
		} else {
			fmt.Printf("Created topic '%s' with %d partitions\n",
				result.Topic,
				getPartitionsForTopic(result.Topic, config.Topics))
		}
	}

	return nil
}

// Helper to find partition count for a topic
func getPartitionsForTopic(name string, configs []TopicConfig) int {
	for _, config := range configs {
		if config.Name == name {
			return config.NumPartitions
		}
	}
	return 1 // Default
}
