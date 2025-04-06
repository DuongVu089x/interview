package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	domainuserconnection "github.com/DuongVu089x/interview/customer/domain/user_connection"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var USER_CON_ID = "userConId"
var USER_ID = "userId"
var CONNECTED_TIME = "connectedTime"
var USER_AGENT = "userAgent"

type WebSocketHandler struct {
	userConnRepository domainuserconnection.Repository
	wsServer           *WSServer
}

func NewWebSocketHandler(userConnRepository domainuserconnection.Repository, wsServer *WSServer) *WebSocketHandler {
	return &WebSocketHandler{userConnRepository: userConnRepository, wsServer: wsServer}
}

func (h *WebSocketHandler) pushMessageToDevice(conn *Connection, message string) error {
	if !conn.IsActive() {
		return errors.New("connection is not active")
	}
	err := conn.Send(message)
	if err != nil {
		fmt.Printf("Error sending message to device: %s", err)
	}

	conn.RLock()
	connAnyID := conn.Attached[USER_CON_ID]
	conn.RUnlock()

	if connAnyID == nil {
		return errors.New("user connection id not found")
	}

	connID := connAnyID.(primitive.ObjectID)

	// conn.RLock()
	// userID := conn.Attached[USER_ID]
	// conn.RUnlock()

	// if userID == nil {
	// 	return
	// }

	// Only update connected time if message was sent successfully
	if err == nil {
		// Get the last connected time from connection metadata
		conn.RLock()
		connTime := conn.Attached[CONNECTED_TIME]
		conn.RUnlock()

		// Convert interface{} to time.Time
		connectedTime := connTime.(time.Time)
		now := time.Now()

		// If less than 30 seconds have passed since last update, skip updating
		if now.Sub(connectedTime).Seconds() < 30 {
			return nil
		}

		// Update the connected time to current time
		conn.WLock()
		conn.Attached[CONNECTED_TIME] = now
		conn.WUnlock()

		// Update the user connection status to active
		h.userConnRepository.UpdateUserConnection(&domainuserconnection.UserConnection{
			ID: connID,
		}, &domainuserconnection.UserConnection{
			Status:          domainuserconnection.ConStatus.ACTIVE,
			LastUpdatedTime: &now,
		})
	}

	return err
}

func (h *WebSocketHandler) OnWSConnect(conn *Connection) error {
	now := time.Now()
	host, _ := os.Hostname()
	version := os.Getenv("VERSION")
	ip := conn.GetHeader("X-Real-IP")

	userConn := &domainuserconnection.UserConnection{
		Status:        domainuserconnection.ConStatus.ACTIVE,
		UserID:        "UNKNOWN",
		ConnectedTime: &now,
		WSHost:        host,
		WSHostVersion: version,
		IP:            ip,
		UserAgent:     conn.GetHeader("User-Agent"),

		LastUpdatedTime: &now,

		ConnectionLocalID: conn.ID,
	}
	conn.Attached[CONNECTED_TIME] = now
	conn.Attached[USER_AGENT] = userConn.UserAgent

	userConn, err := h.userConnRepository.CreateUserConnection(userConn)
	var outputMsg *WSOutputMessage
	if err != nil {
		outputMsg = &WSOutputMessage{
			Topic: TopicEnum.NONE,
			Content: map[string]any{
				"status":  "error",
				"message": "Failed to create user connection",
			},
		}
	} else {
		conn.Attached[USER_CON_ID] = userConn.ID
		outputMsg = &WSOutputMessage{
			Topic: TopicEnum.CONNECTED,
			Content: map[string]any{
				"status":  "connected",
				"data":    userConn,
				"message": "Welcome to the server! Host: " + host + ", IP: " + ip + ", Time: " + now.Format(time.RFC3339),
			},
		}
	}
	err = h.pushMessageToDevice(conn, outputMsg.toString())
	if err != nil {
		fmt.Printf("OnWSConnect Error pushing message to device: %s", err)
	}
	return err
}

