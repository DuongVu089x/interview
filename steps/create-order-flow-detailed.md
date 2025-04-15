# Create Order Flow - Detailed Implementation Guide

## Overview

This document provides a detailed breakdown of the order creation flow, including architectural decisions, implementation details, and rationale for each component.

## Architectural Decisions

### 1. Why Domain Services?

Domain Services encapsulate complex business logic that doesn't naturally fit into entities:

```go
// domain/order/service.go
type OrderService interface {
    ValidateOrder(order *Order) error
    CalculateTotal(items []OrderItem) float64
    ApplyVoucher(order *Order, voucher *Voucher) error
}
```

**Rationale:**

1. Separates business rules from entities
2. Makes business logic reusable across use cases
3. Easier to test and maintain
4. Single responsibility principle

### 2. Why Repositories?

Repositories abstract data persistence and provide domain-oriented collection interfaces:

```go
// domain/order/repository.go
type OrderRepository interface {
    Create(ctx context.Context, order *Order) error
    Update(ctx context.Context, order *Order) error
    FindByID(ctx context.Context, id string) (*Order, error)
    FindByCustomerID(ctx context.Context, customerID string) ([]*Order, error)
}
```

**Rationale:**

1. Decouples domain logic from data access
2. Makes storage implementation replaceable
3. Simplifies testing with mocks
4. Provides domain-specific query methods

### 3. Why Use Cases?

Use Cases orchestrate the flow of data and coordinate between different services:

```go
// application/order/usecase.go
type CreateOrderUseCase struct {
    orderService       domain.OrderService
    orderRepo         domain.OrderRepository
    customerClient    pb.CustomerServiceClient
    voucherService   domain.VoucherService
    outboxRepo       domain.OutboxRepository
    transactionMgr   domain.TransactionManager
}
```

**Rationale:**

1. Implements specific application features
2. Coordinates multiple services
3. Handles data transformation
4. Manages transaction boundaries

## Detailed Implementation Steps

### 1. Domain Layer Implementation

#### 1.1 Order Entity

```go
// domain/order/entity.go
type Order struct {
    ID            primitive.ObjectID `bson:"_id,omitempty"`
    CustomerID    string            `bson:"customer_id"`
    Items         []OrderItem       `bson:"items"`
    TotalAmount   float64           `bson:"total_amount"`
    Status        OrderStatus       `bson:"status"`
    VoucherCode   *string           `bson:"voucher_code,omitempty"`
    Discount      float64           `bson:"discount"`
    CreatedAt     time.Time         `bson:"created_at"`
    UpdatedAt     time.Time         `bson:"updated_at"`
}

func (o *Order) Validate() error {
    if o.CustomerID == "" {
        return ErrCustomerIDRequired
    }
    if len(o.Items) == 0 {
        return ErrEmptyOrder
    }
    return nil
}
```

#### 1.2 Order Service Implementation

```go
// domain/order/service.go
type orderService struct {
    repo domain.OrderRepository
}

func (s *orderService) ValidateOrder(order *Order) error {
    // Basic validation
    if err := order.Validate(); err != nil {
        return err
    }

    // Business rule validations
    if err := s.validateOrderAmount(order); err != nil {
        return err
    }

    return nil
}

func (s *orderService) CalculateTotal(items []OrderItem) float64 {
    var total float64
    for _, item := range items {
        total += item.Price * float64(item.Quantity)
    }
    return total
}
```

### 2. Infrastructure Layer

#### 2.1 MongoDB Repository Implementation

```go
// infrastructure/mongodb/order_repository.go
type mongoOrderRepository struct {
    *mongodb.BaseAdapter[*Order]
}

func (r *mongoOrderRepository) Create(ctx context.Context, order *Order) error {
    return r.GetWriteDB().Insert(ctx, "orders", order)
}

func (r *mongoOrderRepository) FindByID(ctx context.Context, id string) (*Order, error) {
    var order Order
    err := r.GetReadDB().QueryOne(ctx, "orders",
        bson.M{"_id": id}, &order)
    return &order, err
}
```

#### 2.2 Outbox Pattern Implementation

