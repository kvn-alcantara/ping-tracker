package repository

import (
	"crypto/tls"
	"fmt"
	"net/http"
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
	pinger.SetPrivileged(false)
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

type HttpPinger struct {
	client *http.Client
}

func NewHttpPinger() *HttpPinger {
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				MinVersion:         tls.VersionTLS12,
			},
		},
	}
	return &HttpPinger{client: client}
}

func (h *HttpPinger) Ping(ip string) (time.Duration, error) {
	start := time.Now()

	resp, err := h.client.Get("https://" + ip)
	if err != nil {
		resp, err = h.client.Get("http://" + ip)
		if err != nil {
			return 0, err
		}
	}
	defer resp.Body.Close()

	duration := time.Since(start)
	return duration, nil
}

var _ Pinger = (*HttpPinger)(nil)

var _ Pinger = (*ProBingPinger)(nil)
