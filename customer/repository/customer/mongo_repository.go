package customer

import (
	"context"

	domaincustomer "github.com/DuongVu089x/interview/customer/domain/customer"
	"github.com/DuongVu089x/interview/customer/infrastructure/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	*mongodb.BaseAdapter[*domaincustomer.Customer]
}

const (
	databaseName   = "customers"
	collectionName = "customers"
)

func NewMongoRepository(writeDB, readDB *mongo.Client) domaincustomer.Repository {
	baseAdapter := mongodb.NewBaseAdapter[*domaincustomer.Customer](writeDB, readDB, databaseName)
	return &MongoRepository{
		BaseAdapter: baseAdapter,
	}
}

func (r *MongoRepository) GetCustomer(ctx context.Context, userId string) (*domaincustomer.Customer, error) {
	var customer domaincustomer.Customer
	err := r.GetReadDB().QueryOne(
		ctx,
		collectionName,
		domaincustomer.Customer{
			UserId: userId,
		},
		&customer,
	)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *MongoRepository) CreateCustomer(ctx context.Context, customer *domaincustomer.Customer) (*domaincustomer.Customer, error) {
	customers, err := r.GetWriteDB().Insert(ctx, collectionName, customer)
	if err != nil {
		return nil, err
	}
	return customers[0], nil
}

func (r *MongoRepository) UpdateCustomer(ctx context.Context, customer *domaincustomer.Customer) error {
	filter := bson.M{"_id": customer.ID}
	update := bson.M{"$set": customer}
	return r.GetWriteDB().Update(ctx, collectionName, filter, update)
}

func (r *MongoRepository) DeleteCustomer(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	return r.GetWriteDB().Delete(ctx, collectionName, filter)
}

// GetCustomers retrieves multiple customers based on a filter
func (r *MongoRepository) GetCustomers(ctx context.Context, filter any) ([]*domaincustomer.Customer, error) {
	var customers []*domaincustomer.Customer
	err := r.GetReadDB().Query(
		ctx,
		collectionName,
		filter,
		&customers,
	)
	if err != nil {
		return nil, err
	}
	return customers, nil
}
