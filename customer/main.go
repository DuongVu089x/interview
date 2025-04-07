package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DuongVu089x/interview/customer/api/rest/notification"
	"github.com/DuongVu089x/interview/customer/application/consumer"
	"github.com/DuongVu089x/interview/customer/component/appctx"
	"github.com/DuongVu089x/interview/customer/component/server"
	"github.com/DuongVu089x/interview/customer/config"
	"github.com/DuongVu089x/interview/customer/infrastructure/kafka"
	"github.com/DuongVu089x/interview/customer/middleware"
	"github.com/DuongVu089x/interview/customer/websocket"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	userconnrepository "github.com/DuongVu089x/interview/customer/repository/user_connection"
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

// Function to setup websocket
func setupWebsocket(cfg *config.Config, mainDB *mongo.Client) *websocket.WSServer {
	wsServer := websocket.NewWSServer("customer")
	wsServer.Expose(cfg.Server.WSPort)

	wsRoute := wsServer.NewRoute("/notifications")

	userConnRepository := userconnrepository.NewMongoRepository(mainDB)
	handler := websocket.NewWebSocketHandler(userConnRepository, wsServer)
	wsRoute.OnConnect = handler.OnWSConnect
	wsRoute.OnMessage = handler.OnWSMessage
	wsRoute.OnClose = handler.OnWSClose

	go func() {
		wsServer.Start()
	}()

	return wsServer
}

// Function to initialize Kafka consumer
func initKafkaConsumer(cfg *config.Config) (*kafka.RetryableConsumer, error) {

	// Initialize Kafka producer with config
	producerConfig := kafka.ProducerConfig{
		BootstrapServers: cfg.Kafka.BootstrapServers,
		SecurityProtocol: cfg.Kafka.SecurityProtocol,
		DefaultTopic:     cfg.Kafka.DefaultTopic,
	}

	// Initialize Kafka consumer with config
	consumerConfig := kafka.ConsumerConfig{
		BootstrapServers:            cfg.Kafka.BootstrapServers,
		SecurityProtocol:            cfg.Kafka.SecurityProtocol,
		GroupID:                     "my-consumer-group",
		AutoOffsetReset:             "earliest",
		SessionTimeoutMs:            45000,
		HeartbeatIntervalMs:         14000, // Should be lower than session timeout
		MaxPollIntervalMs:           300000,
		PartitionAssignmentStrategy: "roundrobin",
		EnableAutoCommit:            true,
		AutoCommitIntervalMs:        5000,
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

	return consumer, nil
}

func main() {
	cfg := config.LoadConfig()

	// Initialize infrastructure
	mainDB, err := initMainDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize main database: %v", err)
	}

	readDB, err := initReadDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize read database: %v", err)
	}

	kafkaConsumer, err := initKafkaConsumer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka consumer: %v", err)
	}

	wsServer := setupWebsocket(cfg, mainDB)

	appCtx := appctx.NewAppContext(mainDB, readDB, nil, kafkaConsumer, nil, wsServer)

	// Initialize consumer service
	notificationConsumer := consumer.NewNotificationConsumer(appCtx)
	consumerService := consumer.NewConsumerService(appCtx, notificationConsumer)

	// Start consumer service in a goroutine with context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := consumerService.SetupConsumers(ctx); err != nil {
			log.Fatalf("Failed to start consumer service: %v", err)
		}
	}()

	// Initialize HTTP server
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.ConfigureCORS())
	e.Use(middleware.RequestLogger())

	// Register routes
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})
	notificationHandler := notification.NewHandler(appCtx)
	notification.RegisterRoutes(e, notificationHandler)

	// Print routes for debugging
	middleware.PrintRegisteredRoutes(e)

	// Start server with graceful shutdown
	srv := server.NewServer(e, cfg.Server.Port, appCtx)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
