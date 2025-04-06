package order

// Service defines the business operations for orders
type Service interface {
	ValidateOrder(order *Order) error
	CalculateTotal(items []OrderItem) float64

	// Calculate total of customer
	CalculateTotalOfCustomer(customerID string, status OrderStatus) (float64, error)

	GetOrderByCustomerID(customerID string, conditions map[string]any) ([]Order, error)
	GetOrder(id string) (*Order, error)

	CreateOrder(order *Order) error
	UpdateOrder(order *Order) error
	DeleteOrder(id string) error
}
