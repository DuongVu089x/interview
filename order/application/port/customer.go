package port

import "context"

// CustomerPort defines the interface for customer service operations
type CustomerPort interface {
	// CheckCustomerExists verifies if a customer exists by their user ID
	CheckCustomerExists(ctx context.Context, userID string) (bool, error)
}
