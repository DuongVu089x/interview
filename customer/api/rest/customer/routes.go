package customer

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, h *RestHandler) {
	g := e.Group("/customers")
	g.GET("/:id", h.HandleGetCustomer)
	// g.POST("", h.HandleCreateCustomer)
	// g.PUT("/:id", h.HandleUpdateCustomer)
	// g.DELETE("/:id", h.HandleDeleteCustomer)
}
