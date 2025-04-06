package websocket

import (
	"errors"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type Connection struct {
	ID       int
	Attached map[string]any
	Key      string

	rootConn *websocket.Conn
	isActive bool
	lock     *sync.Mutex
	rwlock   *sync.RWMutex

	timeout      int
	sentMsgCount int
	readMsgCount int
}

func newConnection(id int, ws *websocket.Conn, timeout int) *Connection {
	return &Connection{
		ID:       id,
		Attached: make(map[string]any),

		rootConn: ws,
		isActive: true,
		lock:     &sync.Mutex{},
		rwlock:   &sync.RWMutex{},
		timeout:  timeout,

		sentMsgCount: 0,
		readMsgCount: 0,
	}
}

// send message to client side
func (c *Connection) Send(message string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.isActive {
		return errors.New("connection is not active")
	}

	if c.timeout > 0 {
		// Set a deadline for how long the write operation can take
		// If the write takes longer than timeout seconds, it will return an error
		c.rootConn.SetWriteDeadline(time.Now().Add(time.Duration(c.timeout) * time.Second))
	}

	n, err := c.rootConn.Write([]byte(message))
	if err == nil && n > 0 {
		c.sentMsgCount++
	}

	return err
}

// Read reads a message from the websocket connection
// It returns the message as a string and any error encountered
// If a panic occurs during reading, it will be recovered and an error will be returned
func (c *Connection) Read() (message string, err error) {
	// Recover from any panics that occur during reading
	defer func() {
		if rec := recover(); rec != nil {
			message = ""
			err = errors.New("read error")
		}
	}()

	// Create a buffer to store the message with max payload size
	var msg = make([]byte, c.rootConn.MaxPayloadBytes)

	// Read from the websocket connection into the buffer
	n, err := c.rootConn.Read(msg)
	if err == nil && n > 0 {
		// If read was successful, increment the read message counter
		c.lock.Lock()
		c.readMsgCount++
		c.lock.Unlock()
	}

	// Return the message converted to string, trimmed to actual bytes read
	return string(msg[:n]), err
}

// IsActive check status of connection
func (c *Connection) IsActive() bool {
	return c.isActive
}

// CLose close the connection
func (c *Connection) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.isActive {
		c.rootConn.Close()
		c.isActive = false
	}
}

// Deactive deactive status of connection
func (c *Connection) Deactive() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.isActive = false
}

// GetHeader get header from request
func (c *Connection) GetHeader(name string) string {
	hdr := c.rootConn.Request().Header
	if hdr != nil {
		return hdr.Get(name)
	}
	return ""
}

func (con *Connection) WLock() {
	con.rwlock.Lock()
}

func (con *Connection) WUnlock() {
	con.rwlock.Unlock()
}

func (con *Connection) RLock() {
	con.rwlock.RLock()
}

func (con *Connection) RUnlock() {
	con.rwlock.RUnlock()
}
