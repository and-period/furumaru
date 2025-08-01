package request

import "time"

type PingRequest struct {
	ID        string       `json:"id"`
	Type      string       `json:"type"`
	Resource  string       `json:"resource"`
	Payload   *PingPayload `json:"data"`
	CreatedAt time.Time    `json:"created_at"`
	Reason    string       `json:"reason"`
}

type PingPayload struct {
	ID        string    `json:"id"`
	Resource  string    `json:"resource"`
	URL       string    `json:"url"`
	Active    bool      `json:"active"`
	Events    []string  `json:"event_list"`
	CreatedAt time.Time `json:"created_at"`
}
