package appctx

import (
	"context"

	"github.com/DuongVu089x/interview/order/infrastructure/kafka"
	pb "github.com/DuongVu089x/interview/order/proto/customer"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type AppContext interface {
	GetMainDBConnection() *mongo.Client
	GetReadMainDBConnection() *mongo.Client

	GetKafkaProducer() *kafka.Producer
	GetKafkaConsumer() *kafka.RetryableConsumer

	// GetRedisClient() *redis.Client

	GetCustomerClient() pb.CustomerServiceClient

	GetDefaultContext() context.Context

	WithContext(c context.Context) AppContext

	// Observability methods
	GetLogger() *zap.Logger
	GetTracer() trace.Tracer
	WithLogger(logger *zap.Logger) AppContext
}

type appCtx struct {
	mainDB     *mongo.Client
	readMainDB *mongo.Client

	ctx context.Context

	kafkaProducer *kafka.Producer
	kafkaConsumer *kafka.RetryableConsumer

	// redisClient *redis.Client

	customerClient pb.CustomerServiceClient

	// Observability components
	logger *zap.Logger
	tracer trace.Tracer
}

func NewAppContext(
	mainDB *mongo.Client,
	readMainDB *mongo.Client,
	kafkaProducer *kafka.Producer,
	kafkaConsumer *kafka.RetryableConsumer,
	// redisClient *redis.Client,
	customerClient pb.CustomerServiceClient,
	logger *zap.Logger,
	tracer trace.Tracer,
) *appCtx {
	return &appCtx{
		mainDB:        mainDB,
		readMainDB:    readMainDB,
		kafkaProducer: kafkaProducer,
		kafkaConsumer: kafkaConsumer,
		// redisClient:    redisClient,
		customerClient: customerClient,
		logger:         logger,
		tracer:         tracer,
	}
}

func (ctx *appCtx) GetMainDBConnection() *mongo.Client {
	return ctx.mainDB
}

func (ctx *appCtx) GetReadMainDBConnection() *mongo.Client {
	return ctx.readMainDB
}

func (ctx *appCtx) GetKafkaProducer() *kafka.Producer {
	return ctx.kafkaProducer
}

func (ctx *appCtx) GetKafkaConsumer() *kafka.RetryableConsumer {
	return ctx.kafkaConsumer
}

// func (ctx *appCtx) GetRedisClient() *redis.Client {
// 	return ctx.redisClient
// }

func (ctx *appCtx) GetCustomerClient() pb.CustomerServiceClient {
	return ctx.customerClient
}

func (ctx *appCtx) GetDefaultContext() context.Context {
	if ctx.ctx == nil {
		ctx.ctx = context.Background()
	}
	return ctx.ctx
}

// WithContext creates a new AppContext with the given context
func (ctx *appCtx) WithContext(c context.Context) AppContext {
	clone := *ctx
	clone.ctx = c
	return &clone
}

// GetLogger returns the logger instance
func (ctx *appCtx) GetLogger() *zap.Logger {
	return ctx.logger
}

// GetTracer returns the tracer instance
func (ctx *appCtx) GetTracer() trace.Tracer {
	return ctx.tracer
}

// WithLogger creates a new AppContext with the given logger
func (ctx *appCtx) WithLogger(logger *zap.Logger) AppContext {
	clone := *ctx
	clone.logger = logger
	return &clone
}
