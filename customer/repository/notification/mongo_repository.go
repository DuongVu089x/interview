package notification

import (
	"context"

	domainnotification "github.com/DuongVu089x/interview/customer/domain/notification"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *MongoRepository) GetNotifications(userId string, offset, limit int64) ([]*domainnotification.Notification, error) {
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}})

	if offset > 0 {
		opts.SetSkip(offset)
	}
	if limit > 0 {
		opts.SetLimit(limit)
	}

	cursor, err := r.collection.Find(context.Background(), bson.M{"user_id": userId}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var notifications []*domainnotification.Notification
	if err = cursor.All(context.Background(), &notifications); err != nil {
		return nil, err
	}

	return notifications, nil
}
