package order

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, handler *Handler) {
	e.GET("/order/:id", handler.GetOrder)
	e.GET("/user/:userId/orders", handler.GetOrdersByUserID)
	e.POST("/order", handler.CreateOrder)
}
