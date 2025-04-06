package notification

type Repository interface {
	GetNotification(id string) (*Notification, error)
	CreateNotification(notification *Notification) error
	MarkAsReadNotification(id string) error
}
