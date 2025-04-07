package notification

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetNotification(id string) (*Notification, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Notification), args.Error(1)
}

func (m *MockRepository) GetNotifications(userId string, offset, limit int64) ([]*Notification, error) {
	args := m.Called(userId, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Notification), args.Error(1)
}

func (m *MockRepository) CreateNotification(notification *Notification) error {
	args := m.Called(notification)
	return args.Error(0)
}

func (m *MockRepository) MarkAsReadNotification(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGetNotification(t *testing.T) {
	mockRepo := new(MockRepository)
	now := time.Now()
	isRead := false
	objID := primitive.NewObjectID()

	expectedNotification := &Notification{
		ID:          &objID,
		CreatedTime: &now,
		UserID:      "user123",
		IsRead:      &isRead,
		Topic:       string(Topic.ANNOUNCEMENT),
		Title:       "Test Notification",
		Description: "This is a test notification",
	}

	// Test successful retrieval
	mockRepo.On("GetNotification", "existing_id").Return(expectedNotification, nil)
	mockRepo.On("GetNotification", "non_existing_id").Return(nil, errors.New("notification not found"))

	t.Run("Success - Get Existing Notification", func(t *testing.T) {
		notification, err := mockRepo.GetNotification("existing_id")
		assert.NoError(t, err)
		assert.Equal(t, expectedNotification, notification)
	})

	t.Run("Error - Get Non-existing Notification", func(t *testing.T) {
		notification, err := mockRepo.GetNotification("non_existing_id")
		assert.Error(t, err)
		assert.Nil(t, notification)
		assert.Equal(t, "notification not found", err.Error())
	})
}

func TestGetNotifications(t *testing.T) {
	mockRepo := new(MockRepository)
	now := time.Now()
	isRead := false
	objID1 := primitive.NewObjectID()
	objID2 := primitive.NewObjectID()

	notifications := []*Notification{
		{
			ID:          &objID1,
			CreatedTime: &now,
			UserID:      "user123",
			IsRead:      &isRead,
			Topic:       string(Topic.ANNOUNCEMENT),
			Title:       "Test Notification 1",
		},
		{
			ID:          &objID2,
			CreatedTime: &now,
			UserID:      "user123",
			IsRead:      &isRead,
			Topic:       string(Topic.PROMOTION),
			Title:       "Test Notification 2",
		},
	}

	mockRepo.On("GetNotifications", "user123", int64(0), int64(10)).Return(notifications, nil)
	mockRepo.On("GetNotifications", "non_existing_user", int64(0), int64(10)).Return([]*Notification{}, nil)

	t.Run("Success - Get User Notifications", func(t *testing.T) {
		result, err := mockRepo.GetNotifications("user123", 0, 10)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, notifications, result)
	})

	t.Run("Success - Get Empty Notifications", func(t *testing.T) {
		result, err := mockRepo.GetNotifications("non_existing_user", 0, 10)
		assert.NoError(t, err)
		assert.Empty(t, result)
	})
}

func TestCreateNotification(t *testing.T) {
	mockRepo := new(MockRepository)
	now := time.Now()
	isRead := false
	objID := primitive.NewObjectID()

	notification := &Notification{
		ID:          &objID,
		CreatedTime: &now,
		UserID:      "user123",
		IsRead:      &isRead,
		Topic:       string(Topic.ANNOUNCEMENT),
		Title:       "Test Notification",
		Description: "This is a test notification",
	}

	mockRepo.On("CreateNotification", notification).Return(nil)
	mockRepo.On("CreateNotification", mock.AnythingOfType("*notification.Notification")).Return(errors.New("invalid notification"))

	t.Run("Success - Create Valid Notification", func(t *testing.T) {
		err := mockRepo.CreateNotification(notification)
		assert.NoError(t, err)
	})

	t.Run("Error - Create Invalid Notification", func(t *testing.T) {
		invalidNotification := &Notification{}
		err := mockRepo.CreateNotification(invalidNotification)
		assert.Error(t, err)
		assert.Equal(t, "invalid notification", err.Error())
	})
}

func TestMarkAsReadNotification(t *testing.T) {
	mockRepo := new(MockRepository)

	mockRepo.On("MarkAsReadNotification", "existing_id").Return(nil)
	mockRepo.On("MarkAsReadNotification", "non_existing_id").Return(errors.New("notification not found"))

	t.Run("Success - Mark Existing Notification as Read", func(t *testing.T) {
		err := mockRepo.MarkAsReadNotification("existing_id")
		assert.NoError(t, err)
	})

	t.Run("Error - Mark Non-existing Notification as Read", func(t *testing.T) {
		err := mockRepo.MarkAsReadNotification("non_existing_id")
		assert.Error(t, err)
		assert.Equal(t, "notification not found", err.Error())
	})
}
