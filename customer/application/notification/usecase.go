package customer

import (
	"context"

	domainnotification "github.com/DuongVu089x/interview/customer/domain/notification"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ReadUseCase handles read operations that don't require the dispatcher
type ReadUseCase struct {
	notificationRepository domainnotification.Repository
	mapper                 *Mapper
	tracer                 trace.Tracer
}

func NewReadUseCase(notificationRepository domainnotification.Repository) *ReadUseCase {
	return &ReadUseCase{
		notificationRepository: notificationRepository,
		mapper:                 NewMapper(),
		tracer:                 otel.GetTracerProvider().Tracer("notification.usecase"),
	}
}

// GetNotificationsResponse represents the response for getting notifications
type GetNotificationsResponse struct {
	Notifications []*NotificationDTO `json:"notifications"`
	Total         int64              `json:"total"`
}

func (u *ReadUseCase) GetNotifications(ctx context.Context, logger *zap.Logger, userId string, page, limit int64) (*GetNotificationsResponse, error) {
	// Create span for database operation
	ctx, span := u.tracer.Start(ctx, "notification.repository.GetNotifications")
	defer span.End()

	offset := (page - 1) * limit
	span.SetAttributes(
		attribute.String("user_id", userId),
		attribute.Int64("offset", offset),
		attribute.Int64("limit", limit),
	)

	notifications, err := u.notificationRepository.GetNotifications(ctx, userId, offset, limit)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		return nil, err
	}

	// Convert domain entities to DTOs
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
