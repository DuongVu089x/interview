package customer

import (
	"testing"
	"time"

	domainnotification "github.com/DuongVu089x/interview/customer/domain/notification"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMapperToEntity(t *testing.T) {
	mapper := NewMapper()

	request := CreateNotificationRequest{
		UserID:      "user123",
		Topic:       "ANNOUNCEMENT",
		Title:       "Test Title",
		Description: "Test Description",
		Link:        "https://example.com",
	}

	entity := mapper.ToEntity(request)

	assert.Equal(t, request.UserID, entity.UserID)
	assert.Equal(t, request.Topic, entity.Topic)
	assert.Equal(t, request.Title, entity.Title)
	assert.Equal(t, request.Description, entity.Description)
	assert.Equal(t, request.Link, entity.Link)
	assert.NotNil(t, entity.IsRead)
	assert.False(t, *entity.IsRead)
}

func TestMapperToDTO(t *testing.T) {
	mapper := NewMapper()
	now := time.Now()
	objID := primitive.NewObjectID()
	isRead := false

	entity := &domainnotification.Notification{
		ID:              &objID,
		CreatedTime:     &now,
		LastUpdatedTime: &now,
		Code:            "code123",
		UserID:          "user123",
		IsRead:          &isRead,
		Topic:           "ANNOUNCEMENT",
		Title:           "Test Title",
		Description:     "Test Description",
		Link:            "https://example.com",
	}

	dto := mapper.ToDTO(entity)

	assert.Equal(t, objID.Hex(), dto.ID)
	assert.Equal(t, now, dto.CreatedTime)
	assert.Equal(t, now, dto.LastUpdatedTime)
	assert.Equal(t, entity.Code, dto.Code)
	assert.Equal(t, entity.UserID, dto.UserID)
	assert.Equal(t, *entity.IsRead, dto.IsRead)
	assert.Equal(t, entity.Topic, dto.Topic)
	assert.Equal(t, entity.Title, dto.Title)
	assert.Equal(t, entity.Description, dto.Description)
	assert.Equal(t, entity.Link, dto.Link)
}

func TestMapperToDTOList(t *testing.T) {
	mapper := NewMapper()
	now := time.Now()
	objID1 := primitive.NewObjectID()
	objID2 := primitive.NewObjectID()
	isRead := false

	entities := []*domainnotification.Notification{
		{
			ID:          &objID1,
			CreatedTime: &now,
			UserID:      "user123",
			IsRead:      &isRead,
			Topic:       "ANNOUNCEMENT",
			Title:       "Title 1",
			Description: "Description 1",
		},
		{
			ID:          &objID2,
			CreatedTime: &now,
			UserID:      "user456",
			IsRead:      &isRead,
			Topic:       "EVENT",
			Title:       "Title 2",
			Description: "Description 2",
		},
	}

	dtos := mapper.ToDTOList(entities)

	assert.Equal(t, len(entities), len(dtos))
	assert.Equal(t, objID1.Hex(), dtos[0].ID)
	assert.Equal(t, entities[0].UserID, dtos[0].UserID)
	assert.Equal(t, objID2.Hex(), dtos[1].ID)
	assert.Equal(t, entities[1].UserID, dtos[1].UserID)
}

func TestMapperToEntityFromDTO(t *testing.T) {
	mapper := NewMapper()
	now := time.Now()
	objID := primitive.NewObjectID().Hex()

	dto := &NotificationDTO{
		ID:              objID,
		CreatedTime:     now,
		LastUpdatedTime: now,
		Code:            "code123",
		UserID:          "user123",
		IsRead:          true,
		Topic:           "ANNOUNCEMENT",
		Title:           "Test Title",
		Description:     "Test Description",
		Link:            "https://example.com",
	}

	entity := mapper.ToEntityFromDTO(dto)

	assert.Equal(t, objID, entity.ID.Hex())
	assert.Equal(t, now, *entity.CreatedTime)
	assert.Equal(t, now, *entity.LastUpdatedTime)
	assert.Equal(t, dto.Code, entity.Code)
	assert.Equal(t, dto.UserID, entity.UserID)
	assert.Equal(t, dto.IsRead, *entity.IsRead)
	assert.Equal(t, dto.Topic, entity.Topic)
	assert.Equal(t, dto.Title, entity.Title)
	assert.Equal(t, dto.Description, entity.Description)
	assert.Equal(t, dto.Link, entity.Link)
}

func TestMapperNilHandling(t *testing.T) {
	mapper := NewMapper()

	// Test nil handling in ToDTO
	assert.Nil(t, mapper.ToDTO(nil))

	// Test nil handling in ToDTOList
	assert.Nil(t, mapper.ToDTOList(nil))

	// Test empty list handling in ToDTOList
	emptyList := make([]*domainnotification.Notification, 0)
	assert.Equal(t, 0, len(mapper.ToDTOList(emptyList)))

	// Test nil handling in ToEntityFromDTO
	assert.Nil(t, mapper.ToEntityFromDTO(nil))
}
