package order

import (
	"context"
	"fmt"
	"time"

	appcontext "github.com/DuongVu089x/interview/order/component/appctx"
	"github.com/DuongVu089x/interview/order/domain"

	domainidgen "github.com/DuongVu089x/interview/order/domain/id_gen"
	domainorder "github.com/DuongVu089x/interview/order/domain/order"
	pb "github.com/DuongVu089x/interview/order/proto/customer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UseCase struct {
	mapper       *Mapper
	orderService domainorder.Service

	// helper repository
	idgenService   domainidgen.Service
	customerClient pb.CustomerServiceClient
}

func NewOrderUseCase(
	orderService domainorder.Service,
	idgenService domainidgen.Service,
	customerClient pb.CustomerServiceClient,
) *UseCase {

	mapper := &Mapper{}

	return &UseCase{
		mapper:         mapper,
		orderService:   orderService,
		idgenService:   idgenService,
		customerClient: customerClient,
	}
}

func (uc *UseCase) CreateOrder(ctx appcontext.AppContext, reqCtx context.Context, req CreateOrderRequest) (*OrderResponse, error) {
	// Start a new span for the CreateOrder operation
	spanCtx, span := ctx.GetTracer().Start(reqCtx, "order.CreateOrder")
	defer span.End()

	// Add relevant attributes to the span
	span.SetAttributes(
		attribute.String("user_id", req.UserID),
		attribute.Int("items_count", len(req.Items)),
	)

	// Check if customer exists
	// Start a new span for the GetCustomer operation
	getCustomerCtx, getCustomerSpan := ctx.GetTracer().Start(spanCtx, "order.GetCustomer")
	defer getCustomerSpan.End()

	// Add customer ID to the span
	getCustomerSpan.SetAttributes(
		attribute.String("user_id", req.UserID),
		attribute.String("operation", "GetCustomer"),
		attribute.String("service", "order-service"),
	)

	// Create metadata with trace context
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(getCustomerCtx, carrier)
	md := metadata.New(map[string]string{})
	for k, v := range carrier {
		md.Set(k, v)
	}
	getCustomerCtx = metadata.NewOutgoingContext(getCustomerCtx, md)

	customerResp, err := uc.customerClient.GetCustomer(getCustomerCtx, &pb.GetCustomerRequest{
		UserId: req.UserID,
	})
	if err != nil {
		getCustomerSpan.RecordError(err)
		getCustomerSpan.SetStatus(otelcodes.Error, err.Error())
		if st, ok := status.FromError(err); ok && st.Code() == grpccodes.NotFound {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("failed to check customer existence: %w", err)
	}
	if !customerResp.Exists {
		err := fmt.Errorf("customer not found")
		getCustomerSpan.RecordError(err)
		getCustomerSpan.SetStatus(otelcodes.Error, "customer not found")
		return nil, err
	}

	// Convert DTO to domain entity
	order := uc.mapper.ToEntity(req)

	// Calculate total using domain service
	order.TotalAmount = uc.orderService.CalculateTotal(order.Items)

	if err := uc.orderService.ValidateOrder(order); err != nil {
		span.RecordError(err)
		span.SetStatus(otelcodes.Error, err.Error())
		return nil, err
	}

	// Generate order ID
	id, _, err := uc.idgenService.GenerateID("ORDER")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelcodes.Error, err.Error())
		return nil, err
	}
	order.OrderID = id
	order.OrderCode = fmt.Sprintf("O%08d", id)

	// Save order
	if err := uc.orderService.CreateOrder(order); err != nil {
		span.RecordError(err)
		span.SetStatus(otelcodes.Error, err.Error())
		return nil, err
	}

	// Send order to Kafka
	err = ctx.GetKafkaProducer().Publish(domain.Message{
		Key:   fmt.Sprintf("ORDER_CREATED_%d", order.OrderID),
		Topic: "orders-topic",
		Value: domain.MessageValue{
			Meta: &domain.MetaData{
				MessageID: fmt.Sprintf("ORDER_CREATED_%d", order.OrderID),
				ServiceID: "order-service",
				Timestamp: time.Now().UnixNano(),
			},
			MessageCode: "ORDER_CREATED",
			Payload: map[string]any{
				"order_id": fmt.Sprintf("%d", order.OrderID),
				"user_id":  order.UserID,
				"amount":   order.TotalAmount,
				"status":   order.Status,
			},
		},
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelcodes.Error, err.Error())
		fmt.Println("Error sending order to Kafka:", err)
	}

	// Add order details to span
	span.SetAttributes(
		attribute.String("order_id", fmt.Sprintf("%d", order.OrderID)),
		attribute.String("order_code", order.OrderCode),
		attribute.Float64("total_amount", order.TotalAmount),
		attribute.String("status", string(order.Status)),
	)

	// Convert domain entity to response DTO
	response := uc.mapper.ToResponse(order)
	return &response, nil
}

func (uc *UseCase) GetOrder(id int64) (*OrderResponse, error) {
	order, err := uc.orderService.GetOrder(id)
	if err != nil {
		return nil, err
	}

	// Convert domain entity to response DTO
	response := uc.mapper.ToResponse(order)
	return &response, nil
}

// GetOrdersByUserID retrieves all orders for a specific user
func (uc *UseCase) GetOrdersByUserID(req GetOrdersByUserIDRequest) (*OrderListResponse, error) {
	// Prepare conditions map for filtering
	conditions := make(map[string]any)
	if req.Status != "" {
		conditions["status"] = domainorder.OrderStatus(req.Status)
	}

	// Get orders from domain service
	orders, err := uc.orderService.GetOrderByCustomerID(req.UserID, conditions)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	// Convert domain entities to response DTOs
	var orderResponses []OrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, uc.mapper.ToResponse(&order))
	}

	// Return the list of orders
	return &OrderListResponse{
		Orders: orderResponses,
		Count:  len(orderResponses),
	}, nil
}
