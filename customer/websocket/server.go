package websocket

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type WSServer struct {
	Name    string
	Port    string
	Timeout int

	server *http.Server
	mux    *http.ServeMux
	lock   *sync.Mutex

	routing   map[string]*wsRoute
	idCounter int
	closedCon int
}

func NewWSServer(name string) *WSServer {
	return &WSServer{
		Name:    name,
		lock:    &sync.Mutex{},
		routing: make(map[string]*wsRoute),
		mux:     &http.ServeMux{},
	}
}

// genConId generate new connection id
func (wss *WSServer) genConId() int {
	wss.lock.Lock()
	defer wss.lock.Unlock()

	wss.idCounter++
	return wss.idCounter
}

func (wss *WSServer) NewRoute(path string) *wsRoute {
	wsr := newWsRoute()
	wss.routing[path] = wsr

	// bind connection to this path
	wss.mux.Handle(path, websocket.Handler(func(conn *websocket.Conn) {
		fmt.Println("New connection: " + conn.RemoteAddr().String())
		con := newConnection(wss.genConId(), conn, wss.Timeout)
		wsr.addCon(con)

		if wsr.OnConnect != nil {
			wsr.OnConnect(con)
		}

		wsr.transferring(con)

	}))

	return wsr
}

// GetRoute get coordinator of given path
func (wss *WSServer) GetRoute(path string) *wsRoute {
	return wss.routing[path]
}

// Expose expose port
func (wss *WSServer) Expose(port string) {
	wss.Port = port
}

// GetActiveCon get active connection number in server
func (wss *WSServer) GetActiveCon() int {
	wss.lock.Lock()
	defer wss.lock.Unlock()

	return wss.idCounter - wss.closedCon
}

// Start start server
func (wss *WSServer) Start() {
	if wss.server != nil {
		return
	}

	wss.server = &http.Server{
		Addr:         fmt.Sprintf(":%s", wss.Port),
		Handler:      wss.mux,
		ReadTimeout:  1800 * time.Second, // 30 minutes
		WriteTimeout: 1800 * time.Second, // 30 minutes
	}

	fmt.Printf("[WS Server] %s trying to listen on port %s\n", wss.Name, wss.Port)
	if err := wss.server.ListenAndServe(); err != nil {
		panic(err)
	}
}
