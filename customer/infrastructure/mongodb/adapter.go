package mongodb

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// MongoAdapter provides MongoDB operations with generic support
type MongoAdapter[T any] struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewMongoAdapter creates a new MongoDB adapter using an existing client
func NewMongoAdapter[T any](client *mongo.Client, dbName string) *MongoAdapter[T] {
	return &MongoAdapter[T]{
		client: client,
		db:     client.Database(dbName),
	}
}

// Ensure MongoAdapter implements DatabasePort
// Note: We don't need to check implementation at compile time since we use generics
// The actual implementation will be checked when concrete types are used

// Query implements the DatabasePort interface
func (m *MongoAdapter[T]) Query(ctx context.Context, collection string, filter any, results any, opts ...*options.FindOptions) error {
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
func (m *MongoAdapter[T]) QueryOne(ctx context.Context, collection string, filter any, result any, opts ...*options.FindOneOptions) error {
	err := m.db.Collection(collection).FindOne(ctx, filter, opts...).Decode(result)
	if err != nil {
		return fmt.Errorf("failed to execute find one query: %w", err)
	}
	return nil
}

// FindOneAndUpdate executes a findOneAndUpdate operation and decodes the result
func (m *MongoAdapter[T]) FindOneAndUpdate(ctx context.Context, collection string, filter any, update any, result any, opts ...*options.FindOneAndUpdateOptions) error {
	err := m.db.Collection(collection).FindOneAndUpdate(ctx, filter, update, opts...).Decode(result)
	if err != nil {
		return fmt.Errorf("failed to execute find one and update: %w", err)
	}
	return nil
}

// newList return new object with same type of TemplateObject
func (m *MongoAdapter[T]) newList(limit int) interface{} {
	t := reflect.TypeOf(new(T)).Elem()
	return reflect.MakeSlice(reflect.SliceOf(t), 0, limit).Interface()
}

// convertToBson converts an entity to a BSON document
func (m *MongoAdapter[T]) convertToBson(entity any) (bson.M, error) {
	// First marshal the entity to BSON
	data, err := bson.Marshal(entity)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal entity: %w", err)
	}

	// Then unmarshal to bson.M
	var doc bson.M
	err = bson.Unmarshal(data, &doc)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal to bson.M: %w", err)
	}

	// Handle created_time if not set
	if doc["created_time"] == nil {
		doc["created_time"] = time.Now()
	}

	return doc, nil
}

// convertToObject converts a BSON document back to an entity
func (m *MongoAdapter[T]) convertToObject(doc bson.M) (T, error) {
	var result T
	data, err := bson.Marshal(doc)
	if err != nil {
		return result, fmt.Errorf("failed to marshal document: %w", err)
	}

	err = bson.Unmarshal(data, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal to entity: %w", err)
	}

	return result, nil
}

// Insert implements the DatabasePort interface
func (m *MongoAdapter[T]) Insert(ctx context.Context, collection string, documents ...any) ([]T, error) {
	if len(documents) == 1 {
		// Handle single document insert
		doc := documents[0]
		bsonDoc, err := m.convertToBson(doc)
		if err != nil {
			return nil, fmt.Errorf("failed to convert document to BSON: %w", err)
		}

		result, err := m.db.Collection(collection).InsertOne(ctx, bsonDoc)
		if err != nil {
			return nil, fmt.Errorf("failed to insert document: %w", err)
		}

		// Set the ID back in the document
		bsonDoc["_id"] = result.InsertedID
		entity, err := m.convertToObject(bsonDoc)
		if err != nil {
			return nil, fmt.Errorf("failed to convert BSON back to entity: %w", err)
		}

		return []T{entity}, nil
	}

	// Handle multiple documents
	bsonDocs := make([]any, len(documents))
	for i, doc := range documents {
		bsonDoc, err := m.convertToBson(doc)
		if err != nil {
			return nil, fmt.Errorf("failed to convert document at index %d to BSON: %w", i, err)
		}
		bsonDocs[i] = bsonDoc
	}

	result, err := m.db.Collection(collection).InsertMany(ctx, bsonDocs)
	if err != nil {
		return nil, fmt.Errorf("failed to insert documents: %w", err)
	}

	// Convert all documents back with their IDs
	var entities []T
	for i, insertedID := range result.InsertedIDs {
		bsonDoc := bsonDocs[i].(bson.M)
		bsonDoc["_id"] = insertedID
		entity, err := m.convertToObject(bsonDoc)
		if err != nil {
			return nil, fmt.Errorf("failed to convert BSON back to entity at index %d: %w", i, err)
		}
		entities = append(entities, entity)
	}

	return entities, nil
}

func (m *MongoAdapter[T]) Update(ctx context.Context, collection string, filter any, update any, opts ...*options.UpdateOptions) error {
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

func (m *MongoAdapter[T]) Delete(ctx context.Context, collection string, filter any, opts ...*options.DeleteOptions) error {
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

func (m *MongoAdapter[T]) Upsert(ctx context.Context, collection string, filter any, update any) error {
	opts := options.Update().SetUpsert(true)
	_, err := m.db.Collection(collection).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to upsert document: %w", err)
	}
	return nil
}

func (m *MongoAdapter[T]) Incr(ctx context.Context, collection string, filter any, field string, amount int64) error {
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
