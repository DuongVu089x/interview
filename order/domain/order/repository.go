package order

type Repository interface {
	GetOrder(id int64) (*Order, error)
	GetOrders(conditions Order) ([]Order, error)

	CreateOrder(order *Order) error
	UpdateOrder(order *Order) error
	DeleteOrder(id string) error
}
