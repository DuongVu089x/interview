package port

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// DatabasePort defines the interface for database operations
type DatabasePort interface {
	// Query returns multiple documents matching the filter
	Query(ctx context.Context, collection string, filter any, results any, opts ...*options.FindOptions) error

	// QueryOne returns a single document matching the filter
	QueryOne(ctx context.Context, collection string, filter any, result any, opts ...*options.FindOneOptions) error

	// Insert inserts one or more documents
	Insert(ctx context.Context, collection string, documents ...any) error

	// Update updates documents matching the filter
	Update(ctx context.Context, collection string, filter any, update any, opts ...*options.UpdateOptions) error

	// Delete deletes documents matching the filter
	Delete(ctx context.Context, collection string, filter any, opts ...*options.DeleteOptions) error

	// Upsert inserts a document if it doesn't exist, or updates it if it does
	Upsert(ctx context.Context, collection string, filter any, update any) error

	// Incr atomically increments a numeric field by the specified amount
	Incr(ctx context.Context, collection string, filter any, field string, amount int64) error
}
