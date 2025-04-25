package notification

import (
	"net/http"
	"strconv"

	notificationusecase "github.com/DuongVu089x/interview/customer/application/notification"
	"github.com/DuongVu089x/interview/customer/component/appctx"
	"github.com/DuongVu089x/interview/customer/middleware"
	notificationrepository "github.com/DuongVu089x/interview/customer/repository/notification"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Handler struct {
	appCtx              appctx.AppContext
	notificationUseCase *notificationusecase.ReadUseCase
	tracer              trace.Tracer
}

func NewHandler(appCtx appctx.AppContext) *Handler {
	notificationRepo := notificationrepository.NewMongoRepository(appCtx.GetMainDBConnection(), appCtx.GetReadMainDBConnection())
	notificationUseCase := notificationusecase.NewReadUseCase(notificationRepo)

	return &Handler{
		appCtx:              appCtx,
		notificationUseCase: notificationUseCase,
		tracer:              appCtx.GetTracer(),
	}
}

// GetNotifications handles retrieving notifications for a user
func (h *Handler) GetNotifications(c echo.Context) error {
	ctx := c.Request().Context()
	logger := middleware.GetRequestLogger(c)

	// Single span for the entire handler operation
	ctx, span := h.tracer.Start(ctx, "notification.GetNotifications")
	defer span.End()

	// Get and validate parameters
	userId := c.QueryParam("userId")
	if userId == "" {
		logger.Warn("Missing user ID in request")
		return echo.NewHTTPError(http.StatusBadRequest, "User ID is required")
	}

	page, err := strconv.ParseInt(c.QueryParam("page"), 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.ParseInt(c.QueryParam("limit"), 10, 64)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Add essential request context
	span.SetAttributes(
		attribute.String("user_id", userId),
		attribute.Int64("page", page),
		attribute.Int64("limit", limit),
	)

	// Execute use case
	response, err := h.notificationUseCase.GetNotifications(ctx, logger, userId, page, limit)
	if err != nil {
		logger.Error("Failed to get notifications",
			zap.String("user_id", userId),
			zap.Error(err),
		)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get notifications")
	}

	logger.Info("Successfully retrieved notifications",
		zap.String("user_id", userId),
		zap.Int("count", len(response.Notifications)),
	)

	return c.JSON(http.StatusOK, response)
}
