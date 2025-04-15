package notification

import "context"

type Repository interface {
	GetNotification(ctx context.Context, id string) (*Notification, error)
	GetNotifications(ctx context.Context, userId string, offset, limit int64) ([]*Notification, error)
	CreateNotification(ctx context.Context, notification *Notification) (*Notification, error)
	MarkAsReadNotification(ctx context.Context, id string) error
}
