package id_gen

import (
	"context"

	domainidgen "github.com/DuongVu089x/interview/order/domain/id_gen"
	"github.com/DuongVu089x/interview/order/infrastructure/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	writeDB *mongodb.MongoAdapter
}

const (
	databaseName   = "orders"
	collectionName = "id_gen"
)

func NewMongoRepository(writeDB *mongo.Client) domainidgen.Repository {
	writeAdapter := mongodb.NewMongoAdapter(writeDB, databaseName)

	return &MongoRepository{
		writeDB: writeAdapter,
	}
}

func (r *MongoRepository) GenerateID(key string) (int64, error) {
	var idGen domainidgen.IDGen

	err := r.writeDB.FindOneAndUpdate(
		context.Background(),
		collectionName,
		bson.M{"key": key},
		bson.M{"$inc": bson.M{"value": 1}},
		&idGen,
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	)
	if err != nil {
		return 0, err
	}

	return idGen.Value, nil
}
