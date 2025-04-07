package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	MongoDB MongoDBConfig
	Kafka   KafkaConfig
	Redis   RedisConfig
	Server  ServerConfig
}

// MongoDBConfig holds MongoDB configuration
type MongoDBConfig struct {
	URI      string
	ReadURI  string
	Database string
}

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	BootstrapServers string
	SecurityProtocol string
	DefaultTopic     string
	Topics           []TopicConfig
}

// TopicConfig holds configuration for a Kafka topic
type TopicConfig struct {
	Name              string
	NumPartitions     int
	ReplicationFactor int
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGODB_URI", ""),
			ReadURI:  getEnv("MONGODB_READ_URI", ""),
			Database: getEnv("MONGODB_DATABASE", "orders"),
		},
		Kafka: KafkaConfig{
			BootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS", ""),
			SecurityProtocol: getEnv("KAFKA_SECURITY_PROTOCOL", ""),
			DefaultTopic:     getEnv("KAFKA_DEFAULT_TOPIC", ""),
			Topics: []TopicConfig{
				{
					Name:              "orders-topic",
					NumPartitions:     3,
					ReplicationFactor: 3,
				},
			},
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", ""),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8081"),
		},
	}
}

// Helper function to get an environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Helper function to get an environment variable as an integer with a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
