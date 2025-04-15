package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// BaseAdapter provides common functionality for MongoDB adapters
type BaseAdapter struct {
	writeDB *MongoAdapter
	readDB  *MongoAdapter
}

// NewBaseAdapter creates a new base adapter using existing clients
func NewBaseAdapter(writeClient, readClient *mongo.Client, dbName string) *BaseAdapter {
	return &BaseAdapter{
		writeDB: NewMongoAdapter(writeClient, dbName),
		readDB:  NewMongoAdapter(readClient, dbName),
	}
}

// GetWriteDB returns the write adapter
func (b *BaseAdapter) GetWriteDB() *MongoAdapter {
	return b.writeDB
}

// GetReadDB returns the read adapter
func (b *BaseAdapter) GetReadDB() *MongoAdapter {
	return b.readDB
}
