package localUtil

import (
	"errors"
	"sync"
	"time"
)

// Default configurations for HTTPCodeMonitor.
const (
	DefaultWheelSize          = 150000
	DefaultTimeLimit          = 300
	DefaultRequestChannelSize = 1000
)

// HTTPCodeMonitor manages a set of sliding windows for monitoring HTTP requests.
type HTTPCodeMonitor struct {
	windows sync.Map
}

var defaultHTTPMonitor *HTTPCodeMonitor

// MonitorConfig contains the configuration for a monitor.
type MonitorConfig struct {
	// Path is the HTTP endpoint to monitor.
	Path string

	// LatencyMultiplier is the multiplier for the latency threshold.
	LatencyMultiplier float64

	// FailureThreshold is the failure rate threshold.
	FailureThreshold float64

	// LatencyCount is the latency count threshold.
	LatencyCount int64

	// FailureCount is the failure count threshold.
	FailureCount int64

	// AlertOnLatency determines whether to alert on latency.
	AlertOnLatency bool

	// LatencyThreshold is the latency threshold.
	LatencyThreshold float64

	// TimeLimit is the time limit for the monitor.
	TimeLimit int64

	// The size for requestChannel.
	RequestChannelSize int64

	// WheelSize is the size of the wheel.
	WheelSize int64

	// The minimum number of requests needed to perform alert checks.
	RequestThreshold int64

	// AlertFunc is the function to call when an alert is triggered.
	AlertFunc func(path string, total int64, failRate float64, avg float64)
}

// SlidingWindow represents a sliding window for monitoring HTTP requests.
type SlidingWindow struct {
	// The path of the HTTP endpoint.
	path string

	// The time wheel.
	wheel *TimeWheel

	// The channel for sending requests to the processing goroutine.
	requestChannel chan Request

	// The size for request channel.
	requestChannelSize int64

	// The time limit for the sliding window.
	timeLimit int64

	// The size of the wheel.
	wheelSize int64

	// The total number of requests.
	totalRequests int64

	// The number of failed requests.
	failedRequests int64

	// The total latency.
	totalLatency float64

	// Whether to alert on latency.
	alertOnLatency bool

	// Latency threshold.
	latencyThreshold float64

	// The count of latency exceed latencyThreshold.
	latencyExceedCount int64

	// The latency alert threshold.
	latencyAlertThreshold int64

	// The multiplier for the latency threshold.
	latencyMultiplier float64

	// The count of failures.
	failureExceedCount int64

	// The failure alert threshold.
	failureAlertThreshold int64

	// The failure rate threshold.
	failureThreshold float64

	// The minimum number of requests needed to perform alert checks.
	requestThreshold int64

	// The function to call when an alert is triggered.
	alertFunc func(path string, total int64, failRate float64, avg float64)

	// The alert cooldown time.
	alertCooldownTime int64
}

// Request represents a request.
type Request struct {
	// Whether the request was successful.
	success bool

	// The latency of the request.
	latency float64

	// The previous average latency.
	oldLatency float64

	// Whether the request exists.
	isExist bool
}

// TimeWheel represents a time wheel.
type TimeWheel struct {
	// The data in the time wheel.
	data []Request

	// The current index in the time wheel.
	index int
}

// GetDefaultHTTPMonitor returns the default HTTP monitor.
func GetDefaultHTTPMonitor() *HTTPCodeMonitor {
	return defaultHTTPMonitor
}

// NewHTTPMonitor creates a new HTTP monitor with the given configurations.
func NewHTTPMonitor(configs []*MonitorConfig) error {
	defaultHTTPMonitor = &HTTPCodeMonitor{
		windows: sync.Map{},
	}
	for _, config := range configs {
		if config.LatencyMultiplier <= 0 {
			config.LatencyMultiplier = 5.00
		}
		if config.LatencyCount <= 0 {
			config.LatencyCount = 5
		}
		if config.FailureThreshold <= 0 {
			config.FailureThreshold = 0.1
		}
		if config.FailureCount <= 0 {
			config.FailureCount = 5
		}
		if config.TimeLimit <= 0 {
			config.TimeLimit = DefaultTimeLimit
		}
		if config.WheelSize <= 0 {
			config.WheelSize = DefaultWheelSize
		}
		if config.RequestChannelSize <= 0 {
			config.RequestChannelSize = DefaultRequestChannelSize
		}
		window, err := newTimeBasedWindow(config)
		if err != nil {
			return err
		}
		defaultHTTPMonitor.windows.Store(config.Path, window)
	}
	return nil
}

