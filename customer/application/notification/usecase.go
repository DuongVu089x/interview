package customer

import (
	"context"

	domainnotification "github.com/DuongVu089x/interview/customer/domain/notification"
)

// ReadUseCase handles read operations that don't require the dispatcher
type ReadUseCase struct {
	notificationRepository domainnotification.Repository
	mapper                 *Mapper
}

func NewReadUseCase(notificationRepository domainnotification.Repository) *ReadUseCase {
	return &ReadUseCase{
		notificationRepository: notificationRepository,
		mapper:                 NewMapper(),
	}
}

// GetNotificationsResponse represents the response for getting notifications
type GetNotificationsResponse struct {
	Notifications []*NotificationDTO `json:"notifications"`
	Total         int64              `json:"total"`
}

func (u *ReadUseCase) GetNotifications(ctx context.Context, userId string, page, limit int64) (*GetNotificationsResponse, error) {
	offset := (page - 1) * limit
	notifications, err := u.notificationRepository.GetNotifications(ctx, userId, offset, limit)
	if err != nil {
		return nil, err
	}

	// Use mapper to convert domain entities to DTOs
	dtos := u.mapper.ToDTOList(notifications)

	return &GetNotificationsResponse{
		Notifications: dtos,
		Total:         int64(len(notifications)),
	}, nil
}

func (u *ReadUseCase) GetNotification(ctx context.Context, id string) (*NotificationDTO, error) {
	notification, err := u.notificationRepository.GetNotification(ctx, id)
	if err != nil {
		return nil, err
	}
	return u.mapper.ToDTO(notification), nil
}

func (u *ReadUseCase) MarkAsReadNotification(ctx context.Context, id string) error {
	return u.notificationRepository.MarkAsReadNotification(ctx, id)
}

// WriteUseCase handles operations that require the dispatcher (like creating notifications)
type WriteUseCase struct {
	notificationRepository domainnotification.Repository
	notifier               NotificationDispatcher
	mapper                 *Mapper
}

// NotificationDispatcher interface defines how to dispatch notifications
type NotificationDispatcher interface {
	DispatchToUser(userID string, notification *NotificationDTO) error
}

func NewWriteUseCase(notificationRepository domainnotification.Repository, notifier NotificationDispatcher) *WriteUseCase {
	return &WriteUseCase{
		notificationRepository: notificationRepository,
		notifier:               notifier,
		mapper:                 NewMapper(),
	}
}

func (u *WriteUseCase) CreateNotification(ctx context.Context, request *CreateNotificationRequest) error {
	// Use mapper to convert DTO to domain entity
	notification := u.mapper.ToEntity(*request)

	notification, err := u.notificationRepository.CreateNotification(ctx, notification)
	if err != nil {
		return err
	}

	// Convert back to DTO for dispatching
	notificationDTO := u.mapper.ToDTO(notification)
	return u.notifier.DispatchToUser(notification.UserID, notificationDTO)
}
