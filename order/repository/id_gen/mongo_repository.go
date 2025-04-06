package id_gen

import (
	"context"

	domainidgen "github.com/DuongVu089x/interview/order/domain/id_gen"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Client) domainidgen.Repository {
	return &MongoRepository{
		collection: db.Database("orders").Collection("id_gen"),
	}
}

func (r *MongoRepository) GenerateID(key string) (int64, error) {
	result := r.collection.FindOneAndUpdate(context.Background(),
		bson.M{"key": key},
		bson.M{
			"$inc": bson.M{"value": 1},
		},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After))

	var idGen domainidgen.IDGen
	err := result.Decode(&idGen)
	if err != nil {
		return 0, err
	}

	return idGen.Value, nil
}
