package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	domainuserconnection "github.com/DuongVu089x/interview/customer/domain/user_connection"
	"go.mongodb.org/mongo-driver/bson"
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
	handler := &WebSocketHandler{userConnRepository: userConnRepository, wsServer: wsServer}

	// Start background goroutine to clean user connections every 20 seconds
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			handler.cleanUserConnection()
		}
	}()

	return handler
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
		h.userConnRepository.UpdateUserConnection(context.Background(), &domainuserconnection.UserConnection{
			ID: connID,
		}, &domainuserconnection.UserConnection{
			Status:          domainuserconnection.ConStatus.ACTIVE,
			LastMessageTime: &now,
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

	userConn, err := h.userConnRepository.CreateUserConnection(context.Background(), userConn)
	var outputMsg *WSOutputMessage
	if err != nil {
		outputMsg = &WSOutputMessage{
			Topic: TOPIC_NONE,
			Content: map[string]any{
				"status":  "error",
				"message": "Failed to create user connection",
			},
		}
	} else {
		conn.Attached[USER_CON_ID] = userConn.ID
		outputMsg = &WSOutputMessage{
			Topic: TOPIC_CONNECTED,
			Content: map[string]any{
				"status":  "connected",
				"data":    userConn,
				"message": "Welcome to the server! Host: " + host + ", IP: " + ip + ", Time: " + now.Format(time.RFC3339),
			},
		}
	}
	err = h.pushMessageToDevice(conn, outputMsg.ToString())
	if err != nil {
		fmt.Printf("OnWSConnect Error pushing message to device: %s", err)
	}
	return err
}

func (h *WebSocketHandler) OnWSMessage(conn *Connection, message string) error {
	fmt.Printf("Received message: %s", message)

	if conn.Attached[USER_CON_ID] == nil {
		outputMsg := &WSOutputMessage{
			Topic: TOPIC_NONE,
			Content: map[string]any{
				"status":  "error",
				"message": "User connection not found",
			},
		}
		err := h.pushMessageToDevice(conn, outputMsg.ToString())
		if err != nil {
			fmt.Printf("OnWSMessage Error pushing message to device: %s", err)
		}
		return err
	}

	var msg WSInputMessage
	err := json.Unmarshal([]byte(message), &msg)
	if err != nil {
		outputMsg := &WSOutputMessage{
			Topic: TOPIC_NONE,
			Content: map[string]interface{}{
				"status":  "error",
				"message": "Invalid message format",
			},
		}

		err = h.pushMessageToDevice(conn, outputMsg.ToString())
		if err != nil {
			fmt.Printf("OnWSMessage Error pushing message to device: %s", err)
		}
		return err
	}

	if msg.Topic == "" {
		return errors.New("topic is none")
	}

	switch msg.Topic {
	case TOPIC_AUTHORIZATION:
		var data WSAuthorization
		bytes, _ := json.Marshal(msg.Content)
		err = json.Unmarshal(bytes, &data)
		if err != nil {
			outputMsg := &WSOutputMessage{
				Topic: TOPIC_NONE,
				Content: map[string]any{
					"status":  "error",
					"message": "Invalid message format",
				},
			}
			err = h.pushMessageToDevice(conn, outputMsg.ToString())
			if err != nil {
				fmt.Printf("OnWSMessage Error pushing message to device: %s", err)
			}
			return err
		}

		now := time.Now()
		objID := conn.Attached[USER_CON_ID].(primitive.ObjectID)
		conn.Attached[USER_ID] = data.UserID

		err = h.userConnRepository.UpdateUserConnection(context.Background(), &domainuserconnection.UserConnection{
			ID: objID,
		}, &domainuserconnection.UserConnection{
			UserID:          data.UserID,
			Status:          domainuserconnection.ConStatus.ACTIVE,
			LastUpdatedTime: &now,
		})
		if err != nil {
			outputMsg := &WSOutputMessage{
				Topic: TOPIC_NONE,
				Content: map[string]any{
					"status":  "error",
					"message": "Failed to update user connection",
				},
			}
			err = h.pushMessageToDevice(conn, outputMsg.ToString())
			if err != nil {
				fmt.Printf("OnWSMessage Error pushing message to device: %s", err)
			}
			return err
		}
		outputMsg := &WSOutputMessage{
			Topic: TOPIC_AUTHORIZATION,
			Content: map[string]any{
				"status":    "success",
				"message":   "Authorization successful",
				"errorCode": "",
			},
			Callback: msg.Callback,
		}
		err = h.pushMessageToDevice(conn, outputMsg.ToString())
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

	h.userConnRepository.GetUserConnection(context.Background(), &domainuserconnection.UserConnection{
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

// cleanUserConnection cleans up user connections that haven't been authorized within 3 minutes of connecting
func (h *WebSocketHandler) cleanUserConnection() {
	// log.Println("Start cleaning user connections")
	// defer log.Println("End cleaning user connections")

	route := h.wsServer.GetRoute("/notifications")
	if route == nil {
		log.Println("route not found")
		return
	}

	connMap := route.GetConnectionMap()
	now := time.Now()

	outputMsg := &WSOutputMessage{
		Topic: TOPIC_CONNECTION,
		Content: map[string]interface{}{
			"action":     "PING",
			"serverTime": now,
		},
	}
	for _, conn := range connMap {
		conn.RLock()
		userID := conn.Attached[USER_ID]
		connTime := conn.Attached[CONNECTED_TIME]
		cid := conn.Attached[USER_CON_ID]
		conn.RUnlock()

		// check if client hasn't authorized within 3 minutes of connecting
		if userID == nil && connTime != nil {
			connectedTime := connTime.(time.Time)

			// Close connection if client hasn't authorized within 3 minutes of connecting
			if connectedTime.Add(180 * time.Second).Before(now) {
				outputMsg := &WSOutputMessage{
					Topic: TOPIC_CONNECTION,
					Content: map[string]interface{}{
						"status": "CLOSED",
						"reason": "Too long not authorized.",
					},
				}
				h.pushMessageToDevice(conn, outputMsg.ToString())
				conn.Close()
				return
			}
		}

		// send ping message to device
		err := h.pushMessageToDevice(conn, outputMsg.ToString())
		if err != nil {
			conn.Close()
			if cid != nil {
				connID := cid.(primitive.ObjectID)
				h.userConnRepository.UpdateUserConnection(context.Background(), &domainuserconnection.UserConnection{
					ID: connID,
				}, &domainuserconnection.UserConnection{
					Status:             domainuserconnection.ConStatus.CLOSED,
					DisconnectedTime:   &now,
					DisconnectedReason: "Closed because fail to send ping message.",
				})
			}
		}
	}

	// remove old connections that have not sent a ping message after 1 minute (3 times the job processing time)
	limit := now.Add(-time.Duration(60) * time.Second)
	userConns, err := h.userConnRepository.GetUserConnections(context.Background(), &domainuserconnection.UserConnection{
		Status: domainuserconnection.ConStatus.ACTIVE,
		ComplexQuery: []bson.M{{
			"last_message_time": bson.M{"$lt": limit},
		}},
	}, 0, 100)
	if err != nil {
		return
	}

	for _, userConn := range userConns {
		h.userConnRepository.UpdateUserConnection(context.Background(), &domainuserconnection.UserConnection{
			ID: userConn.ID,
		}, &domainuserconnection.UserConnection{
			Status:             domainuserconnection.ConStatus.CLOSED,
			DisconnectedTime:   &now,
			DisconnectedReason: "Closed by cleaning process.",
		})
	}
}
