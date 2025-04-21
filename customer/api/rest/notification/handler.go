package notification

import (
	"net/http"
	"strconv"

	notificationusecase "github.com/DuongVu089x/interview/customer/application/notification"
	"github.com/DuongVu089x/interview/customer/component/appctx"
	notificationrepository "github.com/DuongVu089x/interview/customer/repository/notification"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	appCtx              appctx.AppContext
	notificationUseCase *notificationusecase.ReadUseCase
}

func NewHandler(appCtx appctx.AppContext) *Handler {
	notificationRepo := notificationrepository.NewMongoRepository(appCtx.GetMainDBConnection(), appCtx.GetReadMainDBConnection())
	notificationUseCase := notificationusecase.NewReadUseCase(notificationRepo)

	return &Handler{
		appCtx:              appCtx,
		notificationUseCase: notificationUseCase,
	}
}

// GetNotifications handles retrieving notifications for a user
func (h *Handler) GetNotifications(c echo.Context) error {
	userId := c.QueryParam("userId")
	if userId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID is required")
	}

	page, _ := strconv.ParseInt(c.QueryParam("page"), 10, 64)
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.ParseInt(c.QueryParam("limit"), 10, 64)
	if limit < 1 {
		limit = 10
	}

	response, err := h.notificationUseCase.GetNotifications(c.Request().Context(), userId, page, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get notifications")
	}

	return c.JSON(http.StatusOK, response)
}
