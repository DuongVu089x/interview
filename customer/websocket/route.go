package websocket

import (
	"errors"
	"sync"
)

type wsRoute struct {
	OnConnect func(con *Connection) error
	OnMessage func(con *Connection, message string) error
	OnClose   func(con *Connection, err error) error

	conMapMutex sync.Mutex
	connMap     map[int]*Connection
	payloadSize int
}

func newWsRoute() *wsRoute {
	return &wsRoute{
		connMap:     make(map[int]*Connection),
		payloadSize: 1024,
		conMapMutex: sync.Mutex{},
	}
}

func (wsr *wsRoute) addCon(con *Connection) {
	wsr.conMapMutex.Lock()
	defer wsr.conMapMutex.Unlock()

	wsr.connMap[con.ID] = con
}

func (wsr *wsRoute) removeCon(id int) {
	wsr.conMapMutex.Lock()
	defer wsr.conMapMutex.Unlock()

	delete(wsr.connMap, id)
}

// keep chatting with client
func (wsr *wsRoute) transferring(conn *Connection) {

	// recover from panic make sure connection is closed
	defer func() {
		if err := recover(); err != nil {
			if conn != nil {
				conn.Deactive()
				wsr.removeCon(conn.ID)
				if wsr.OnClose != nil {
					wsr.OnClose(conn, errors.New("panic from transferring & recover"))
				}
			}
		}
	}()

	// keep chatting with client
	for {
		// override max payload size
		if conn.rootConn.MaxPayloadBytes != wsr.payloadSize {
			conn.rootConn.MaxPayloadBytes = wsr.payloadSize
		}

		payload, err := conn.Read()
		if err != nil {
			conn.Deactive()
			wsr.removeCon(conn.ID)
			if wsr.OnClose != nil {
				wsr.OnClose(conn, err)
			}
			return
		}
		if wsr.OnMessage != nil {
			wsr.OnMessage(conn, payload)
		}
	}
}

func (wsr *wsRoute) GetConnection(conId int) *Connection {
	wsr.conMapMutex.Lock()
	defer wsr.conMapMutex.Unlock()
	return wsr.connMap[conId]
}

func (wsr *wsRoute) GetConnectionMap() map[int]*Connection {
	wsr.conMapMutex.Lock()
	defer wsr.conMapMutex.Unlock()

	var tempMap = make(map[int]*Connection, len(wsr.connMap))
	for id, conn := range wsr.connMap {
		tempMap[id] = conn
	}
	return tempMap
}
