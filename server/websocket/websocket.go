package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSMessage struct {
	Type    string `json:"type"`    // broadcast or private
	Target  string `json:"target"`  // UserId for private messages
	Content string `json:"content"` // The actual data
	Sender  string `json:"sender"`  // Set by the server
}

type Client struct {
	UserId string
	Conn   *websocket.Conn
	Send   chan WSMessage
}

type Hub struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan WSMessage
	Direct     chan WSMessage
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan WSMessage),
		Direct:     make(chan WSMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.UserId] = client

		case client := <-h.Unregister:
			if _, ok := h.Clients[client.UserId]; ok {
				delete(h.Clients, client.UserId)
				close(client.Send)
			}

		case msg := <-h.Broadcast:
			for id, client := range h.Clients {
				select {
				case client.Send <- msg:
				default:
					close(client.Send)
					delete(h.Clients, id)
				}
			}

		case msg := <-h.Direct:
			if client, ok := h.Clients[msg.Target]; ok {
				select {
				case client.Send <- msg:
				default:
					close(client.Send)
					delete(h.Clients, msg.Target)
				}
			}
		}
	}
}

// broadcast to all
func handleBroadcast(hub *Hub, msg WSMessage) {
	hub.Broadcast <- msg
}

// process then broadcast
func handleProcessedBroadcast(hub *Hub, msg WSMessage) {
	msg.Content = "ALARM: " + msg.Content // Example processing
	hub.Broadcast <- msg
}

// broadcast direct to client x
func handlePrivate(hub *Hub, msg WSMessage) {
	hub.Direct <- msg
}

// process then broadcast to client x
func handleProcessedPrivate(hub *Hub, msg WSMessage) {
	msg.Content = "SECURE MSG: " + msg.Content // Example processing
	hub.Direct <- msg
}

func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		var msg WSMessage
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			break
		}
		msg.Sender = c.UserId

		// ROUTING LOGIC
		switch msg.Type {
		case "broadcast":
			handleBroadcast(hub, msg)
		case "broadcast_special":
			handleProcessedBroadcast(hub, msg)
		case "private":
			handlePrivate(hub, msg)
		case "private_special":
			handleProcessedPrivate(hub, msg)
		}
	}
}

func (c *Client) WritePump() {
	for msg := range c.Send {
		c.Conn.WriteJSON(msg)
	}
}

func HandleWebsocket(hub *Hub, ctx *gin.Context) {
	conn, _ := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	userId := ctx.Query("id") // Example: /ws?id=alice

	client := &Client{
		UserId: userId,
		Conn:   conn,
		Send:   make(chan WSMessage, 256),
	}

	hub.Register <- client

	go client.WritePump()
	go client.ReadPump(hub)
}
