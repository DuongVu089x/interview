package customer

import (
	"context"
	"fmt"

	domaincustomer "github.com/DuongVu089x/interview/customer/domain/customer"
)

type UseCase struct {
	repo domaincustomer.Repository
}

func NewUseCase(repo domaincustomer.Repository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

func (u *UseCase) GetCustomer(ctx context.Context, userId string) (*domaincustomer.Customer, error) {
	customer, err := u.repo.GetCustomer(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
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
