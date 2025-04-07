package entity

import "time"

type Host struct {
	URL     string
	Status  string
	Latency time.Duration
	Error   error
}
