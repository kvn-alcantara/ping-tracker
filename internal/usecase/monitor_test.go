package usecase

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockPinger struct {
	mock.Mock
}

func (m *MockPinger) Ping(ip string) (time.Duration, error) {
	args := m.Called(ip)
	return args.Get(0).(time.Duration), args.Error(1)
}

type MockDisplay struct {
	mock.Mock
}

func (m *MockDisplay) PrintHeader(header string) {
	m.Called(header)
}

func (m *MockDisplay) ClearScreen() {
	m.Called()
}

func (m *MockDisplay) PrintStatus(url, status string, latency time.Duration) {
	m.Called(url, status, latency)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Log(message string) {
	m.Called(message)
}

func setupMonitor(urls []string) (*MockPinger, *MockDisplay, *MockLogger, *Monitor) {
	mockPinger := new(MockPinger)
	mockDisplay := new(MockDisplay)
	mockLogger := new(MockLogger)
	monitor := NewMonitor(mockPinger, mockDisplay, mockLogger, urls)
	return mockPinger, mockDisplay, mockLogger, monitor
}

func runMonitorWithTimeout(monitor *Monitor, duration time.Duration) {
	done := make(chan bool)
	go func() {
		time.Sleep(duration)
		close(done)
	}()
	monitor.StartMonitoring()
	monitor.Run(done)
}

func TestMonitorStartMonitoring(t *testing.T) {
	urls := []string{"example.com", "test.com"}
	mockPinger, mockDisplay, mockLogger, monitor := setupMonitor(urls)

	mockPinger.On("Ping", mock.Anything).Return(50*time.Millisecond, nil).Maybe()
	mockLogger.On("Log", mock.Anything).Return().Maybe()
	mockDisplay.On("PrintHeader", mock.Anything).Return().Maybe()
	mockDisplay.On("ClearScreen").Return().Maybe()
	mockDisplay.On("PrintStatus", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()

	runMonitorWithTimeout(monitor, 100*time.Millisecond)
}

func TestMonitorRunWithOfflineHost(t *testing.T) {
	urls := []string{"example.com"}
	mockPinger, mockDisplay, mockLogger, monitor := setupMonitor(urls)

	mockPinger.On("Ping", mock.Anything).Return(0*time.Millisecond, fmt.Errorf("host unreachable")).Maybe()
	mockLogger.On("Log", mock.Anything).Return().Maybe()
	mockDisplay.On("PrintHeader", mock.Anything).Return().Maybe()
	mockDisplay.On("ClearScreen").Return().Maybe()
	mockDisplay.On("PrintStatus", mock.Anything, "Offline", mock.Anything).Return().Maybe()

	runMonitorWithTimeout(monitor, 100*time.Millisecond)
}

func TestMonitorTickerUpdatesDisplay(t *testing.T) {
	urls := []string{"example.com"}
	mockPinger, mockDisplay, mockLogger, monitor := setupMonitor(urls)

	mockPinger.On("Ping", mock.Anything).Return(50*time.Millisecond, nil).Maybe()
	mockLogger.On("Log", mock.Anything).Return().Maybe()

	mockDisplay.On("PrintHeader", mock.Anything).Return().Maybe()
	mockDisplay.On("ClearScreen").Return().Maybe()
	mockDisplay.On("PrintStatus",
		mock.MatchedBy(func(url string) bool { return url == "example.com" }),
		mock.Anything,
		mock.Anything,
	).Return().Maybe()

	done := make(chan bool)
	go func() {
		time.Sleep(600 * time.Millisecond)
		close(done)
	}()

	monitor.StartMonitoring()
	monitor.Run(done)
}
