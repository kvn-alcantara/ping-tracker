package repository

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

// Pinger is an interface that defines the Ping method for checking the reachability of an IP address.
type Pinger interface {
	Ping(ip string) (time.Duration, error)
}

// ProBingPinger is an implementation of the Pinger interface using the pro-bing library.
type ProBingPinger struct{}

// NewProBingPinger creates a new instance of ProBingPinger.
func NewProBingPinger() *ProBingPinger {
	return &ProBingPinger{}
}

// Ping sends ICMP echo requests to the specified IP address and returns the average round-trip time or an error.
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

// HTTPPinger is an implementation of the Pinger interface using HTTP requests.
type HTTPPinger struct {
	client *http.Client
}

// NewHTTPPinger creates a new instance of HTTPPinger with a configured HTTP client.
func NewHTTPPinger() *HTTPPinger {
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				MinVersion:         tls.VersionTLS12,
			},
		},
	}
	return &HTTPPinger{client: client}
}

// Ping sends an HTTP GET request to the specified IP address and returns the round-trip time or an error.
// It tries HTTPS first, and if it fails, it falls back to HTTP.
func (h *HTTPPinger) Ping(ip string) (time.Duration, error) {
	start := time.Now()

	resp, err := h.client.Get("https://" + ip)
	if err != nil {
		resp, err = h.client.Get("http://" + ip)
		if err != nil {
			return 0, err
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("received non-success status code %d from %s", resp.StatusCode, ip)
	}

	duration := time.Since(start)
	return duration, nil
}

var _ Pinger = (*HTTPPinger)(nil)

var _ Pinger = (*ProBingPinger)(nil)
