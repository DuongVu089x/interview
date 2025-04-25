package config

import (
	"os"
	"strconv"
)

// CustomerServiceConfig holds customer service configuration
type CustomerServiceConfig struct {
	Host string
	Port string
}

// Config holds all configuration for the application
type Config struct {
	MongoDB         MongoDBConfig
	Kafka           KafkaConfig
	Redis           RedisConfig
	Server          ServerConfig
	CustomerService CustomerServiceConfig
	GRPC            GRPCConfig
	Observability   ObservabilityConfig `mapstructure:"observability"`
}

// ObservabilityConfig holds all observability-related configuration
type ObservabilityConfig struct {
	Jaeger     JaegerConfig `mapstructure:"jaeger"`
	Prometheus PrometheusConfig
	Logging    LoggingConfig
}

// JaegerConfig holds Jaeger tracing configuration
type JaegerConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	AgentHost   string `mapstructure:"agent_host"`
	AgentPort   string `mapstructure:"agent_port"`
	ServiceName string `mapstructure:"service_name"`
}

// PrometheusConfig holds Prometheus metrics configuration
type PrometheusConfig struct {
	Port     string
	Endpoint string
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level      string
	OutputPath string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
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

// GRPCConfig holds gRPC configuration
type GRPCConfig struct {
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
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", ""),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8081"),
		},
		CustomerService: CustomerServiceConfig{
			Host: getEnv("CUSTOMER_SERVICE_HOST", "localhost"),
			Port: getEnv("CUSTOMER_SERVICE_PORT", "8080"),
		},
		GRPC: GRPCConfig{
			Port: getEnv("GRPC_PORT", "50051"),
		},
		Observability: ObservabilityConfig{
			Jaeger: JaegerConfig{
				ServiceName: getEnv("JAEGER_SERVICE_NAME", "order-service"),
				AgentHost:   getEnv("JAEGER_AGENT_HOST", "localhost"),
				AgentPort:   getEnv("JAEGER_AGENT_PORT", "6831"),
				Enabled:     getEnvAsBool("JAEGER_ENABLED", true),
			},
			Prometheus: PrometheusConfig{
				Port:     getEnv("PROMETHEUS_PORT", "2112"),
				Endpoint: getEnv("PROMETHEUS_ENDPOINT", "/metrics"),
			},
			Logging: LoggingConfig{
				Level:      getEnv("LOG_LEVEL", "info"),
				OutputPath: getEnv("LOG_OUTPUT_PATH", "logs/order-service.log"),
				MaxSize:    getEnvAsInt("LOG_MAX_SIZE", 100),  // 100MB
				MaxBackups: getEnvAsInt("LOG_MAX_BACKUPS", 3), // 3 backups
				MaxAge:     getEnvAsInt("LOG_MAX_AGE", 28),    // 28 days
				Compress:   getEnvAsBool("LOG_COMPRESS", true),
			},
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
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// Helper function to get an environment variable as a boolean with a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}
