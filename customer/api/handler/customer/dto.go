package customer

import (
	"time"

	domaincustomer "github.com/DuongVu089x/interview/customer/domain/customer"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CustomerResponse represents the customer data returned to clients
type CustomerResponse struct {
	ID        string    `json:"id,omitempty"`
	UserId    string    `json:"userId,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

// CreateCustomerRequest represents the request body for creating a customer
type CreateCustomerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone,omitempty"`
}

// UpdateCustomerRequest represents the request body for updating a customer
type UpdateCustomerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone,omitempty"`
}

// ToCustomerResponse converts a domain customer to a customer response DTO
func ToCustomerResponse(customer *domaincustomer.Customer) *CustomerResponse {
	if customer == nil {
		return nil
	}

	response := &CustomerResponse{
		UserId:    customer.UserId,
		Name:      customer.Name,
		Email:     customer.Email,
		Phone:     customer.Phone,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}

	if customer.ID != nil {
		response.ID = customer.ID.Hex()
	}

	return response
}

// ToDomainCustomer converts a create customer request to a domain customer
func (req *CreateCustomerRequest) ToDomainCustomer() *domaincustomer.Customer {
	return &domaincustomer.Customer{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}
}

// ToDomainCustomer converts an update customer request to a domain customer
func (req *UpdateCustomerRequest) ToDomainCustomer(id string) *domaincustomer.Customer {
	objectID, _ := primitive.ObjectIDFromHex(id)
	return &domaincustomer.Customer{
		ID:    &objectID,
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}
}