// newTimeBasedWindow creates a new time-based window with the given configuration.
func newTimeBasedWindow(config *MonitorConfig) (*SlidingWindow, error) {
	if config.LatencyMultiplier <= 0 {
		return nil, errors.New("invalid latency multiplier in configuration")
	}
	if config.FailureThreshold <= 0 {
		return nil, errors.New("invalid failure threshold in configuration")
	}
	if config.LatencyCount <= 0 {
		return nil, errors.New("invalid latency count in configuration")
	}
	if config.FailureCount <= 0 {
		return nil, errors.New("invalid failure count in configuration")
	}
	base := &SlidingWindow{
		path:                  config.Path,
		failureThreshold:      config.FailureThreshold,
		latencyMultiplier:     config.LatencyMultiplier,
		latencyAlertThreshold: config.LatencyCount,
		failureAlertThreshold: config.FailureCount,
		alertOnLatency:        config.AlertOnLatency,
		latencyThreshold:      config.LatencyThreshold,
		timeLimit:             config.TimeLimit,
		wheelSize:             config.WheelSize,
		wheel: &TimeWheel{
			data:  make([]Request, config.WheelSize),
			index: 0,
		},
		requestThreshold: config.RequestThreshold,
		alertFunc:        config.AlertFunc,
		requestChannel:   make(chan Request, config.RequestChannelSize),
	}
	go base.processRequests()
	go base.rotateWheel()
	return base, nil
}

// Record records a request in the monitor.
func (m *HTTPCodeMonitor) Record(endpoint string, success bool, latency float64) error {
	value, exists := m.windows.Load(endpoint)
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

// updateLatencyExceedCount updates the count of latency exceedances.
func (sw *SlidingWindow) updateLatencyExceedCount(req Request) {
	if !req.isExist {
		return
	}
	decrement := func(c *int64) {
		*c--
		if *c < 0 {
			*c = 0
		}
	}
	decrement(&sw.totalRequests)
	if !req.success {
		decrement(&sw.failedRequests)
	}
	sw.totalLatency -= req.latency
	if sw.totalLatency < 0 {
		sw.totalLatency = 0
	}
	if req.latency > req.oldLatency*sw.latencyMultiplier {
		decrement(&sw.latencyExceedCount)
	}
}

// processRequests processes requests sent through the request channel.
func (sw *SlidingWindow) processRequests() {
	for req := range sw.requestChannel {
		if req.isExist {
			sw.doRecord(req.success, req.latency)
		}
	}
}

// doRecord records a request in the window.
func (sw *SlidingWindow) doRecord(success bool, latency float64) {
	if !success {
		sw.failedRequests++
	}
	sw.totalRequests++
	sw.totalLatency += latency
	avgLat := sw.averageLatency()
	failRate := sw.failureRate()
	sw.wheel.push(Request{
		success:    success,
		latency:    latency,
		oldLatency: avgLat,
		isExist:    true,
	})
	if sw.alertOnLatency {
		sw.updateCounterAndAlert(
			latency > sw.latencyThreshold && latency > avgLat*sw.latencyMultiplier,
			&sw.latencyExceedCount,
			sw.latencyAlertThreshold,
		)
	}
	sw.updateCounterAndAlert(
		failRate > sw.failureThreshold,
		&sw.failureExceedCount,
		sw.failureAlertThreshold,
	)
}

// record records a request in the window.
func (sw *SlidingWindow) record(success bool, latency float64) {
	sw.requestChannel <- Request{success: success, latency: latency}
}

// updateCounterAndAlert updates the counter and potentially triggers an alert.
func (sw *SlidingWindow) updateCounterAndAlert(condition bool, count *int64, threshold int64) {
	if sw.requestThreshold > 0 && sw.totalRequests < sw.requestThreshold {
		return
	}

	if condition {
		*count++
		if sw.alertFunc != nil && *count >= threshold {
			now := time.Now().Unix()
			if now-sw.alertCooldownTime >= 60*60 {
				sw.alertFunc(sw.path, sw.totalRequests, sw.failureRate(), sw.averageLatency())
				sw.alertCooldownTime = now
			}
			*count = 0
		}
	} else if *count > 0 {
		*count--
	} else {
		*count = 0
	}
}

// failureRate computes the failure rate.
func (sw *SlidingWindow) failureRate() float64 {
	if sw.totalRequests == 0 {
		return 0
	}
	return float64(sw.failedRequests) / float64(sw.totalRequests)
}

// averageLatency computes the average latency.
func (sw *SlidingWindow) averageLatency() float64 {
	if sw.totalRequests == 0 {
		return 0
	}
	return sw.totalLatency / float64(sw.totalRequests)
}

// push pushes a request onto the time wheel.
func (tw *TimeWheel) push(req Request) {
	tw.data[tw.index] = req
	tw.index = (tw.index + 1) % len(tw.data)
}

// rotateWheel rotates the wheel.
func (sw *SlidingWindow) rotateWheel() {
	ticker := time.NewTicker(time.Duration(sw.timeLimit) * time.Second / time.Duration(sw.wheelSize))
	defer ticker.Stop()
	for range ticker.C {
		oldReq := sw.wheel.data[sw.wheel.index]
		sw.updateLatencyExceedCount(oldReq)
		sw.wheel.push(Request{})
	}
}
