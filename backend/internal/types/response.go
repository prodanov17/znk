package types

import "time"

type ErrorResponse struct {
	StatusCode int       `json:"status"`
	Error      string    `json:"error"`
	Path       string    `json:"path"`
	Timestamp  time.Time `json:"timestamp"`
}
