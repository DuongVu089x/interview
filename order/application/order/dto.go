package order

import (
	"time"
)

// DTOs (Data Transfer Objects) for input/output
type CreateOrderRequest struct {
	OrderID string    `json:"orderId" validate:"omitempty"`
	UserID  string    `json:"userId" validate:"required"`
	Items   []ItemDTO `json:"items" validate:"required,dive,required"`
}

type ItemDTO struct {
	ProductID string  `json:"productId" validate:"required"`
	Quantity  int     `json:"quantity" validate:"required,gt=0"`
	Price     float64 `json:"price" validate:"required,gt=0"`
}

type OrderResponse struct {
	OrderID     int64     `json:"orderId"`
	OrderCode   string    `json:"orderCode"`
	UserID      string    `json:"userId"`
	Items       []ItemDTO `json:"items"`
	TotalAmount float64   `json:"totalAmount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}

// GetOrdersByUserIDRequest defines request parameters for getting orders by user ID
type GetOrdersByUserIDRequest struct {
	UserID string `json:"userId,omitempty" validate:"omitempty"`
	Status string `json:"status,omitempty" validate:"omitempty"`
}

// OrderListResponse represents a collection of orders
type OrderListResponse struct {
	Orders []OrderResponse `json:"orders,omitempty"`
	Count  int             `json:"count,omitempty"`
}
