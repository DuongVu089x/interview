package customer

import (
	domainnotification "github.com/DuongVu089x/interview/customer/domain/notification"
)

type UseCase struct {
	notificationRepository domainnotification.Repository
	notifier               NotificationDispatcher
}

func NewUseCase(notificationRepository domainnotification.Repository, notifier NotificationDispatcher) *UseCase {
	return &UseCase{notificationRepository: notificationRepository, notifier: notifier}
}

type NotificationDispatcher interface {
	DispatchToUser(userID string, notification *domainnotification.Notification) error
}

func (u *UseCase) GetNotification(id string) (*domainnotification.Notification, error) {
	return u.notificationRepository.GetNotification(id)
}

func (u *UseCase) MarkAsReadNotification(id string) error {
	return u.notificationRepository.MarkAsReadNotification(id)
}

func (u *UseCase) CreateNotification(notification *domainnotification.Notification) error {
	err := u.notificationRepository.CreateNotification(notification)
	if err != nil {
		return err
	}

	return u.notifier.DispatchToUser(notification.UserId, notification)
}
