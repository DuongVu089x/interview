package order

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, handler *Handler) {
	e.POST("/order", handler.CreateOrder)
	e.GET("/order/:id", handler.GetOrder)
}
