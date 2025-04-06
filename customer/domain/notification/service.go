package notification

type Service interface {
	MarkAsReadNotification(id string) error
}
