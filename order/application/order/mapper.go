package order

import domainorder "github.com/DuongVu089x/interview/order/domain/order"

// Mapper converts between DTOs and Domain entities
type Mapper struct{}

func (m *Mapper) ToEntity(dto CreateOrderRequest) *domainorder.Order {
	return &domainorder.Order{
		UserID: dto.UserID,
		Items:  m.toOrderItems(dto.Items),
		Status: domainorder.StatusPending,
	}
}

func (m *Mapper) toOrderItems(dtos []ItemDTO) []domainorder.OrderItem {
	items := make([]domainorder.OrderItem, len(dtos))
	for i, dto := range dtos {
		items[i] = domainorder.OrderItem{
			ProductID: dto.ProductID,
			Quantity:  dto.Quantity,
			Price:     dto.Price,
		}
	}
	return items
}

func (m *Mapper) ToResponse(order *domainorder.Order) OrderResponse {
	return OrderResponse{
		OrderID:     order.OrderID,
		OrderCode:   order.OrderCode,
		UserID:      order.UserID,
		Items:       m.toItemDTOs(order.Items),
		TotalAmount: order.TotalAmount,
		Status:      string(order.Status),
		CreatedAt:   order.CreatedAt,
	}
}

func (m *Mapper) toItemDTOs(items []domainorder.OrderItem) []ItemDTO {
	dto := make([]ItemDTO, len(items))
	for i, item := range items {
		dto[i] = ItemDTO{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}
	return dto
}
