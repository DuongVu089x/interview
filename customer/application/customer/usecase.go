package customer

import (
	"context"
	"fmt"

	domaincustomer "github.com/DuongVu089x/interview/customer/domain/customer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type UseCase struct {
	repo   domaincustomer.Repository
	tracer trace.Tracer
}

func NewUseCase(repo domaincustomer.Repository, tracer trace.Tracer) *UseCase {
	return &UseCase{
		repo:   repo,
		tracer: tracer,
	}
}

func (u *UseCase) GetCustomer(ctx context.Context, userId string) (*domaincustomer.Customer, error) {
	// The context already has the trace from gRPC interceptor
	// Just create one span for the database operation
	ctx, span := u.tracer.Start(ctx, "customer-service: customer.CustomerService")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", userId),
		attribute.String("operation", "GetCustomer"),
	)

	customer, err := u.repo.GetCustomer(ctx, userId)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	// Set customer details in span
	span.SetAttributes(
		attribute.String("customer_id", customer.ID.Hex()),
		attribute.String("customer_name", customer.Name),
		attribute.String("customer_email", customer.Email),
	)

	return customer, nil
}

func (u *UseCase) CreateCustomer(ctx context.Context, customer *domaincustomer.Customer) (*domaincustomer.Customer, error) {
	customer, err := u.repo.CreateCustomer(ctx, customer)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}
	return customer, nil
}

func (u *UseCase) UpdateCustomer(ctx context.Context, customer *domaincustomer.Customer) error {
	if err := u.repo.UpdateCustomer(ctx, customer); err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}
	return nil
}

func (u *UseCase) DeleteCustomer(ctx context.Context, id string) error {
	if err := u.repo.DeleteCustomer(ctx, id); err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}
	return nil
}
