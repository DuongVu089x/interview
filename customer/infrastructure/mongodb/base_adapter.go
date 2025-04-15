package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// BaseAdapter provides common functionality for MongoDB adapters
type BaseAdapter[T any] struct {
	writeDB *MongoAdapter[T]
	readDB  *MongoAdapter[T]
}

// NewBaseAdapter creates a new base adapter using existing clients
func NewBaseAdapter[T any](writeClient, readClient *mongo.Client, dbName string) *BaseAdapter[T] {
	return &BaseAdapter[T]{
		writeDB: NewMongoAdapter[T](writeClient, dbName),
		readDB:  NewMongoAdapter[T](readClient, dbName),
	}
}

// GetWriteDB returns the write adapter
func (b *BaseAdapter[T]) GetWriteDB() *MongoAdapter[T] {
	return b.writeDB
}

// GetReadDB returns the read adapter
func (b *BaseAdapter[T]) GetReadDB() *MongoAdapter[T] {
	return b.readDB
}
