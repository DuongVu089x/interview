package userconnection

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	domainuserconnection "github.com/DuongVu089x/interview/customer/domain/user_connection"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Client) domainuserconnection.Repository {
	return &MongoRepository{collection: db.Database("customer").Collection("user_connections")}
}

func (r *MongoRepository) GetUserConnection(query *domainuserconnection.UserConnection) (*domainuserconnection.UserConnection, error) {
	userConn := &domainuserconnection.UserConnection{}
	err := r.collection.FindOne(context.Background(), query).Decode(userConn)
	if err != nil {
		return nil, err
	}
	return userConn, nil
}

func (r *MongoRepository) GetUserConnections(query *domainuserconnection.UserConnection, offset, limit int64) ([]*domainuserconnection.UserConnection, error) {
	opts := &options.FindOptions{}
	if offset > 0 {
		opts.SetSkip(offset)
	}
	if limit > 0 {
		opts.SetLimit(limit)
	}

	cursor, err := r.collection.Find(context.Background(), query, opts)
	if err != nil {
		return nil, err
	}

	userConns := make([]*domainuserconnection.UserConnection, 0)
	err = cursor.All(context.Background(), &userConns)
	if err != nil {
		return nil, err
	}
	return userConns, nil
}

func (r *MongoRepository) CreateUserConnection(userConn *domainuserconnection.UserConnection) (*domainuserconnection.UserConnection, error) {
	result, err := r.collection.InsertOne(context.Background(), userConn)
	if err != nil {
		return nil, err
	}

	userConn.ID = result.InsertedID.(primitive.ObjectID)
	return userConn, nil
}

func (r *MongoRepository) UpdateUserConnection(query *domainuserconnection.UserConnection, updating *domainuserconnection.UserConnection) error {
	_, err := r.collection.UpdateOne(context.Background(), query, bson.M{"$set": updating})
	return err
}

func (r *MongoRepository) DeleteUserConnection(userConn *domainuserconnection.UserConnection) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": userConn.ID})
	return err
}
