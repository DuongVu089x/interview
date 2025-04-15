package customer

import "context"

// Repository defines the interface for customer data access
type Repository interface {
	GetCustomer(ctx context.Context, userId string) (*Customer, error)
	CreateCustomer(ctx context.Context, customer *Customer) (*Customer, error)
	UpdateCustomer(ctx context.Context, customer *Customer) error
	DeleteCustomer(ctx context.Context, id string) error
}
