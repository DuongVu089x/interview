# Create Order Flow Task Breakdown

## Overview

This task implements a transactional order creation flow with customer validation, voucher/promotion validation, inventory checks, and outbox pattern for event notifications.

## Architecture Components

-   Domain Services
-   Use Cases
-   Repositories
-   Event Handlers
-   External Services Integration

## Task Breakdown

### 1. Domain Layer Implementation (Order Service)

-   [ ] Create Order Validation Service

    -   [ ] Implement order structure validation
    -   [ ] Add business rule validations
    -   [ ] Create quantity validation logic
    -   [ ] Add total amount calculation

-   [ ] Create Inventory Service

    -   [ ] Implement stock check methods
    -   [ ] Add reserve stock functionality
    -   [ ] Create stock update methods

-   [ ] Create Voucher/Promotion Service
    -   [ ] Implement voucher validation logic
    -   [ ] Add promotion calculation methods
    -   [ ] Create discount application rules

### 2. Infrastructure Layer Setup

-   [ ] Create Order Repository

    -   [ ] Implement CRUD operations
    -   [ ] Add transaction support
    -   [ ] Create query methods

-   [ ] Create Outbox Repository

    -   [ ] Implement message storage
    -   [ ] Add message status tracking
    -   [ ] Create query methods for pending messages

-   [ ] Setup Kafka Producer
    -   [ ] Implement message publishing
    -   [ ] Add error handling
    -   [ ] Create retry mechanism

### 3. Application Layer (Use Cases)

-   [ ] Create OrderCreationUseCase

    ```go
    type OrderCreationUseCase struct {
        orderService       domain.OrderService
        customerClient     pb.CustomerServiceClient
        voucherService    domain.VoucherService
        inventoryService  domain.InventoryService
        outboxRepository  domain.OutboxRepository
        transactionMgr    domain.TransactionManager
    }
    ```

-   [ ] Implement Validation Steps

    -   [ ] Customer validation using gRPC
    -   [ ] Voucher/promotion validation
    -   [ ] Stock availability check
    -   [ ] Order total calculation

-   [ ] Implement Transaction Flow

    -   [ ] Begin transaction
    -   [ ] Create order record
    -   [ ] Update inventory
    -   [ ] Create outbox message
    -   [ ] Commit transaction

-   [ ] Implement Event Publishing
    -   [ ] Create order created event
    -   [ ] Setup outbox message processor
    -   [ ] Implement Kafka publishing

### 4. External Service Integration

-   [ ] Setup Customer Service Client

    -   [ ] Implement gRPC client configuration
    -   [ ] Add retry mechanism
    -   [ ] Create timeout handling

-   [ ] Setup Inventory Service Client
    -   [ ] Implement stock check methods
    -   [ ] Add reservation system
    -   [ ] Create rollback mechanism

### 5. Testing

-   [ ] Unit Tests

    -   [ ] Test order validation
    -   [ ] Test voucher validation
    -   [ ] Test inventory checks
    -   [ ] Test transaction flow

-   [ ] Integration Tests
    -   [ ] Test customer service integration
    -   [ ] Test inventory service integration
    -   [ ] Test Kafka publishing
    -   [ ] Test transaction rollback

## Implementation Flow

1. **Validation Phase** (Use Case Layer)

    ```go
    func (uc *OrderCreationUseCase) validateOrder(ctx context.Context, order *Order) error {
        // Customer validation
        if valid, err := uc.customerClient.ValidateCustomer(ctx, order.CustomerID); !valid {
            return ErrInvalidCustomer
        }

        // Voucher validation
        if err := uc.voucherService.ValidateVoucher(ctx, order.VoucherCode); err != nil {
            return ErrInvalidVoucher
        }

        // Stock validation
        if err := uc.inventoryService.CheckStock(ctx, order.Items); err != nil {
            return ErrInsufficientStock
        }

        return nil
    }
    ```

2. **Transaction Phase** (Use Case Layer)

    ```go
    func (uc *OrderCreationUseCase) executeTransaction(ctx context.Context, order *Order) error {
        return uc.transactionMgr.ExecuteInTx(ctx, func(ctx context.Context) error {
            // Create order
            if err := uc.orderRepository.Create(ctx, order); err != nil {
                return err
            }

            // Update inventory
            if err := uc.inventoryService.UpdateStock(ctx, order.Items); err != nil {
                return err
            }

            // Create outbox message
            message := createOrderCreatedMessage(order)
            return uc.outboxRepository.Create(ctx, message)
        })
    }
    ```

3. **Event Publishing Phase** (Background Process)
    ```go
    func (p *OutboxProcessor) processMessages(ctx context.Context) {
        messages, err := p.outboxRepository.GetPendingMessages(ctx)
        for _, msg := range messages {
            if err := p.kafkaProducer.Publish(msg); err != nil {
                continue // Will be retried next cycle
            }
            p.outboxRepository.MarkAsPublished(ctx, msg.ID)
        }
    }
    ```

## Dependencies

-   MongoDB for order storage
-   Kafka for event publishing
-   gRPC for customer service communication
-   Redis for distributed locking (optional)

## Error Handling Strategy

1. **Validation Errors**

    - Return specific error types for each validation failure
    - Include detailed error messages for client feedback

2. **Transaction Errors**

    - Automatic rollback on any error
    - Retry mechanism for transient failures
    - Dead letter queue for failed messages

3. **Publishing Errors**
    - Retry mechanism with exponential backoff
    - Error logging and monitoring
    - Manual intervention process for failed messages

## Monitoring Considerations

-   [ ] Add transaction timing metrics
-   [ ] Monitor outbox message queue length
-   [ ] Track validation failure rates
-   [ ] Monitor Kafka publishing success rate

## Open Questions

1. Should we implement compensating transactions for partial failures?
2. What is the retry strategy for failed outbox messages?
3. How do we handle long-running transactions?
4. What is the timeout strategy for external service calls?
