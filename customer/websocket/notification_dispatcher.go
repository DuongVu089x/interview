package websocket

import (
	"errors"

	domainnotification "github.com/DuongVu089x/interview/customer/domain/notification"
	domainuserconnection "github.com/DuongVu089x/interview/customer/domain/user_connection"
)

// NotificationDispatcher handles dispatching notifications to connected websocket clients
type NotificationDispatcher struct {
	routePath        string            // Path for the websocket route
	wsServer         *WSServer         // Reference to the websocket server
	webSocketHandler *WebSocketHandler // Reference to the websocket handler
}

// NewNotificationDispatcher creates a new notification dispatcher instance
func NewNotificationDispatcher(wsServer *WSServer, routePath string, webSocketHandler *WebSocketHandler) *NotificationDispatcher {
	return &NotificationDispatcher{wsServer: wsServer, routePath: routePath, webSocketHandler: webSocketHandler}
}

// DispatchToUser sends a notification to all active connections for a specific user
func (n *NotificationDispatcher) DispatchToUser(userID string, notification *domainnotification.Notification) error {
	// Get the websocket route
	route := n.wsServer.GetRoute(n.routePath)
	if route == nil {
		return errors.New("route not found")
	}

	// Get all active connections for the user
	userConns, err := n.webSocketHandler.userConnRepository.GetUserConnections(&domainuserconnection.UserConnection{
		UserID: userID,
		Status: domainuserconnection.ConStatus.ACTIVE,
	}, 0, 100)
	if err != nil {
		return err
	}

	// Create the notification message
	outputMsg := &WSOutputMessage{
		Topic: TopicEnum.ANNOUNCEMENT,
		Content: map[string]any{
			"topic":       notification.Topic,
			"title":       notification.Title,
			"description": notification.Description,
			"link":        notification.Link,
		},
	}

	// Send notification to each active connection
	for _, userConn := range userConns {
		con, err := n.webSocketHandler.GetConnByID(int(userConn.ConnectionLocalID))
		if err != nil {
			continue
		}

		if con == nil {
			continue
		}

		err2 := n.webSocketHandler.pushMessageToDevice(con, outputMsg.toString())
		if err2 != nil {
			return err2
		}
	}
	return nil
}
