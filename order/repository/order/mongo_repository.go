package order

import (
	"context"

	domainorder "github.com/DuongVu089x/interview/order/domain/order"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Client) domainorder.Repository {
	return &MongoRepository{
		collection: db.Database("orders").Collection("orders"),
	}
}

func (r *MongoRepository) GetOrderByCustomerID(customerID string, conditions map[string]interface{}) ([]domainorder.Order, error) {
	var orders []domainorder.Order

	query := bson.M{"userId": customerID}
	if status, ok := conditions["status"]; ok {
		query["status"] = status
	}

	cursor, err := r.collection.Find(context.Background(), query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *MongoRepository) CreateOrder(order *domainorder.Order) error {
	_, err := r.collection.InsertOne(context.Background(), order)
	return err
}

func (r *MongoRepository) GetOrder(id string) (*domainorder.Order, error) {
	var order domainorder.Order
	err := r.collection.FindOne(context.Background(), bson.M{"order_id": id}).Decode(&order)
	return &order, err
}

func (r *MongoRepository) UpdateOrder(order *domainorder.Order) error {
	_, err := r.collection.UpdateOne(context.Background(), bson.M{"order_id": order.OrderID}, bson.M{"$set": order})
	return err
}

func (r *MongoRepository) DeleteOrder(id string) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"order_id": id})
	return err
}