```go
// infrastructure/outbox/repository.go
type OutboxMessage struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    Topic     string            `bson:"topic"`
    Key       string            `bson:"key"`
    Payload   []byte            `bson:"payload"`
    Status    MessageStatus     `bson:"status"`
    CreatedAt time.Time         `bson:"created_at"`
    UpdatedAt time.Time         `bson:"updated_at"`
}

type outboxRepository struct {
    *mongodb.BaseAdapter[*OutboxMessage]
}

func (r *outboxRepository) Create(ctx context.Context, msg *OutboxMessage) error {
    msg.Status = MessageStatusPending
    msg.CreatedAt = time.Now()
    return r.GetWriteDB().Insert(ctx, "outbox", msg)
}
```

### 3. Application Layer (Use Cases)

#### 3.1 Create Order Use Case

```go
// application/order/create_order.go
type CreateOrderRequest struct {
    CustomerID  string      `json:"customer_id"`
    Items      []OrderItem `json:"items"`
    VoucherCode *string     `json:"voucher_code,omitempty"`
}

func (uc *CreateOrderUseCase) Execute(ctx context.Context, req *CreateOrderRequest) error {
    // 1. Validate customer
    if err := uc.validateCustomer(ctx, req.CustomerID); err != nil {
        return err
    }

    // 2. Create order entity
    order := &Order{
        CustomerID: req.CustomerID,
        Items:     req.Items,
        Status:    OrderStatusPending,
        CreatedAt: time.Now(),
    }

    // 3. Apply voucher if present
    if req.VoucherCode != nil {
        if err := uc.applyVoucher(ctx, order, *req.VoucherCode); err != nil {
            return err
        }
    }

    // 4. Execute transaction
    return uc.transactionMgr.ExecuteInTx(ctx, func(ctx context.Context) error {
        // Create order
        if err := uc.orderRepo.Create(ctx, order); err != nil {
            return err
        }

        // Create outbox message
        msg := &OutboxMessage{
            Topic:   "order_created",
            Key:     order.ID.Hex(),
            Payload: createOrderCreatedEvent(order),
        }
        return uc.outboxRepo.Create(ctx, msg)
    })
}
```

#### 3.2 Outbox Message Processor

```go
// application/order/outbox_processor.go
type OutboxProcessor struct {
    outboxRepo    domain.OutboxRepository
    kafkaProducer domain.MessageProducer
}

func (p *OutboxProcessor) ProcessPendingMessages(ctx context.Context) error {
    messages, err := p.outboxRepo.GetPendingMessages(ctx)
    if err != nil {
        return err
    }

    for _, msg := range messages {
        if err := p.processMessage(ctx, msg); err != nil {
            log.Printf("Failed to process message %s: %v", msg.ID, err)
            continue
        }
    }
    return nil
}

func (p *OutboxProcessor) processMessage(ctx context.Context, msg *OutboxMessage) error {
    // Publish to Kafka
    if err := p.kafkaProducer.Publish(msg.Topic, msg.Key, msg.Payload); err != nil {
        return err
    }

    // Mark as published
    return p.outboxRepo.MarkAsPublished(ctx, msg.ID)
}
```

### 4. External Service Integration

#### 4.1 Customer Service Client

```go
// infrastructure/grpc/customer_client.go
type customerClient struct {
    client pb.CustomerServiceClient
}

func (c *customerClient) ValidateCustomer(ctx context.Context, customerID string) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    _, err := c.client.GetCustomer(ctx, &pb.GetCustomerRequest{
        UserId: customerID,
    })
    return err
}
```

### 5. Error Handling

#### 5.1 Domain Errors

```go
// domain/order/errors.go
var (
    ErrCustomerIDRequired = errors.New("customer ID is required")
    ErrEmptyOrder        = errors.New("order must contain at least one item")
    ErrInvalidVoucher    = errors.New("invalid voucher code")
    ErrInsufficientStock = errors.New("insufficient stock")
)
```

#### 5.2 Application Errors

```go
// application/order/errors.go
type OrderError struct {
    Code    string
    Message string
    Cause   error
}

func (e *OrderError) Error() string {
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

var (
    ErrCustomerNotFound = &OrderError{
        Code:    "CUSTOMER_NOT_FOUND",
        Message: "customer not found",
    }
    ErrVoucherExpired = &OrderError{
        Code:    "VOUCHER_EXPIRED",
        Message: "voucher has expired",
    }
)
```

