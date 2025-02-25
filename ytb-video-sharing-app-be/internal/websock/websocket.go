package websock

import (
	"fmt"
	"net/http"
	"sync"
	"time"
	"ytb-video-sharing-app-be/utils"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WebSocketServerInterface interface {
	HandleWebSocketConnection(w http.ResponseWriter, r *http.Request)
	SendMessage(accountID int64, event *EventMessage)
	CloseSingleConnection(accountID int64, connID string)
}

type webSocketServer struct {
	clients    map[int64]map[string]*websocket.Conn
	upgrader   websocket.Upgrader
	clientMux  sync.Mutex
	keyManager *utils.KeyManager
}

type EventMessage struct {
	Title     string `json:"title,omitempty"`
	SharedBy  string `json:"shared_by,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
}

func NewWebSocketServer(keyManager *utils.KeyManager) (WebSocketServerInterface, *http.ServeMux) {
	server := &webSocketServer{
		clients: make(map[int64]map[string]*websocket.Conn),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// TODO: allow for only FE
				return true
			},
		},
		keyManager: keyManager,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", server.HandleWebSocketConnection)

	return server, mux
}

func (s *webSocketServer) HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	// check user have access token (if not, not initialize websocket)
	accessToken := r.Header.Get("Sec-WebSocket-Protocol")

	claims, errToken := utils.ValidateToken(accessToken, s.keyManager)

	if errToken != nil {
		fmt.Println("WebSocket upgrade error:", errToken)
		http.Error(w, errToken.Error(), http.StatusUnauthorized)
		return
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		http.Error(w, errToken.Error(), http.StatusInternalServerError)
		return
	}

	connID := uuid.New().String()

	// lock resource which manages connections
	s.clientMux.Lock()
	if _, exists := s.clients[claims.AccountID]; !exists {
		s.clients[claims.AccountID] = make(map[string]*websocket.Conn)
	}
	s.clients[claims.AccountID][connID] = conn
	s.clientMux.Unlock()

	// send conn id to user connect
	conn.WriteJSON(map[string]string{"conn_id": connID})
	fmt.Println("send message oke!")

	heartBeatCh := make(chan struct{})
	doneCh := make(chan struct{})

	// monitor heart beat
	go s.monitorHeartBeat(claims.AccountID, connID, conn, heartBeatCh, doneCh)

	go s.listenMessage(claims.AccountID, connID, conn, heartBeatCh, doneCh)
}

// monitor heart beat to check user active
func (s *webSocketServer) monitorHeartBeat(accountID int64, connID string, conn *websocket.Conn, heartBeatCh chan struct{}, doneCh chan struct{}) {
	timer := time.NewTimer(5 * time.Minute)

	for {
		select {
		case <-heartBeatCh:
			// reset timer if receive heart beat
			timer.Reset(5 * time.Minute)
		case <-timer.C:
			fmt.Printf("Heartbeat timeout for account: %v, conn: %s\n", accountID, connID)
			s.shutdown(accountID, connID, conn, doneCh)
			return
		case <-doneCh:
			// if receive signal close connection, exit go routine monitor heart beat
			return
		}
	}
}

// listen message from user to determine user is active or inactive
func (s *webSocketServer) listenMessage(accountID int64, connID string, conn *websocket.Conn, heartBeatCh chan struct{}, doneCh chan struct{}) {
	defer s.shutdown(accountID, connID, conn, doneCh)

	for {
		select {
		case <-doneCh:
			return
		default:
			_, message, err := conn.ReadMessage()

			if err != nil {
				fmt.Println("WebSocket read error:", err)
				return
			}

			if string(message) == "ping" {
				// receive ping, signal reset timer
				heartBeatCh <- struct{}{}
			}
		}
	}
}

// graceful shutdown connection websocket
func (s *webSocketServer) shutdown(accountID int64, connID string, conn *websocket.Conn, doneCh chan struct{}) {
	s.clientMux.Lock()
	defer s.clientMux.Unlock()

	delete(s.clients[accountID], connID)
	if len(s.clients[accountID]) == 0 {
		delete(s.clients, accountID)
	}

	close(doneCh)
	conn.Close()
	fmt.Println("Closed WebSocket for user:", accountID)
}

func (s *webSocketServer) SendMessage(accountID int64, event *EventMessage) {
	if userConns, exists := s.clients[accountID]; exists {
		for _, conn := range userConns {
			conn.WriteJSON(event)
		}
	}
}

func (s *webSocketServer) CloseSingleConnection(accountID int64, connID string) {
	s.clientMux.Lock()
	defer s.clientMux.Unlock()
	s.clients[accountID][connID].Close()
	delete(s.clients[accountID], connID)
	if len(s.clients[accountID]) == 0 {
		delete(s.clients, accountID)
	}
}
