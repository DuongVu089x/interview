package order

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID          *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	OrderID     int64               `json:"orderId,omitempty" bson:"order_id,omitempty"`
	OrderCode   string              `json:"orderCode,omitempty" bson:"order_code,omitempty"`
	UserID      string              `json:"userId,omitempty" bson:"user_id,omitempty"`
	Items       []OrderItem         `json:"items,omitempty" bson:"items,omitempty"`
	TotalAmount float64             `json:"totalAmount,omitempty" bson:"total_amount,omitempty"`
	Status      OrderStatus         `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt   time.Time           `json:"createdAt,omitempty" bson:"created_at,omitempty"`
}

type OrderItem struct {
	ProductID string  `json:"productId,omitempty" bson:"product_id,omitempty"`
	Quantity  int     `json:"quantity,omitempty" bson:"quantity,omitempty"`
	Price     float64 `json:"price,omitempty" bson:"price,omitempty"`
}

type OrderStatus string

const (
	StatusPending OrderStatus = "pending"
	StatusPaid    OrderStatus = "paid"
	StatusShipped OrderStatus = "shipped"
)
