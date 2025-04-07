package websocket

import (
	"encoding/json"
)

// TopicEnumValue represents the type of WebSocket message
type TopicEnumValue string

const (
	// TOPIC_NONE represents an undefined topic
	TOPIC_NONE TopicEnumValue = "NONE"
	// TOPIC_CONNECTED represents a connection established message
	TOPIC_CONNECTED TopicEnumValue = "CONNECTED"
	// TOPIC_CONNECTION represents connection status related messages
	TOPIC_CONNECTION TopicEnumValue = "CONNECTION"
	// TOPIC_PING represents ping/heartbeat messages
	TOPIC_PING TopicEnumValue = "PING"
	// TOPIC_AUTHORIZATION represents authentication related messages
	TOPIC_AUTHORIZATION TopicEnumValue = "AUTHORIZATION"
	// TOPIC_ANNOUNCEMENT represents system announcement messages
	TOPIC_ANNOUNCEMENT TopicEnumValue = "ANNOUNCEMENT"
	// TOPIC_EVENT represents general event messages
	TOPIC_EVENT TopicEnumValue = "EVENT"
)

// IsValid checks if the topic is a valid enum value
func (t TopicEnumValue) IsValid() bool {
	switch t {
	case TOPIC_NONE, TOPIC_CONNECTED, TOPIC_CONNECTION, TOPIC_PING,
		TOPIC_AUTHORIZATION, TOPIC_ANNOUNCEMENT, TOPIC_EVENT:
		return true
	}
	return false
}

// String returns the string representation of the Topic
func (t TopicEnumValue) String() string {
	return string(t)
}

// WSMessage represents the base structure for both input and output WebSocket messages
type WSMessage struct {
	Topic    TopicEnumValue `json:"topic,omitempty"`
	Content  map[string]any `json:"content,omitempty"`
	Callback string         `json:"callback,omitempty"`
}

// WSInputMessage represents a message received from the client
type WSInputMessage struct {
	Topic    TopicEnumValue `json:"topic,omitempty"`
	Content  map[string]any `json:"content,omitempty"`
	Callback string         `json:"callback,omitempty"`
}

// WSOutputMessage represents a message sent to the client
type WSOutputMessage struct {
	Topic    TopicEnumValue `json:"topic,omitempty"`
	Content  map[string]any `json:"content,omitempty"`
	Callback string         `json:"callback,omitempty"`
}

// ToString converts the output message to a JSON string
func (o *WSOutputMessage) ToString() string {
	bytes, err := json.Marshal(o)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// WSAuthorization represents the authorization payload in WebSocket messages
type WSAuthorization struct {
	SessionToken string `json:"session_token"`
	Type         string `json:"type"`
	UserID       string `json:"user_id"`
}
