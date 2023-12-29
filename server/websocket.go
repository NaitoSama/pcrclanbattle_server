package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"pcrclanbattle_server/common"
	"pcrclanbattle_server/db"
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

// run start ws server
func (server *WebSocketServer) run() {
	for {
		select {
		case client := <-server.register: // a new conn
			lock.Lock()
			server.clients[client] = true
			lock.Unlock()
		case client := <-server.unregister: // a conn closed
			if _, ok := server.clients[client]; ok {
				lock.Lock()
				delete(server.clients, client)
				lock.Unlock()
				close(client.send)
			}
		case message := <-server.broadcast: // a broadcast event occurred
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

// HandleConnection upgrade http request to websocket connection
func (server *WebSocketServer) HandleConnection(context *gin.Context) {
	conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	client := &Client{conn: conn, send: make(chan []byte, 8)}
	server.register <- client

	// send conn boss state data when it's first time
	lock.RLock()
	data, _ := json.Marshal(db.Cache.Bosses)
	lock.RUnlock()
	client.send <- data

	go client.write()
	client.read()
}

// read receive messages sent from user
func (client *Client) read() {
	defer func() {
		client.conn.Close()
	}()

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			break
		}
		// handle the content from ws client
		err = informationDiversion(message)
		if err != nil {
			feedback, _ := json.Marshal(gin.H{"result": err.Error()})
			client.send <- []byte(feedback)
		}
	}
}

// write send message to user
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

func informationDiversion(message []byte) error {
	data := make(map[string]string)
	json.Unmarshal(message, &data)
	//if err != nil {
	//	return err
	//}
	dataType, ok := data["type"]
	if !ok {
		return errors.New("need type")
	}
	token, ok := data["token"]
	if !ok {
		return errors.New("need token")
	}
	user := common.MyClaims{}
	user.UserID, user.UserName, user.UserAuthority, ok = common.ParseJWT(token)
	if !ok {
		return errors.New("permission denied")
	}
	switch dataType {
	case "attack":
		err := AttackBoss(message, user.UserName)
		if err != nil {
			return err
		}
		return nil
	case "revise":
		return nil
	case "undo":
		return nil
	case "imin":
		return nil
	case "imout":
		return nil
	default:
		return errors.New("unknown type")
	}
}