func (h *WebSocketHandler) OnWSMessage(conn *Connection, message string) error {
	fmt.Printf("Received message: %s", message)

	if conn.Attached[USER_CON_ID] == nil {
		outputMsg := &WSOutputMessage{
			Topic: TopicEnum.NONE,
			Content: map[string]any{
				"status":  "error",
				"message": "User connection not found",
			},
		}
		err := h.pushMessageToDevice(conn, outputMsg.toString())
		if err != nil {
			fmt.Printf("OnWSMessage Error pushing message to device: %s", err)
		}
		return err
	}

	var msg WSInputMessage
	err := json.Unmarshal([]byte(message), &msg)
	if err != nil {
		outputMsg := &WSOutputMessage{
			Topic: TopicEnum.NONE,
			Content: map[string]interface{}{
				"status":  "error",
				"message": "Invalid message format",
			},
		}

		err = h.pushMessageToDevice(conn, outputMsg.toString())
		if err != nil {
			fmt.Printf("OnWSMessage Error pushing message to device: %s", err)
		}
		return err
	}

	if string(msg.Topic) == "" {
		return errors.New("topic is none")
	}

	switch msg.Topic {
	case TopicEnum.AUTHORIZATION:
		var data WSAuthorization
		bytes, _ := json.Marshal(msg.Content)
		err = json.Unmarshal(bytes, &data)
		if err != nil {
			outputMsg := &WSOutputMessage{
				Topic: TopicEnum.NONE,
				Content: map[string]any{
					"status":  "error",
					"message": "Invalid message format",
				},
			}
			err = h.pushMessageToDevice(conn, outputMsg.toString())
			if err != nil {
				fmt.Printf("OnWSMessage Error pushing message to device: %s", err)
			}
			return err
		}

		now := time.Now()
		objID := conn.Attached[USER_CON_ID].(primitive.ObjectID)

		err = h.userConnRepository.UpdateUserConnection(&domainuserconnection.UserConnection{
			ID: objID,
		}, &domainuserconnection.UserConnection{
			UserID:          data.UserID,
			Status:          domainuserconnection.ConStatus.ACTIVE,
			LastUpdatedTime: &now,
		})
		if err != nil {
			outputMsg := &WSOutputMessage{
				Topic: TopicEnum.NONE,
				Content: map[string]any{
					"status":  "error",
					"message": "Failed to update user connection",
				},
			}
			err = h.pushMessageToDevice(conn, outputMsg.toString())
			if err != nil {
				fmt.Printf("OnWSMessage Error pushing message to device: %s", err)
			}
			return err
		}
		outputMsg := &WSOutputMessage{
			Topic: TopicEnum.AUTHORIZATION,
			Content: map[string]any{
				"status":    "success",
				"message":   "Authorization successful",
				"errorCode": "",
			},
			Callback: msg.Callback,
		}
		err = h.pushMessageToDevice(conn, outputMsg.toString())
		if err != nil {
			fmt.Printf("OnWSMessage Error pushing message to device: %s", err)
		}
		return err
	}

	return nil
}

func (h *WebSocketHandler) OnWSClose(conn *Connection, err error) error {
	if conn.Attached[USER_CON_ID] == nil {
		return errors.New("user connection id not found")
	}

	if err == nil || err.Error() == "" {
		return errors.New("error is nil")
	}

	now := time.Now()

	connAnyID := conn.Attached[USER_CON_ID]
	if connAnyID == nil {
		return errors.New("user connection id not found")
	}
	connID := connAnyID.(primitive.ObjectID)

	h.userConnRepository.GetUserConnection(&domainuserconnection.UserConnection{
		ID:              connID,
		Status:          domainuserconnection.ConStatus.CLOSED,
		LastUpdatedTime: &now,
	})

	return nil
}

func (h *WebSocketHandler) GetConnByID(connID int) (*Connection, error) {
	route := h.wsServer.GetRoute("/notifications")
	if route == nil {
		return nil, errors.New("route not found")
	}

	return route.GetConnection(connID), nil
}

// WSOutputMessage
type TopicEnumValue string

// WSInputMessage
type WSInputMessage struct {
	Topic    TopicEnumValue `json:"topic,omitempty"`
	Content  map[string]any `json:"content,omitempty"`
	Callback string         `json:"callback,omitempty"`
}

type WSOutputMessage struct {
	Topic    TopicEnumValue `json:"topic,omitempty"`
	Content  map[string]any `json:"content,omitempty"`
	Callback string         `json:"callback,omitempty"`
}

type WSAuthorization struct {
	SessionToken string `json:"session_token"`
	Type         string `json:"type"`
	UserID       string `json:"user_id"`
}

func (o *WSOutputMessage) toString() string {
	bytes, err := json.Marshal(o)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// Define topic enum values as package variables
var TopicEnum = struct {
	NONE          TopicEnumValue
	CONNECTED     TopicEnumValue
	CONNECTION    TopicEnumValue
	PING          TopicEnumValue
	AUTHORIZATION TopicEnumValue
	ANNOUNCEMENT  TopicEnumValue
	EVENT         TopicEnumValue
}{
	NONE:          "NONE",
	CONNECTED:     "CONNECTED",
	CONNECTION:    "CONNECTION",
	PING:          "PING",
	AUTHORIZATION: "AUTHORIZATION",
	ANNOUNCEMENT:  "ANNOUNCEMENT",
	EVENT:         "EVENT",
}
