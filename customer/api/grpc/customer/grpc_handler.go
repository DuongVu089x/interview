package customer

import (
	"context"
	"log"
	"time"

	customerhandler "github.com/DuongVu089x/interview/customer/api/handler/customer"
	"github.com/DuongVu089x/interview/customer/component/appctx"
	pb "github.com/DuongVu089x/interview/customer/proto/customer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcHandler struct {
	pb.UnimplementedCustomerServiceServer
	handler *customerhandler.Handler
}

func NewGrpcHandler(appCtx appctx.AppContext) *GrpcHandler {
	return &GrpcHandler{
		handler: customerhandler.NewHandler(appCtx),
	}
}

func (h *GrpcHandler) GetCustomer(ctx context.Context, req *pb.GetCustomerRequest) (*pb.GetCustomerResponse, error) {
	customer, err := h.handler.GetCustomer(ctx, req.UserId)
	if err != nil {
		switch err {
		case customerhandler.ErrCustomerNotFound:
			return &pb.GetCustomerResponse{
				Customer: nil,
				Exists:   false,
			}, nil
		case customerhandler.ErrEmptyID:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			log.Printf("Error getting customer: %v", err)
			return nil, status.Errorf(codes.Internal, "Failed to get customer: %v", err)
		}
	}

	return &pb.GetCustomerResponse{
		Customer: &pb.Customer{
			Id:        customer.ID,
			Name:      customer.Name,
			Email:     customer.Email,
			Phone:     customer.Phone,
			CreatedAt: customer.CreatedAt.Format(time.RFC3339),
			UpdatedAt: customer.UpdatedAt.Format(time.RFC3339),
		},
		Exists: true,
	}, nil
}
