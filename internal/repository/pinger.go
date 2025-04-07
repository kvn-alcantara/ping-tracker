package repository

import (
	"fmt"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

type Pinger interface {
	Ping(ip string) (time.Duration, error)
}

type ProBingPinger struct{}

func NewProBingPinger() *ProBingPinger {
	return &ProBingPinger{}
}

func (p *ProBingPinger) Ping(ip string) (time.Duration, error) {
	pinger, err := probing.NewPinger(ip)
	if err != nil {
		return 0, err
	}
	pinger.SetPrivileged(true)
	pinger.Count = 3
	pinger.Timeout = 1 * time.Second

	err = pinger.Run()
	if err != nil {
		return 0, err
	}

	stats := pinger.Statistics()
	if stats.PacketsRecv > 0 {
		return stats.AvgRtt, nil
	}
	return 0, fmt.Errorf("no ping reply received from %s", ip)
}

var _ Pinger = (*ProBingPinger)(nil)
