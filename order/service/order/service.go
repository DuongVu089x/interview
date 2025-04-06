package order

import (
	"errors"

	domainorder "github.com/DuongVu089x/interview/order/domain/order"
)

type Service struct {
	orderRepo domainorder.Repository
}

func NewOrderService(orderRepo domainorder.Repository) domainorder.Service {
	return &Service{orderRepo: orderRepo}
}

func (s *Service) GetOrderByCustomerID(customerID string, conditions map[string]interface{}) ([]domainorder.Order, error) {
	return s.orderRepo.GetOrderByCustomerID(customerID, conditions)
}

func (s *Service) GetOrder(id string) (*domainorder.Order, error) {
	return s.orderRepo.GetOrder(id)
}

func (s *Service) CalculateTotal(items []domainorder.OrderItem) float64 {
	total := 0.0
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

func (s *Service) CreateOrder(order *domainorder.Order) error {
	return s.orderRepo.CreateOrder(order)
}

func (s *Service) UpdateOrder(order *domainorder.Order) error {
	return s.orderRepo.UpdateOrder(order)
}

func (s *Service) DeleteOrder(id string) error {
	return s.orderRepo.DeleteOrder(id)
}

func (s *Service) CalculateTotalOfCustomer(customerID string, status domainorder.OrderStatus) (float64, error) {
	total := 0.0

	orders, err := s.GetOrderByCustomerID(customerID, map[string]interface{}{"status": status})
	if err != nil {
		return 0, err
	}

	for _, order := range orders {
		total += order.TotalAmount
	}
	return total, nil
}

func (s *Service) ValidateOrder(order *domainorder.Order) error {
	if order.UserID == "" {
		return errors.New("user ID is required")
	}

	if len(order.Items) == 0 {
		return errors.New("items are required")
	}

	for _, item := range order.Items {
		if item.ProductID == "" {
			return errors.New("product ID is required")
		}

		if item.Quantity <= 0 {
			return errors.New("quantity must be greater than 0")
		}

	}

	if order.TotalAmount <= 0 {
		return errors.New("total amount must be greater than 0")
	}

	return nil
}
