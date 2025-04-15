package order

import (
	"context"

	domainorder "github.com/DuongVu089x/interview/order/domain/order"
	"github.com/DuongVu089x/interview/order/infrastructure/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	*mongodb.BaseAdapter
}

const (
	databaseName   = "orders"
	collectionName = "orders"
)

func NewMongoRepository(writeDB, readDB *mongo.Client) domainorder.Repository {
	baseAdapter := mongodb.NewBaseAdapter(writeDB, readDB, databaseName)
	return &MongoRepository{
		BaseAdapter: baseAdapter,
	}
}

func (r *MongoRepository) GetOrder(id int64) (*domainorder.Order, error) {
	var order domainorder.Order
	err := r.GetReadDB().QueryOne(
		context.Background(),
		collectionName,
		bson.M{"order_id": id},
		&order,
	)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *MongoRepository) GetOrders(conditions domainorder.Order) ([]domainorder.Order, error) {
	filter := bson.M{}
	if conditions.OrderID != 0 {
		filter["order_id"] = conditions.OrderID
	}
	if conditions.UserID != "" {
		filter["user_id"] = conditions.UserID
	}
	if conditions.Status != "" {
		filter["status"] = conditions.Status
	}

	var orders []domainorder.Order
	err := r.GetReadDB().Query(
		context.Background(),
		collectionName,
		filter,
		&orders,
		options.Find().SetSort(bson.M{"created_at": -1}),
	)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *MongoRepository) CreateOrder(order *domainorder.Order) error {
	return r.GetWriteDB().Insert(context.Background(), collectionName, order)
}

func (r *MongoRepository) UpdateOrder(order *domainorder.Order) error {
	filter := bson.M{"order_id": order.OrderID}
	update := bson.M{"$set": order}
	return r.GetWriteDB().Update(context.Background(), collectionName, filter, update)
}

func (r *MongoRepository) DeleteOrder(id string) error {
	filter := bson.M{"order_id": id}
	return r.GetWriteDB().Delete(context.Background(), collectionName, filter)
}
