package websock

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// pongWait is how long we will await a pong response from client
	pongWait = 30 * time.Second
	// pingInterval has to be less than pongWait, We cant multiply by 0.9 to get 90% of time
	// Because that can make decimals, so instead *9 / 10 to get 90%
	// The reason why it has to be less than PingRequency is becuase otherwise it will send a new Ping before getting response
	pingInterval = (pongWait * 9) / 10
)

type ClientList map[string]*Client

type Client struct {
	connection *websocket.Conn

	manager *Manager

	// egress is used to avoid concurrent writes on the WebSocket
	egress chan Event

	connID string
}

func NewClient(conn *websocket.Conn, manager *Manager, connID string) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
		connID:     connID,
	}
}

func (c *Client) ReadMessages() {
	defer func() {
		// graceful Close the Connection once this function is done
		c.manager.removeClient(c, c.connID)
	}()

	// Set Max Size of Messages in Bytes
	c.connection.SetReadLimit(512)
	// Configure Wait time for Pong response, use Current time + pongWait
	// This has to be done here to set the first initial timer.
	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}
	// Configure how to handle Pong responses
	c.connection.SetPongHandler(c.pongHandler)

	for {
		_, payload, err := c.connection.ReadMessage()

		if err != nil {
			// If Connection is closed, we will Recieve an error here
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}

			break
		}

		var request Event

		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error marshalling message: %v", err)
			break // Breaking the connection here might be harsh xD
		}
		// Route the Event
		if err := c.manager.routeEvent(request, c); err != nil {
			log.Println("Error handeling Message: ", err)
		}
	}
}

// pongHandler is used to handle PongMessages for the Client
func (c *Client) pongHandler(pongMsg string) error {
	// Current time + Pong Wait time
	log.Println("pong")
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}

func (c *Client) WriteMessages() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		// graceful close if this triggers a closing
		c.manager.removeClient(c, c.connID)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			// ok will be fail in case egress channel closed
			if !ok {
				// manager has closed this channel, notify FE
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed: ", err)
				}

				// close go routine
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return // closes the connection, should we really
			}
			// Write a Regular text message to the connection
			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
			}

			log.Println("sent message")
		case <-ticker.C:
			log.Println("ping")
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("writemsg error: %v\n", err)
				return
			}
		}
	}
}
