package notification

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, handler *Handler) {
	e.GET("/api/notifications", handler.GetNotifications)
}
