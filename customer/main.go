package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DuongVu089x/interview/customer/component/appctx"
	"github.com/DuongVu089x/interview/customer/config"
	"github.com/DuongVu089x/interview/customer/domain"
	"github.com/DuongVu089x/interview/customer/infrastructure/kafka"
	"github.com/DuongVu089x/interview/customer/websocket"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	notificationusecase "github.com/DuongVu089x/interview/customer/application/notification"
	notificationdomain "github.com/DuongVu089x/interview/customer/domain/notification"
	notificationrepository "github.com/DuongVu089x/interview/customer/repository/notification"
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

	return consumer, nil
}

func registerConsumerKafka(appCtx appctx.AppContext) error {
	mainDB := appCtx.GetMainDBConnection()
	consumer := appCtx.GetKafkaConsumer()
	wsServer := appCtx.GetWebSocketServer()

	ctx, cancel := context.WithCancel(context.Background())

	// Register different handlers for different topics
	err := consumer.RegisterHandler("orders-topic", func(msg domain.Message) error {

		fmt.Printf("Processing order: key: %s, topic: %s, partition: %d, offset: %d\n", msg.Key, msg.Topic, msg.Partition, msg.Offset)
		fmt.Printf("Processing order: payload: %v\n", msg.Value.Payload)

		payload := msg.Value.Payload
		payloadMap := payload.(map[string]any)

		notificationRepository := notificationrepository.NewMongoRepository(mainDB)
		userConnRepository := userconnrepository.NewMongoRepository(mainDB)

		notificationHandler := websocket.NewWebSocketHandler(userConnRepository, wsServer)
		notificationDispatcher := websocket.NewNotificationDispatcher(wsServer, "/notifications", notificationHandler)
		notificationUsecase := notificationusecase.NewUseCase(notificationRepository, notificationDispatcher)
		notificationUsecase.CreateNotification(&notificationdomain.Notification{
			Topic:       "order-created",
			Title:       "Order Created",
			Description: "Order created successfully",
			Link:        fmt.Sprintf("localhost:8081/order/%f", payloadMap["order_id"].(float64)),
			UserId:      payloadMap["user_id"].(string),
		})

		return nil
	})
	if err != nil {
		log.Fatalf("Failed to register handler: %s", err)
	}

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

	return nil
}

func main() {
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

	kafkaConsumer, err := initKafkaConsumer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka consumer: %v", err)
		return
	}

	wsServer := setupWebsocket(cfg, mainDB)

	appCtx := appctx.NewAppContext(mainDB, readDB, nil, kafkaConsumer, nil, wsServer)

	// register consumer kafka
	registerConsumerKafka(appCtx)

	e := echo.New()

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)
	if err := e.Start(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	e.Start(serverAddr)
}
