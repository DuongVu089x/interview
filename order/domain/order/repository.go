package order

type Repository interface {
	GetOrderByCustomerID(customerID string, conditions map[string]any) ([]Order, error)
	GetOrder(id string) (*Order, error)

	CreateOrder(order *Order) error
	UpdateOrder(order *Order) error
	DeleteOrder(id string) error
}
