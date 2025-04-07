package order

import (
	"fmt"
	"net/http"
	"strconv"

	orderusecase "github.com/DuongVu089x/interview/order/application/order"
	"github.com/DuongVu089x/interview/order/component/appctx"
	idgenrepository "github.com/DuongVu089x/interview/order/repository/id_gen"
	orderrepository "github.com/DuongVu089x/interview/order/repository/order"
	idgenservice "github.com/DuongVu089x/interview/order/service/id_gen"
	orderservice "github.com/DuongVu089x/interview/order/service/order"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	appCtx       appctx.AppContext
	orderUseCase *orderusecase.UseCase
	validator    *CustomValidator
}

func NewHandler(appCtx appctx.AppContext) *Handler {
	orderRepo := orderrepository.NewMongoRepository(appCtx.GetMainDBConnection())
	orderService := orderservice.NewOrderService(orderRepo)

	idgenRepo := idgenrepository.NewMongoRepository(appCtx.GetMainDBConnection())
	idgenService := idgenservice.NewIDGenService(idgenRepo)
	orderUseCase := orderusecase.NewOrderUseCase(orderService, idgenService)
	return &Handler{
		appCtx:       appCtx,
		orderUseCase: orderUseCase,
		validator:    NewCustomValidator(),
	}
}

// CreateOrder handles order creation requests
func (h *Handler) CreateOrder(c echo.Context) error {
	var req orderusecase.CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	if err := h.validator.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	response, err := h.orderUseCase.CreateOrder(h.appCtx, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create order "+err.Error())
	}

	return c.JSON(http.StatusCreated, response)
}

// GetOrder handles single order retrieval
func (h *Handler) GetOrder(c echo.Context) error {
	id := c.Param("id")

	orderID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid order ID")
	}

	order, err := h.orderUseCase.GetOrder(orderID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get order")
	}
	return c.JSON(http.StatusOK, order)
}

// GetOrdersByUserID handles retrieval of all orders for a specific user
func (h *Handler) GetOrdersByUserID(c echo.Context) error {
	// Extract user ID from path parameter
	userID := c.Param("userId")
	fmt.Printf("GetOrdersByUserID handler called with userID: %s\n", userID)

	if userID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID is required")
	}

	// Extract optional status from query parameter
	status := c.QueryParam("status")
	fmt.Printf("GetOrdersByUserID status filter: %s\n", status)

	// Create the request
	req := orderusecase.GetOrdersByUserIDRequest{
		UserID: userID,
		Status: status,
	}

	// Call the use case
	response, err := h.orderUseCase.GetOrdersByUserID(req)
	if err != nil {
		fmt.Printf("GetOrdersByUserID error: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get orders: "+err.Error())
	}

	fmt.Printf("GetOrdersByUserID success, returning %d orders\n", response.Count)
	return c.JSON(http.StatusOK, response)
}
