package monitor

import (
	"errors"
	"sync"
)

// HTTPCodeMonitor manages a set of sliding windows for monitoring HTTP requests.
type HTTPCodeMonitor struct {
	Windows sync.Map
}

var defaultHTTPMonitor *HTTPCodeMonitor

// GetDefaultHTTPMonitor returns the default HTTP monitor.
func GetDefaultHTTPMonitor() *HTTPCodeMonitor {
	return defaultHTTPMonitor
}

// NewHTTPMonitor creates a new HTTP monitor with the given configurations.
func NewHTTPMonitor(configs []*Config) error {
	defaultHTTPMonitor = &HTTPCodeMonitor{
		Windows: sync.Map{},
	}
	for _, config := range configs {
		if config.LatencyMultiplier <= 0 {
			config.LatencyMultiplier = defaultLatencyMultiplier
		}
		if config.LatencyCount <= 0 {
			config.LatencyCount = defaultLatencyCount
		}
		if config.FailureThreshold <= 0 {
			config.FailureThreshold = defaultFailureThreshold
		}
		if config.FailureCount <= 0 {
			config.FailureCount = defaultFailureCount
		}
		if config.TimeLimit <= 0 {
			config.TimeLimit = defaultTimeLimit
		}
		if config.WheelSize <= 0 {
			config.WheelSize = defaultWheelSize
		}
		if config.RequestChannelSize <= 0 {
			config.RequestChannelSize = defaultRequestChannelSize
		}
		if config.RequestThreshold <= 0 {
			config.RequestThreshold = defaultRequestThreshold
		}
		window, err := newTimeBasedWindow(config)
		if err != nil {
			return err
		}
		defaultHTTPMonitor.Windows.Store(config.Path, window)
	}
	return nil
}

// Record records a request in the monitor.
func (m *HTTPCodeMonitor) Record(endpoint string, success bool, latency float64) error {
	value, exists := m.Windows.Load(endpoint)
	if !exists {
		return errors.New("endpoint not found: " + endpoint)
	}
	window, ok := value.(*SlidingWindow)
	if !ok {
		return errors.New("type assertion failed for endpoint: " + endpoint)
	}
	window.record(success, latency)
	return nil
}
