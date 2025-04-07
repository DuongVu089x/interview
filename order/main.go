package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DuongVu089x/interview/order/api/middleware"
	"github.com/DuongVu089x/interview/order/api/rest/order"
	"github.com/DuongVu089x/interview/order/component/appctx"
	"github.com/DuongVu089x/interview/order/config"
	"github.com/DuongVu089x/interview/order/infrastructure/kafka"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Function to initialize main database connection
func initMainDB(cfg *config.Config) (*mongo.Client, error) {
	db, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	err = db.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	return db, nil
}

// Function to initialize read database connection
func initReadDB(cfg *config.Config) (*mongo.Client, error) {
	db, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoDB.ReadURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	err = db.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	return db, nil
}

// Function to initialize Kafka producer
func initKafkaProducer(cfg *config.Config) (*kafka.Producer, error) {
	// Convert config format
	kafkaTopics := make([]kafka.TopicConfig, 0, len(cfg.Kafka.Topics))

	for _, topic := range cfg.Kafka.Topics {
		kafkaTopics = append(kafkaTopics, kafka.TopicConfig{
			Name:              topic.Name,
			NumPartitions:     topic.NumPartitions,
			ReplicationFactor: topic.ReplicationFactor,
		})
	}

	producerConfig := kafka.ProducerConfig{
		BootstrapServers: cfg.Kafka.BootstrapServers,
		SecurityProtocol: cfg.Kafka.SecurityProtocol,
		DefaultTopic:     cfg.Kafka.DefaultTopic,
		Topics:           kafkaTopics,
	}

	// Create topics before producing
	err := kafka.CreateTopics(producerConfig)
	if err != nil {
		log.Printf("Warning: Topic creation failed: %s", err)
		return nil, fmt.Errorf("failed to create producer: %s", err)
	}

	producer, err := kafka.NewProducer(producerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %s", err)
	}

	return producer, nil
}

// Function to initialize Kafka consumer
func initKafkaConsumer(cfg *config.Config) (*kafka.RetryableConsumer, error) {
	// Convert config format
	kafkaTopics := make([]kafka.TopicConfig, 0, len(cfg.Kafka.Topics))
	for _, topic := range cfg.Kafka.Topics {
		kafkaTopics = append(kafkaTopics, kafka.TopicConfig{
			Name:              topic.Name,
			NumPartitions:     topic.NumPartitions,
			ReplicationFactor: topic.ReplicationFactor,
		})
	}

	producerConfig := kafka.ProducerConfig{
		BootstrapServers: cfg.Kafka.BootstrapServers,
		SecurityProtocol: cfg.Kafka.SecurityProtocol,
		DefaultTopic:     cfg.Kafka.DefaultTopic,
		Topics:           kafkaTopics,
	}

	// Initialize Kafka consumer with config
	consumerConfig := kafka.ConsumerConfig{
		BootstrapServers: cfg.Kafka.BootstrapServers,
		SecurityProtocol: cfg.Kafka.SecurityProtocol,
		GroupID:          "my-consumer-group",
		AutoOffsetReset:  "earliest",
	}

	retryConfig := kafka.RetryConfig{
		RetryTopicSuffix:    "-retry",
		DLQTopicSuffix:      "-dlq",
		MaxRetryAttempts:    3,
		RetryBackoffInitial: 1 * time.Second,
		RetryBackoffMax:     1 * time.Minute,
		RetryBackoffFactor:  2,
	}

	// Create consumer without storing the return in a variable
	consumer, err := kafka.NewRetryableConsumer(consumerConfig, producerConfig, retryConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %s", err)
	}

	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start consuming messages in a goroutine
	go func() {
		defer func() {
			consumer.Close()
			cancel()
		}()

		if err := consumer.Start(ctx); err != nil && err != context.Canceled {
			log.Fatalf("Consumer error: %s", err)
		}
	}()

	return consumer, nil
}

func initRedis(cfg *config.Config) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %v", err)
	}

	return redisClient, nil
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize infrastructure
	mainDB, err := initMainDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize main database: %v", err)
		return
	}

	readDB, err := initReadDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize read database: %v", err)
		return
	}

	kafkaProducer, err := initKafkaProducer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
		return
	}
	defer kafkaProducer.Close()

	// kafkaConsumer, err := initKafkaConsumer(cfg)
	// if err != nil {
	// 	log.Fatalf("Failed to initialize Kafka consumer: %v", err)
	// 	return
	// }

	redisClient, err := initRedis(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
		return
	}

	// Initialize AppContext
	appctx := appctx.NewAppContext(mainDB, readDB, kafkaProducer, nil, redisClient)

	// Initialize Echo
	e := echo.New()

	// Add middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.ConfigureCORS())
	e.Use(middleware.RequestLogger())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	// Initialize handlers
	orderHandler := order.NewHandler(appctx)

	// Register routes
	order.RegisterRoutes(e, orderHandler)

	// Print all registered routes for debugging
	middleware.PrintRegisteredRoutes(e)

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)
	if err := e.Start(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	e.Start(serverAddr)
}
