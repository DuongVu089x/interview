package customer

import (
	domainnotification "github.com/DuongVu089x/interview/customer/domain/notification"
)

// ReadUseCase handles read operations that don't require the dispatcher
type ReadUseCase struct {
	notificationRepository domainnotification.Repository
}

func NewReadUseCase(notificationRepository domainnotification.Repository) *ReadUseCase {
	return &ReadUseCase{notificationRepository: notificationRepository}
}

// GetNotificationsResponse represents the response for getting notifications
type GetNotificationsResponse struct {
	Notifications []*domainnotification.Notification `json:"notifications"`
	Total         int64                              `json:"total"`
}

func (u *ReadUseCase) GetNotifications(userId string, page, limit int64) (*GetNotificationsResponse, error) {
	offset := (page - 1) * limit
	notifications, err := u.notificationRepository.GetNotifications(userId, offset, limit)
	if err != nil {
		return nil, err
	}

	return &GetNotificationsResponse{
		Notifications: notifications,
		Total:         int64(len(notifications)),
	}, nil
}

func (u *ReadUseCase) GetNotification(id string) (*domainnotification.Notification, error) {
	return u.notificationRepository.GetNotification(id)
}

func (u *ReadUseCase) MarkAsReadNotification(id string) error {
	return u.notificationRepository.MarkAsReadNotification(id)
}

// WriteUseCase handles operations that require the dispatcher (like creating notifications)
type WriteUseCase struct {
	notificationRepository domainnotification.Repository
	notifier               NotificationDispatcher
}

type NotificationDispatcher interface {
	DispatchToUser(userID string, notification *domainnotification.Notification) error
}

func NewWriteUseCase(notificationRepository domainnotification.Repository, notifier NotificationDispatcher) *WriteUseCase {
	return &WriteUseCase{
		notificationRepository: notificationRepository,
		notifier:               notifier,
	}
}

func (u *WriteUseCase) CreateNotification(notification *domainnotification.Notification) error {
	err := u.notificationRepository.CreateNotification(notification)
	if err != nil {
		return err
	}

	return u.notifier.DispatchToUser(notification.UserId, notification)
}
