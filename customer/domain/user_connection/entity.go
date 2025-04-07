package userconnection

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserConnection represent
type UserConnection struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	LastUpdatedTime *time.Time         `json:"lastUpdatedTime,omitempty" bson:"last_updated_time,omitempty"`
	CreatedTime     *time.Time         `json:"createdTime,omitempty" bson:"created_time,omitempty"`

	WSHost        string `json:"wsHost,omitempty" bson:"ws_host,omitempty"`
	WSHostVersion string `json:"wsHostVersion,omitempty" bson:"ws_host_version,omitempty"`
	IP            string `json:"ip,omitempty" bson:"ip,omitempty"`
	UserID        string `json:"userId,omitempty" bson:"user_id,omitempty"`
	UserAgent     string `json:"userAgent,omitempty" bson:"user_agent,omitempty"`

	ConnectedTime   *time.Time `json:"connectedTime,omitempty" bson:"connected_time,omitempty"`
	LastMessageTime *time.Time `json:"lastMessageTime,omitempty" bson:"last_message_time,omitempty"`

	DisconnectedTime   *time.Time `json:"disconnectedTime,omitempty" bson:"disconnected_time,omitempty"`
	DisconnectedReason string     `json:"disconnectedReason,omitempty" bson:"disconnected_reason,omitempty"`

	Status ConStatusEnumValue `json:"status,omitempty" bson:"status,omitempty"`

	ConnectionLocalID int `json:"connectionLocalId,omitempty" bson:"connection_local_id,omitempty"` // for websocket local connection id

	ComplexQuery []bson.M `json:"-" bson:"$and,omitempty"`
}

type ConStatusEnumValue string

type ConStatusEnum struct {
	ACTIVE ConStatusEnumValue
	CLOSED ConStatusEnumValue
}

var ConStatus = &ConStatusEnum{
	ACTIVE: ConStatusEnumValue("ACTIVE"),
	CLOSED: ConStatusEnumValue("CLOSED"),
}