## Transaction Management

### Why Transactions?

1. **Data Consistency**

    - Order creation and inventory updates must be atomic
    - Outbox message must be created in the same transaction

2. **Rollback Support**
    - Automatic rollback on errors
    - Prevents partial updates

### Implementation

```go
// infrastructure/mongodb/transaction.go
type TransactionManager struct {
    client *mongo.Client
}

func (tm *TransactionManager) ExecuteInTx(ctx context.Context, fn func(context.Context) error) error {
    session, err := tm.client.StartSession()
    if err != nil {
        return err
    }
    defer session.EndSession(ctx)

    _, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
        return nil, fn(sessCtx)
    })
    return err
}
```

## Monitoring and Observability

### 1. Metrics

```go
// infrastructure/metrics/order_metrics.go
type OrderMetrics struct {
    orderCreationLatency   prometheus.Histogram
    orderCreationFailures  prometheus.Counter
    outboxMessageCount     prometheus.Gauge
}

func (m *OrderMetrics) RecordOrderCreation(duration time.Duration) {
    m.orderCreationLatency.Observe(duration.Seconds())
}
```

### 2. Logging

```go
// infrastructure/logging/order_logger.go
type OrderLogger struct {
    logger *zap.Logger
}

func (l *OrderLogger) LogOrderCreation(ctx context.Context, order *Order) {
    l.logger.Info("order created",
        zap.String("order_id", order.ID.Hex()),
        zap.String("customer_id", order.CustomerID),
        zap.Float64("total_amount", order.TotalAmount),
    )
}
```

## Testing Strategy

### 1. Unit Tests

```go
// domain/order/service_test.go
func TestOrderService_ValidateOrder(t *testing.T) {
    tests := []struct {
        name    string
        order   *Order
        wantErr error
    }{
        {
            name: "valid order",
            order: &Order{
                CustomerID: "123",
                Items: []OrderItem{{
                    ProductID: "456",
                    Quantity:  1,
                    Price:     10.0,
                }},
            },
            wantErr: nil,
        },
        // ... more test cases
    }
    // ... test implementation
}
```

### 2. Integration Tests

```go
// tests/integration/order_creation_test.go
func TestCreateOrder_E2E(t *testing.T) {
    // Setup test dependencies
    ctx := context.Background()
    useCase := setupTestUseCase(t)

    // Create test request
    req := &CreateOrderRequest{
        CustomerID: "test_customer",
        Items: []OrderItem{{
            ProductID: "test_product",
            Quantity:  1,
            Price:     10.0,
        }},
    }

    // Execute test
    err := useCase.Execute(ctx, req)
    require.NoError(t, err)

    // Verify results
    // ... verification steps
}
```

## Deployment Considerations

1. **Database Indexes**

```javascript
// mongodb/indexes.js
db.orders.createIndex({ customer_id: 1 });
db.orders.createIndex({ created_at: 1 });
db.outbox.createIndex({ status: 1, created_at: 1 });
```

2. **Resource Requirements**

-   MongoDB:
    -   Replica set for high availability
    -   Sufficient storage for orders and outbox messages
-   Kafka:
    -   Multiple partitions for parallel processing
    -   Retention policy for order events
-   Application:
    -   CPU: 2 cores minimum
    -   Memory: 4GB minimum
    -   Storage: 20GB minimum

## Security Considerations

1. **Authentication**

    - Customer service gRPC calls must be authenticated
    - MongoDB connections must use authentication

2. **Authorization**

    - Validate customer can create orders
    - Check order amount limits

3. **Data Protection**
    - Encrypt sensitive order data
    - Mask customer information in logs

## Open Questions and TODOs

1. **Performance**

    - [ ] Implement caching for customer validation
    - [ ] Add database query optimization
    - [ ] Consider bulk operations for batch processing

2. **Reliability**

    - [ ] Add circuit breakers for external services
    - [ ] Implement retry policies
    - [ ] Add dead letter queues

3. **Monitoring**
    - [ ] Add detailed transaction tracing
    - [ ] Set up alerting thresholds
    - [ ] Create monitoring dashboards
