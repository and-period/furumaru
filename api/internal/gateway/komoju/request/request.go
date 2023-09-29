package request

import "time"

type PingPayload struct {
	ID        string    `json:"id"`
	Resource  string    `json:"resource"`
	URL       string    `json:"url"`
	Active    bool      `json:"active"`
	Events    []string  `json:"event_list"`
	CreatedAt time.Time `json:"created_at"`
}
