package customer

import (
	"context"
	"errors"

	customerusecase "github.com/DuongVu089x/interview/customer/application/customer"
	"github.com/DuongVu089x/interview/customer/component/appctx"
	customerrepository "github.com/DuongVu089x/interview/customer/repository/customer"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrInvalidID        = errors.New("invalid customer ID")
	ErrEmptyID          = errors.New("customer ID is required")
)

type Handler struct {
	appCtx          appctx.AppContext
	customerUseCase customerusecase.UseCase
}

func NewHandler(appCtx appctx.AppContext) *Handler {
	customerRepo := customerrepository.NewMongoRepository(
		appCtx.GetMainDBConnection(),
		appCtx.GetReadMainDBConnection(),
	)
	customerUseCase := customerusecase.NewUseCase(customerRepo, appCtx.GetTracer())

	return &Handler{
		appCtx:          appCtx,
		customerUseCase: *customerUseCase,
	}
}

// GetCustomer handles retrieving a customer by ID
func (h *Handler) GetCustomer(ctx context.Context, userId string) (*CustomerResponse, error) {
	if userId == "" {
		return nil, ErrEmptyID
	}

	customer, err := h.customerUseCase.GetCustomer(ctx, userId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrCustomerNotFound
		}
		return nil, err
	}

	return ToCustomerResponse(customer), nil
}

// // CreateCustomer handles customer creation
// func (h *Handler) CreateCustomer(ctx context.Context, req *CreateCustomerRequest) (*CustomerResponse, error) {
// 	if err := h.validateCreateRequest(req); err != nil {
// 		return nil, err
// 	}

// 	customer, err := h.customerUseCase.CreateCustomer(ctx, req.ToDomainCustomer())
// 	if err != nil {
// 		return nil, err
// 	}

// 	return ToCustomerResponse(customer), nil
// }

// // UpdateCustomer handles customer updates
// func (h *Handler) UpdateCustomer(ctx context.Context, id string, req *UpdateCustomerRequest) (*CustomerResponse, error) {
// 	if err := h.validateUpdateRequest(req); err != nil {
// 		return nil, err
// 	}

// 	customer := req.ToDomainCustomer(id)
// 	if err := h.customerUseCase.UpdateCustomer(ctx, customer); err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return nil, ErrCustomerNotFound
// 		}
// 		return nil, err
// 	}

// 	return ToCustomerResponse(customer), nil
// }

// // DeleteCustomer handles customer deletion
// func (h *Handler) DeleteCustomer(ctx context.Context, id string) error {
// 	if id == "" {
// 		return ErrEmptyID
// 	}

// 	err := h.customerUseCase.DeleteCustomer(ctx, id)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return ErrCustomerNotFound
// 		}
// 		return err
// 	}

// 	return nil
// }

// // PingDB checks if the database is accessible
// func (h *Handler) PingDB(ctx context.Context) error {
// 	return h.appCtx.GetMainDBConnection().Ping(ctx, nil)
// }

// // validateCreateRequest performs validation on create customer request
// func (h *Handler) validateCreateRequest(req *CreateCustomerRequest) error {
// 	if req == nil {
// 		return errors.New("request data is required")
// 	}
// 	if req.Name == "" {
// 		return errors.New("customer name is required")
// 	}
// 	if req.Email == "" {
// 		return errors.New("customer email is required")
// 	}
// 	return nil
// }

// // validateUpdateRequest performs validation on update customer request
// func (h *Handler) validateUpdateRequest(req *UpdateCustomerRequest) error {
// 	if req == nil {
// 		return errors.New("request data is required")
// 	}
// 	if req.Name == "" {
// 		return errors.New("customer name is required")
// 	}
// 	if req.Email == "" {
// 		return errors.New("customer email is required")
// 	}
// 	return nil
// }
