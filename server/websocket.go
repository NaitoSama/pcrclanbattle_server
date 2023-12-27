package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"pcrclanbattle_server/common"
	"sync"
)

var lock sync.RWMutex

var Server *WebSocketServer

// upgrader Upgrade from request to ws conn
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections (for testing purposes)
		return true
	},
}

// Client represents a connected ws client
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// WebSocketServer
type WebSocketServer struct {
	clients    map[*Client]bool // ws client
	broadcast  chan []byte      // broadcast content
	register   chan *Client     // register a client
	unregister chan *Client     // unregister a client
}

func (server *WebSocketServer) run() {
	for {
		select {
		case client := <-server.register:
			lock.Lock()
			server.clients[client] = true
			lock.Unlock()
		case client := <-server.unregister:
			if _, ok := server.clients[client]; ok {
				lock.Lock()
				delete(server.clients, client)
				lock.Unlock()
				close(client.send)
			}
		case message := <-server.broadcast:
			for client := range server.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(server.clients, client)
				}
			}
		}
	}
}

func (server *WebSocketServer) HandleConnection(context *gin.Context) {
	conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	client := &Client{conn: conn, send: make(chan []byte, 256)}
	server.register <- client

	// todo send conn boss state content when it's first time

	go client.write()
	client.read()
}

func (client *Client) read() {
	defer func() {
		client.conn.Close()
	}()

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			break
		}
		// todo handle the content from ws client
		Server.broadcast <- message
	}
}

func (client *Client) write() {
	defer func() {
		client.conn.Close()
	}()

	for {
		message, ok := <-client.send
		if !ok {
			break
		}
		err := client.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
}

// WSInit Start websocket server
func WSInit() {
	server := &WebSocketServer{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	Server = server
	go Server.run()
	common.Logln(0, "websocket server started")
}
