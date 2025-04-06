package appctx

import (
	"context"

	"github.com/DuongVu089x/interview/customer/infrastructure/kafka"
	"github.com/DuongVu089x/interview/customer/websocket"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppContext interface {
	GetMainDBConnection() *mongo.Client
	GetReadMainDBConnection() *mongo.Client

	GetKafkaProducer() *kafka.Producer
	GetKafkaConsumer() *kafka.RetryableConsumer

	GetRedisClient() *redis.Client

	GetDefaultContext() context.Context

	WithContext(c context.Context) AppContext

	GetWebSocketServer() *websocket.WSServer
}

type appCtx struct {
	mainDB     *mongo.Client
	readMainDB *mongo.Client

	ctx context.Context

	kafkaProducer *kafka.Producer
	kafkaConsumer *kafka.RetryableConsumer

	redisClient *redis.Client

	wsServer *websocket.WSServer
}

func NewAppContext(mainDB *mongo.Client, readMainDB *mongo.Client, kafkaProducer *kafka.Producer, kafkaConsumer *kafka.RetryableConsumer, redisClient *redis.Client, wsServer *websocket.WSServer) *appCtx {
	return &appCtx{
		mainDB:     mainDB,
		readMainDB: readMainDB,

		kafkaProducer: kafkaProducer,
		kafkaConsumer: kafkaConsumer,

		redisClient: redisClient,

		wsServer: wsServer,
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

func (ctx *appCtx) GetRedisClient() *redis.Client {
	return ctx.redisClient
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

func (ctx *appCtx) GetWebSocketServer() *websocket.WSServer {
	return ctx.wsServer
}
