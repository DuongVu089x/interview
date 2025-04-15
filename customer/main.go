package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	customergrpchandler "github.com/DuongVu089x/interview/customer/api/grpc/customer"
	"github.com/DuongVu089x/interview/customer/api/rest/customer"
	"github.com/DuongVu089x/interview/customer/api/rest/notification"
	"github.com/DuongVu089x/interview/customer/application/consumer"
	"github.com/DuongVu089x/interview/customer/component/appctx"
	"github.com/DuongVu089x/interview/customer/component/observability"
	"github.com/DuongVu089x/interview/customer/config"
	"github.com/DuongVu089x/interview/customer/infrastructure/kafka"
	"github.com/DuongVu089x/interview/customer/middleware"
	"github.com/DuongVu089x/interview/customer/websocket"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	pb "github.com/DuongVu089x/interview/customer/proto/customer"
	userconnrepository "github.com/DuongVu089x/interview/customer/repository/user_connection"
	"google.golang.org/grpc"
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
func setupWebsocket(cfg *config.Config, mainDB, readDB *mongo.Client) *websocket.WSServer {
	wsServer := websocket.NewWSServer("customer")
	wsServer.Expose(cfg.Server.WSPort)

	wsRoute := wsServer.NewRoute("/notifications")

	userConnRepository := userconnrepository.NewMongoRepository(mainDB, readDB)
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

// Function to initialize gRPC server
func initGrpcServer(appCtx appctx.AppContext, cfg *config.Config) error {
	// Create a new gRPC server
	server := grpc.NewServer()

	// Register the customer service
	customerHandler := customergrpchandler.NewGrpcHandler(appCtx)
	pb.RegisterCustomerServiceServer(server, customerHandler)

	// Start gRPC server in a goroutine
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPC.Port))
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		log.Printf("Starting gRPC server on port %s", cfg.GRPC.Port)
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	return nil
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize observability
	tp, err := observability.InitTracer(&cfg.Observability)
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}
	defer tp.Shutdown(context.Background())

	logger, err := observability.NewLogger(&cfg.Observability)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	if err := observability.InitMetrics(&cfg.Observability); err != nil {
		logger.Error("Failed to initialize metrics", zap.Error(err))
	}

	// Replace standard logger with zap logger
	zap.ReplaceGlobals(logger)
	logger = logger.Named("customer-service")

	// Get tracer
	tracer := observability.GetTracer("customer-service")

	// Initialize infrastructure with observability
	mainDB, err := initMainDB(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize main database", zap.Error(err))
	}

	readDB, err := initReadDB(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize read database", zap.Error(err))
	}

	kafkaConsumer, err := initKafkaConsumer(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize Kafka consumer", zap.Error(err))
	}

	wsServer := setupWebsocket(cfg, mainDB, readDB)

	// Create AppContext with observability components
	appCtx := appctx.NewAppContext(
		mainDB,
		readDB,
		nil,
		kafkaConsumer,
		nil,
		wsServer,
		logger,
		tracer,
	)

	// Initialize consumer service with observability
	notificationConsumer := consumer.NewNotificationConsumer(appCtx)
	consumerService := consumer.NewConsumerService(appCtx, notificationConsumer)

	// Start consumer service in a goroutine with context and observability
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := consumerService.SetupConsumers(ctx); err != nil {
			appCtx.GetLogger().Fatal("Failed to start consumer service", zap.Error(err))
		}
	}()

	// Initialize HTTP server with observability middleware
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

	customerHandler := customer.NewRestHandler(appCtx)
	customer.RegisterRoutes(e, customerHandler)

	// Initialize gRPC server with observability
	if err := initGrpcServer(appCtx, cfg); err != nil {
		appCtx.GetLogger().Fatal("Failed to initialize gRPC server", zap.Error(err))
	}

	// Start HTTP server
	appCtx.GetLogger().Info("Starting HTTP server", zap.String("port", cfg.Server.Port))
	if err := e.Start(":" + cfg.Server.Port); err != nil {
		appCtx.GetLogger().Fatal("Failed to start server", zap.Error(err))
	}
}
