package customer

import (
	"time"

	domainnotification "github.com/DuongVu089x/interview/customer/domain/notification"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NotificationDTO represents the data transfer object for notifications
type NotificationDTO struct {
	ID              string    `json:"id,omitempty"`
	CreatedTime     time.Time `json:"createdTime,omitempty"`
	LastUpdatedTime time.Time `json:"lastUpdatedTime,omitempty"`
	Code            string    `json:"code,omitempty"`
	UserID          string    `json:"userId,omitempty"`
	IsRead          bool      `json:"isRead,omitempty"`
	Topic           string    `json:"topic,omitempty"`
	Title           string    `json:"title,omitempty"`
	Description     string    `json:"description,omitempty"`
	Link            string    `json:"link,omitempty"`
}

// CreateNotificationRequest represents the request for creating notifications
type CreateNotificationRequest struct {
	UserID      string `json:"userId,omitempty"`
	Topic       string `json:"topic,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Link        string `json:"link,omitempty"`
}

// Mapper converts between DTOs and Domain entities
type Mapper struct{}

// NewMapper creates a new instance of Mapper
func NewMapper() *Mapper {
	return &Mapper{}
}

// ToEntity converts a DTO to a domain entity
func (m *Mapper) ToEntity(dto CreateNotificationRequest) *domainnotification.Notification {
	isRead := false
	return &domainnotification.Notification{
		UserID:      dto.UserID,
		Topic:       dto.Topic,
		Title:       dto.Title,
		Description: dto.Description,
		Link:        dto.Link,
		IsRead:      &isRead,
	}
}

// ToDTO converts a domain entity to a DTO
func (m *Mapper) ToDTO(entity *domainnotification.Notification) *NotificationDTO {
	if entity == nil {
		return nil
	}

	dto := &NotificationDTO{
		UserID:      entity.UserID,
		Topic:       entity.Topic,
		Title:       entity.Title,
		Description: entity.Description,
		Link:        entity.Link,
	}

	if entity.ID != nil {
		dto.ID = entity.ID.Hex()
	}

	if entity.CreatedTime != nil {
		dto.CreatedTime = *entity.CreatedTime
	}

	if entity.LastUpdatedTime != nil {
		dto.LastUpdatedTime = *entity.LastUpdatedTime
	}

	if entity.IsRead != nil {
		dto.IsRead = *entity.IsRead
	}

	if entity.Code != "" {
		dto.Code = entity.Code
	}

	return dto
}

// ToDTOList converts a list of domain entities to a list of DTOs
func (m *Mapper) ToDTOList(entities []*domainnotification.Notification) []*NotificationDTO {
	if entities == nil {
		return nil
	}

	dtos := make([]*NotificationDTO, len(entities))
	for i, entity := range entities {
		dtos[i] = m.ToDTO(entity)
	}
	return dtos
}

// ToEntityFromDTO converts a DTO to a domain entity
func (m *Mapper) ToEntityFromDTO(dto *NotificationDTO) *domainnotification.Notification {
	if dto == nil {
		return nil
	}

	entity := &domainnotification.Notification{
		UserID:      dto.UserID,
		Topic:       dto.Topic,
		Title:       dto.Title,
		Description: dto.Description,
		Link:        dto.Link,
		Code:        dto.Code,
	}

	if dto.ID != "" {
		objectID, err := primitive.ObjectIDFromHex(dto.ID)
		if err == nil {
			entity.ID = &objectID
		}
	}

	// Convert non-zero time values
	if !dto.CreatedTime.IsZero() {
		createdTime := dto.CreatedTime
		entity.CreatedTime = &createdTime
	}

	if !dto.LastUpdatedTime.IsZero() {
		lastUpdatedTime := dto.LastUpdatedTime
		entity.LastUpdatedTime = &lastUpdatedTime
	}

	isRead := dto.IsRead
	entity.IsRead = &isRead

	return entity
}
