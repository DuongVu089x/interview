package notification

import (
	"context"

	domainnotification "github.com/DuongVu089x/interview/customer/domain/notification"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Client) domainnotification.Repository {
	return &MongoRepository{collection: db.Database("customer").Collection("notifications")}
}

func (r *MongoRepository) CreateNotification(notification *domainnotification.Notification) error {
	_, err := r.collection.InsertOne(context.Background(), notification)
	return err
}

func (r *MongoRepository) GetNotification(id string) (*domainnotification.Notification, error) {
	var notification domainnotification.Notification
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&notification)
	return &notification, err
}

func (r *MongoRepository) MarkAsReadNotification(id string) error {
	_, err := r.collection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": bson.M{"isRead": true}})
	return err
}
