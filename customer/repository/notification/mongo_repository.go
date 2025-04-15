package notification

import (
	"context"

	domainnotification "github.com/DuongVu089x/interview/customer/domain/notification"
	"github.com/DuongVu089x/interview/customer/infrastructure/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	*mongodb.BaseAdapter[*domainnotification.Notification]
}

const (
	databaseName   = "customers"
	collectionName = "notifications"
)

func NewMongoRepository(writeDB, readDB *mongo.Client) domainnotification.Repository {
	baseAdapter := mongodb.NewBaseAdapter[*domainnotification.Notification](writeDB, readDB, databaseName)
	return &MongoRepository{
		BaseAdapter: baseAdapter,
	}
}

func (r *MongoRepository) CreateNotification(ctx context.Context, notification *domainnotification.Notification) (*domainnotification.Notification, error) {
	result, err := r.GetWriteDB().Insert(ctx, collectionName, notification)
	if err != nil {
		return nil, err
	}
	return result[0], nil
}

func (r *MongoRepository) GetNotification(ctx context.Context, id string) (*domainnotification.Notification, error) {
	var notification domainnotification.Notification
	err := r.GetReadDB().QueryOne(ctx, collectionName, bson.M{"_id": id}, &notification)
	return &notification, err
}

func (r *MongoRepository) MarkAsReadNotification(ctx context.Context, id string) error {
	return r.GetWriteDB().Update(ctx, collectionName, bson.M{"_id": id}, bson.M{"$set": bson.M{"isRead": true}})
}

func (r *MongoRepository) GetNotifications(ctx context.Context, userId string, offset, limit int64) ([]*domainnotification.Notification, error) {
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}})

	if offset > 0 {
		opts.SetSkip(offset)
	}
	if limit > 0 {
		opts.SetLimit(limit)
	}

	var notifications []*domainnotification.Notification
	err := r.GetReadDB().Query(ctx, collectionName,
		domainnotification.Notification{
			UserID: userId,
		},
		&notifications,
		opts,
	)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}
