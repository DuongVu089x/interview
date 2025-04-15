package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/DuongVu089x/interview/order/application/port"
)

// ErrEmptyFilter is returned when an update or delete operation is attempted with an empty filter
var ErrEmptyFilter = fmt.Errorf("filter cannot be empty for this operation")

// isEmptyFilter checks if the filter is empty or contains no conditions
func isEmptyFilter(filter any) bool {
	if filter == nil {
		return true
	}

	switch f := filter.(type) {
	case bson.M:
		return len(f) == 0
	case bson.D:
		return len(f) == 0
	default:
		// For other types, we'll consider them non-empty
		// This is a conservative approach
		return false
	}
}

type MongoAdapter struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewMongoAdapter creates a new MongoDB adapter using an existing client
func NewMongoAdapter(client *mongo.Client, dbName string) *MongoAdapter {
	return &MongoAdapter{
		client: client,
		db:     client.Database(dbName),
	}
}

// Ensure MongoAdapter implements DatabasePort
var _ port.DatabasePort = (*MongoAdapter)(nil)

// Query implements the DatabasePort interface
func (m *MongoAdapter) Query(ctx context.Context, collection string, filter any, results any, opts ...*options.FindOptions) error {
	cursor, err := m.db.Collection(collection).Find(ctx, filter, opts...)
	if err != nil {
		return fmt.Errorf("failed to execute find query: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, results); err != nil {
		return fmt.Errorf("failed to decode results: %w", err)
	}
	return nil
}

// QueryOne implements the DatabasePort interface
func (m *MongoAdapter) QueryOne(ctx context.Context, collection string, filter any, result any, opts ...*options.FindOneOptions) error {
	err := m.db.Collection(collection).FindOne(ctx, filter, opts...).Decode(result)
	if err != nil {
		return fmt.Errorf("failed to execute find one query: %w", err)
	}
	return nil
}

// FindOneAndUpdate executes a findOneAndUpdate operation and decodes the result
func (m *MongoAdapter) FindOneAndUpdate(ctx context.Context, collection string, filter any, update any, result any, opts ...*options.FindOneAndUpdateOptions) error {
	err := m.db.Collection(collection).FindOneAndUpdate(ctx, filter, update, opts...).Decode(result)
	if err != nil {
		return fmt.Errorf("failed to execute find one and update: %w", err)
	}
	return nil
}

func (m *MongoAdapter) Insert(ctx context.Context, collection string, documents ...any) error {
	_, err := m.db.Collection(collection).InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("failed to insert documents: %w", err)
	}
	return nil
}

func (m *MongoAdapter) Update(ctx context.Context, collection string, filter any, update any, opts ...*options.UpdateOptions) error {
	if isEmptyFilter(filter) {
		return fmt.Errorf("%w: update requires a non-empty filter", ErrEmptyFilter)
	}

	// Validate update document
	updateDoc, ok := update.(bson.M)
	if !ok {
		return fmt.Errorf("update document must be bson.M")
	}

	// Ensure update has operator
	hasOperator := false
	for k := range updateDoc {
		if k[0] == '$' {
			hasOperator = true
			break
		}
	}
	if !hasOperator {
		update = bson.M{"$set": update}
	}

	// By default, use UpdateOne for safety
	_, err := m.db.Collection(collection).UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}
	return nil
}

func (m *MongoAdapter) Delete(ctx context.Context, collection string, filter any, opts ...*options.DeleteOptions) error {
	if isEmptyFilter(filter) {
		return fmt.Errorf("%w: delete requires a non-empty filter", ErrEmptyFilter)
	}

	// By default, use DeleteOne for safety
	_, err := m.db.Collection(collection).DeleteOne(ctx, filter, opts...)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	return nil
}

func (m *MongoAdapter) Upsert(ctx context.Context, collection string, filter any, update any) error {
	opts := options.Update().SetUpsert(true)
	_, err := m.db.Collection(collection).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to upsert document: %w", err)
	}
	return nil
}

func (m *MongoAdapter) Incr(ctx context.Context, collection string, filter any, field string, amount int64) error {
	update := bson.M{
		"$inc": bson.M{
			field: amount,
		},
	}
	_, err := m.db.Collection(collection).UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to increment field: %w", err)
	}
	return nil
}
