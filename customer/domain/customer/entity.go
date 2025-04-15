package customer

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID        *primitive.ObjectID `bson:"_id,omitempty"`
	UserId    string              `bson:"user_id,omitempty"`
	Name      string              `bson:"name,omitempty"`
	Email     string              `bson:"email,omitempty"`
	Phone     string              `bson:"phone,omitempty"`
	CreatedAt time.Time           `bson:"created_at,omitempty"`
	UpdatedAt time.Time           `bson:"updated_at,omitempty"`
}
