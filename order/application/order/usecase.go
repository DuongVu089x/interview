package order

import (
	"fmt"
	"time"

	appcontext "github.com/DuongVu089x/interview/order/component/appctx"
	"github.com/DuongVu089x/interview/order/domain"

	domainidgen "github.com/DuongVu089x/interview/order/domain/id_gen"
	domainorder "github.com/DuongVu089x/interview/order/domain/order"
	pb "github.com/DuongVu089x/interview/order/proto/customer"
	"google.golang.org/grpc/codes"
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

func (uc *UseCase) CreateOrder(ctx appcontext.AppContext, req CreateOrderRequest) (*OrderResponse, error) {
	// Check if customer exists
	customerResp, err := uc.customerClient.GetCustomer(ctx.GetDefaultContext(), &pb.GetCustomerRequest{
		UserId: req.UserID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("failed to check customer existence: %w", err)
	}
	if !customerResp.Exists {
		return nil, fmt.Errorf("customer not found")
	}

	customer := customerResp.Customer

	fmt.Println(customer)

	// Convert DTO to domain entity
	order := uc.mapper.ToEntity(req)

	// Calculate total using domain service
	order.TotalAmount = uc.orderService.CalculateTotal(order.Items)

	if err := uc.orderService.ValidateOrder(order); err != nil {
		return nil, err
	}

	id, _, err := uc.idgenService.GenerateID("ORDER")
	if err != nil {
		return nil, err
	}
	order.OrderID = id
	order.OrderCode = fmt.Sprintf("O%08d", id)

	// Save order
	if err := uc.orderService.CreateOrder(order); err != nil {
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
		fmt.Println("Error sending order to Kafka:", err)
	}

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
