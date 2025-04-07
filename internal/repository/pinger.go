package repository

import "time"

type Pinger interface {
	Ping(ip string) (time.Duration, error)
}
