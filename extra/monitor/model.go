package monitor

import (
	"sync"
	"time"
)

// HTTPCodeMonitor manages a set of sliding windows for monitoring HTTP requests.
type HTTPCodeMonitor struct {
	windows sync.Map
}

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

	// The alert cool down timestamp.
	alertCoolDownTimeByFailed  int64
	alertCoolDownTimeByLatency int64

	lastRotationTime time.Time
	sampleRate       int64
	requestCount     int64
}

// Request represents a request.
type Request struct {
	// Whether the request was successful.
	success bool

	// The latency of the request.
	latency float64

	// The previous average latency.
	oldAvgLatency float64

	// Whether the request exists.
	isValid bool
}

// TimeWheel represents a time wheel.
type TimeWheel struct {
	// The data in the time wheel.
	data []Request

	// The current index in the time wheel.
	index int
}
