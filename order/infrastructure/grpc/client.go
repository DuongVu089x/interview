package grpc

import (
	"context"
	"fmt"

	pb "github.com/DuongVu089x/interview/order/proto/customer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CustomerClient struct {
	client pb.CustomerServiceClient
}

func NewCustomerClient(host string, port string) (*CustomerClient, error) {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to customer service: %v", err)
	}

	return &CustomerClient{
		client: pb.NewCustomerServiceClient(conn),
	}, nil
}

func (c *CustomerClient) CheckCustomerExists(ctx context.Context, userID string) (bool, error) {
	resp, err := c.client.GetCustomer(ctx, &pb.GetCustomerRequest{
		UserId: userID,
	})
	if err != nil {
		return false, fmt.Errorf("failed to check customer: %v", err)
	}

	return resp.Exists, nil
}
