package consumer

import (
	"fmt"

	notificationusecase "github.com/DuongVu089x/interview/customer/application/notification"
	"github.com/DuongVu089x/interview/customer/component/appctx"
	"github.com/DuongVu089x/interview/customer/domain"
	notificationrepository "github.com/DuongVu089x/interview/customer/repository/notification"
	userconnrepository "github.com/DuongVu089x/interview/customer/repository/user_connection"
	"github.com/DuongVu089x/interview/customer/websocket"
)

// NotificationConsumer handles order notifications
type NotificationConsumer struct {
	appCtx appctx.AppContext
}

// NewNotificationConsumer creates a new notification consumer
func NewNotificationConsumer(appCtx appctx.AppContext) *NotificationConsumer {
	return &NotificationConsumer{
		appCtx: appCtx,
	}
}

// Setup implements the Consumer interface
func (c *NotificationConsumer) Setup() error {
	consumer := c.appCtx.GetKafkaConsumer()
	if consumer == nil {
		return fmt.Errorf("kafka consumer is not initialized")
	}

	mainDB := c.appCtx.GetMainDBConnection()
	wsServer := c.appCtx.GetWebSocketServer()

	// Register order notification handler
	err := consumer.RegisterHandler("orders-topic", func(msg domain.Message) error {
		fmt.Printf("Processing order: key: %s, topic: %s, partition: %d, offset: %d\n", msg.Key, msg.Topic, msg.Partition, msg.Offset)
		fmt.Printf("Processing order: payload: %v\n", msg.Value.Payload)

		payload := msg.Value.Payload
		payloadMap := payload.(map[string]any)

		notificationRepository := notificationrepository.NewMongoRepository(mainDB)
		userConnRepository := userconnrepository.NewMongoRepository(mainDB)

		notificationHandler := websocket.NewWebSocketHandler(userConnRepository, wsServer)
		notificationDispatcher := websocket.NewNotificationDispatcher(wsServer, "/notifications", notificationHandler)
		notificationUseCase := notificationusecase.NewWriteUseCase(notificationRepository, notificationDispatcher)

		return notificationUseCase.CreateNotification(&notificationusecase.CreateNotificationRequest{
			Topic:       "order-created",
			Title:       "Order Created",
			Description: "Order created successfully",
			Link:        fmt.Sprintf("localhost:8081/order/%s", payloadMap["order_id"].(string)),
			UserID:      payloadMap["user_id"].(string),
		})
	})

	if err != nil {
		return fmt.Errorf("failed to register order notification handler: %w", err)
	}

	// Subscribe to all topics after registering handlers
	if err := consumer.Subscribe(); err != nil {
		return fmt.Errorf("failed to subscribe to topics: %w", err)
	}

	return nil
}

// Close implements the Consumer interface
func (c *NotificationConsumer) Close() error {
	consumer := c.appCtx.GetKafkaConsumer()
	if consumer == nil {
		return fmt.Errorf("kafka consumer is not initialized")
	}
	return consumer.Close()
}
