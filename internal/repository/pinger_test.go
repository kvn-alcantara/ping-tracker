package repository

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHttpPingerPingSuccess(t *testing.T) {
	server := SetupLocalServer(t)

	pinger := NewHttpPinger()

	rtt, err := pinger.Ping(server.Listener.Addr().String())

	assert.NoError(t, err, "expected no error when pinging a reachable IP")
	assert.Greater(t, rtt, time.Duration(0), "expected RTT to be greater than 0")
}

func TestHttpPingerPingUnreachable(t *testing.T) {
	pinger := NewHttpPinger()

	rtt, err := pinger.Ping("203.0.113.0")

	assert.Error(t, err, "expected an error when pinging an unreachable IP")
	assert.Equal(t, time.Duration(0), rtt, "expected RTT to be 0 for an unreachable IP")
}

func TestHttpPingerPingInvalidIP(t *testing.T) {
	pinger := NewHttpPinger()

	rtt, err := pinger.Ping("invalid-ip")

	assert.Error(t, err, "expected an error when pinging an invalid IP")
	assert.Equal(t, time.Duration(0), rtt, "expected RTT to be 0 for an invalid IP")
}

func SetupLocalServer(t *testing.T) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	t.Cleanup(server.Close)

	return server
}
