package websock

import "encoding/json"

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// receive event to route message into oke handler
// client pass in here when want to response client
type EventHandler func(event Event, c *Client) error

const (
	EventSendMessage = "send_message"
	EventNotif       = "event_notif"
)

type EventNotificationMessage struct {
	Title     string `json:"title"`
	SharedBy  string `json:"shared_by"`
	Thumbnail string `json:"thumbnail"`
}
