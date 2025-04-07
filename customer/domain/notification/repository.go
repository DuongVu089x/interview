package notification

type Repository interface {
	GetNotification(id string) (*Notification, error)
	GetNotifications(userId string, offset, limit int64) ([]*Notification, error)
	CreateNotification(notification *Notification) error
	MarkAsReadNotification(id string) error
}
