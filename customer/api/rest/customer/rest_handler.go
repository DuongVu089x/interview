package customer

import (
	"net/http"

	customerhandler "github.com/DuongVu089x/interview/customer/api/handler/customer"
	"github.com/DuongVu089x/interview/customer/component/appctx"
	"github.com/labstack/echo/v4"
)

type RestHandler struct {
	handler *customerhandler.Handler
}

func NewRestHandler(appCtx appctx.AppContext) *RestHandler {
	return &RestHandler{
		handler: customerhandler.NewHandler(appCtx),
	}
}

func (h *RestHandler) HandleGetCustomer(c echo.Context) error {
	id := c.Param("id")

	customer, err := h.handler.GetCustomer(c.Request().Context(), id)
	if err != nil {
		switch err {
		case customerhandler.ErrCustomerNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		case customerhandler.ErrEmptyID:
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get customer")
		}
	}

	return c.JSON(http.StatusOK, customer)
}

// func (h *RestHandler) HandleCreateCustomer(c echo.Context) error {
// 	var req customerhandler.CreateCustomerRequest
// 	if err := c.Bind(&req); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
// 	}

// 	customer, err := h.handler.CreateCustomer(c.Request().Context(), &req)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
// 	}

// 	return c.JSON(http.StatusCreated, customer)
// }

// func (h *RestHandler) HandleUpdateCustomer(c echo.Context) error {
// 	id := c.Param("id")

// 	var req customerhandler.UpdateCustomerRequest
// 	if err := c.Bind(&req); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
// 	}

// 	customer, err := h.handler.UpdateCustomer(c.Request().Context(), id, &req)
// 	if err != nil {
// 		switch err {
// 		case customerhandler.ErrCustomerNotFound:
// 			return echo.NewHTTPError(http.StatusNotFound, err.Error())
// 		case customerhandler.ErrInvalidID:
// 			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
// 		default:
// 			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update customer")
// 		}
// 	}

// 	return c.JSON(http.StatusOK, customer)
// }

// func (h *RestHandler) HandleDeleteCustomer(c echo.Context) error {
// 	id := c.Param("id")

// 	err := h.handler.DeleteCustomer(c.Request().Context(), id)
// 	if err != nil {
// 		switch err {
// 		case customerhandler.ErrCustomerNotFound:
// 			return echo.NewHTTPError(http.StatusNotFound, err.Error())
// 		case customerhandler.ErrEmptyID:
// 			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
// 		default:
// 			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete customer")
// 		}
// 	}

// 	return c.NoContent(http.StatusNoContent)
// }
