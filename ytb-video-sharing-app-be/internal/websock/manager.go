package websock

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	ErrEventNotSupported = errors.New("this event type is not supported")
)

// Manager is used to hold references to all Clients Registered, and Broadcasting etc
type Manager struct {
	upgrader websocket.Upgrader

	clients ClientList

	mux      sync.Mutex
	handlers map[string]EventHandler
	otps     RetentionMap
}

func NewManager(rentation RetentionMap) (*Manager, *http.ServeMux) {
	m := &Manager{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		clients: make(ClientList),
		otps:    rentation,
	}

	m.setupEventHandlers()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", m.ServeWS)

	return m, mux
}

// setupEventHandlers configures and adds all handlers
func (m *Manager) setupEventHandlers() {
	m.handlers = make(map[string]EventHandler)
	m.handlers[EventSendMessage] = func(e Event, c *Client) error {
		c.egress <- e
		return nil
	}
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	otp := r.URL.Query().Get("otp")

	if otp == "" {
		// Tell the user its not authorized
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	connID := r.URL.Query().Get("connID")
	if otp == "" {
		// Tell the user its not authorized
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing connID in query"))
		return
	}

	// Verify OTP is existing
	if !m.otps.VerifyOTP(otp) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Println("New connection")
	// Begin by upgrading the HTTP request
	conn, err := m.upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn, m, connID)
	m.addClient(client, connID)

	// Start read/write process
	go client.ReadMessages()
	go client.WriteMessages()
}

// addClient will add clients to our clientList
func (m *Manager) addClient(client *Client, connID string) {
	m.mux.Lock()
	defer m.mux.Unlock()

	// Add clients
	m.clients[connID] = client
}

// removeClient will remove the client and clean up
func (m *Manager) removeClient(client *Client, connID string) {
	m.mux.Lock()
	defer m.mux.Unlock()

	// Check if Client exists, then delete it
	if _, ok := m.clients[connID]; ok {
		// close connection
		client.connection.Close()

		// remove
		delete(m.clients, connID)
	}
}

// routeEvent is used to make sure the correct event goes into the correct handler
func (m *Manager) routeEvent(event Event, c *Client) error {
	// Check if Handler is present in Map
	if handler, ok := m.handlers[event.Type]; ok {
		// Execute the handler and return any err
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
}

func (m *Manager) SendBroadCast(event Event, connIDExclusive string) {
	for connID, client := range m.clients {
		if connID != connIDExclusive {
			log.Printf("Sending to client %s", connID)
			client.egress <- event
		}
	}
}
