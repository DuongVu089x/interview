package notification

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID              *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" `
	CreatedTime     *time.Time          `json:"createdTime,omitempty" bson:"created_time,omitempty"`
	LastUpdatedTime *time.Time          `json:"lastUpdatedTime,omitempty" bson:"last_updated_time,omitempty"`

	Code   string `json:"code,omitempty" bson:"code,omitempty" `
	UserID string `json:"userId,omitempty" bson:"user_id,omitempty" `
	IsRead *bool  `json:"isRead,omitempty" bson:"is_read,omitempty" `

	Topic       string `json:"topic,omitempty" bson:"topic,omitempty"`
	Title       string `json:"title,omitempty" bson:"title,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty" `
	Link        string `json:"link,omitempty" bson:"link,omitempty" `
}

type TopicEnumValue string
type TopicEnum struct {
	NONE          TopicEnumValue
	CONNECTED     TopicEnumValue
	CONNECTION    TopicEnumValue
	PING          TopicEnumValue
	AUTHORIZATION TopicEnumValue

	CART_UPDATE  TopicEnumValue
	ANNOUNCEMENT TopicEnumValue
	PROMOTION    TopicEnumValue
	EVENT        TopicEnumValue
}

var Topic = &TopicEnum{
	NONE:          TopicEnumValue("NONE"), // don't use. Just used for upsert & get new id.
	CONNECTED:     TopicEnumValue("CONNECTED"),
	CONNECTION:    TopicEnumValue("CONNECTION"),
	PING:          TopicEnumValue("PING"),
	AUTHORIZATION: TopicEnumValue("AUTHORIZATION"),

	// current valid topic
	CART_UPDATE:  TopicEnumValue("CART_UPDATE"),
	ANNOUNCEMENT: TopicEnumValue("ANNOUNCEMENT"),
	PROMOTION:    TopicEnumValue("PROMOTION"),
	EVENT:        TopicEnumValue("EVENT"),
}
