package entity

import "time"

// Host represents a network host with its URL, status, latency, and any associated error.
type Host struct {
	URL     string
	Status  string
	Latency time.Duration
	Error   error
}
