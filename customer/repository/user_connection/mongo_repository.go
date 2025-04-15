package userconnection

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	domainuserconnection "github.com/DuongVu089x/interview/customer/domain/user_connection"
	"github.com/DuongVu089x/interview/customer/infrastructure/mongodb"
)

type MongoRepository struct {
	*mongodb.BaseAdapter[*domainuserconnection.UserConnection]
}

const (
	databaseName   = "customers"
	collectionName = "user_connections"
)

func NewMongoRepository(writeDB, readDB *mongo.Client) domainuserconnection.Repository {
	baseAdapter := mongodb.NewBaseAdapter[*domainuserconnection.UserConnection](writeDB, readDB, databaseName)

	return &MongoRepository{
		BaseAdapter: baseAdapter,
	}
}

func (r *MongoRepository) GetUserConnection(ctx context.Context, query *domainuserconnection.UserConnection) (*domainuserconnection.UserConnection, error) {
	userConn := &domainuserconnection.UserConnection{}
	err := r.GetReadDB().QueryOne(ctx, collectionName, query, userConn)
	if err != nil {
		return nil, err
	}
	return userConn, nil
}

func (r *MongoRepository) GetUserConnections(ctx context.Context, query *domainuserconnection.UserConnection, offset, limit int64) ([]*domainuserconnection.UserConnection, error) {
	opts := &options.FindOptions{}
	if offset > 0 {
		opts.SetSkip(offset)
	}
	if limit > 0 {
		opts.SetLimit(limit)
	}

	var userConns []*domainuserconnection.UserConnection
	err := r.GetReadDB().Query(ctx, collectionName, query, &userConns, opts)
	if err != nil {
		return nil, err
	}

	return userConns, nil
}

func (r *MongoRepository) CreateUserConnection(ctx context.Context, userConn *domainuserconnection.UserConnection) (*domainuserconnection.UserConnection, error) {
	result, err := r.GetWriteDB().Insert(ctx, collectionName, userConn)
	if err != nil {
		return nil, err
	}

	return result[0], nil
}

func (r *MongoRepository) UpdateUserConnection(ctx context.Context, query *domainuserconnection.UserConnection, updating *domainuserconnection.UserConnection) error {
	return r.GetWriteDB().Update(ctx, collectionName, query, bson.M{"$set": updating})
}

func (r *MongoRepository) DeleteUserConnection(ctx context.Context, userConn *domainuserconnection.UserConnection) error {
	return r.GetWriteDB().Delete(ctx, collectionName, userConn)
}
