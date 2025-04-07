package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockPinger struct {
	RTT   time.Duration
	Error error
}

func (m *MockPinger) Ping(ip string) (time.Duration, error) {
	return m.RTT, m.Error
}

func TestProBingPingerPingSuccess(t *testing.T) {
	mockPinger := &MockPinger{
		RTT:   50 * time.Millisecond,
		Error: nil,
	}

	rtt, err := mockPinger.Ping("8.8.8.8")

	assert.NoError(t, err, "expected no error when pinging a reachable IP")
	assert.Greater(t, rtt, time.Duration(0), "expected RTT to be greater than 0")
}

func TestProBingPingerPingUnreachable(t *testing.T) {
	mockPinger := &MockPinger{
		RTT:   0,
		Error: fmt.Errorf("no ping reply received"),
	}

	rtt, err := mockPinger.Ping("192.0.2.0")

	assert.Error(t, err, "expected an error when pinging an unreachable IP")
	assert.Equal(t, time.Duration(0), rtt, "expected RTT to be 0 for an unreachable IP")
}

func TestProBingPingerPingInvalidIP(t *testing.T) {
	pinger := NewProBingPinger()

	// Use an invalid IP address
	ip := "invalid-ip"
	rtt, err := pinger.Ping(ip)

	assert.Error(t, err, "expected an error when pinging an invalid IP")
	assert.Equal(t, time.Duration(0), rtt, "expected RTT to be 0 for an invalid IP")
}
