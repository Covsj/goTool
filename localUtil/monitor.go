package localUtil

import (
	"errors"
	"sync"
	"time"
)

const (
	WheelSize        = 2000
	defaultTimeLimit = 4000
)

// MonitorConfig Configuration structure for the Monitor
type MonitorConfig struct {
	Path           string
	LatCond        float64
	FailCond       float64
	LatThresh      int64
	FailThresh     int64
	AlertOnLatency bool
	LatencyCond    float64
	TimeLimit      int64
	AlertFunc      func(path string, total int64, failRate float64, avg float64)
}

var DefaultHttpMonitor *HttpCodeMonitor

// NewHttpMonitor Function to create a new HTTP monitor
func NewHttpMonitor(configs []MonitorConfig) error {
	DefaultHttpMonitor = &HttpCodeMonitor{
		windows: sync.Map{},
	}
	for _, config := range configs {
		// Default values for the configuration
		if config.LatCond <= 0 {
			config.LatCond = 5.00
		}
		if config.LatThresh <= 0 {
			config.LatThresh = 5
		}
		if config.FailCond <= 0 {
			config.FailCond = 0.1
		}
		if config.FailThresh <= 0 {
			config.FailThresh = 5
		}
		// Default value for the time limit
		if config.TimeLimit <= 0 {
			config.TimeLimit = defaultTimeLimit
		}
		// Create a new time-based window
		window, err := NewTimeBasedWindow(config)
		if err != nil {
			return err
		}
		// Store the window in the monitor
		DefaultHttpMonitor.windows.Store(config.Path, window)
	}
	return nil
}

// NewTimeBasedWindow Function to create a new time-based window
func NewTimeBasedWindow(config MonitorConfig) (*slidingWindow, error) {
	// Check for valid parameters
	if config.LatCond <= 0 || config.FailCond <= 0 || config.LatThresh <= 0 || config.FailThresh <= 0 {
		return nil, errors.New("NewTimeBasedWindow invalid parameters")
	}
	base := &slidingWindow{
		path:                  config.Path,
		failedCondition:       config.FailCond,
		latencyCondition:      config.LatCond,
		latencyAlertThreshold: config.LatThresh,
		failedAlertThreshold:  config.FailThresh,
		latencyAlertOn:        config.AlertOnLatency,
		latencyCond:           config.LatencyCond,
		timeLimit:             config.TimeLimit,
		wheel: &timeWheel{
			data:  make([]request, WheelSize),
			index: 0,
		},
		alertFunc:         config.AlertFunc,
		alertCooldownTime: 0,
	}
	// Start the rotation of the wheel
	go base.rotateWheel()
	return base, nil
}

// Record Method to record a request in the monitor
func (m *HttpCodeMonitor) Record(endpoint string, success bool, latency float64) error {
	value, exists := m.windows.Load(endpoint)
	if !exists {
		return errors.New("endpoint not found endpoint:" + endpoint)
	}
	window, ok := value.(*slidingWindow)
	if !ok {
		return errors.New("type assertion failed endpoint:" + endpoint)
	}
	window.record(success, latency)
	return nil
}

// Method to update the count of latency exceed count
func (sw *slidingWindow) updateLatExceedCount(req request) {
	if !req.isExist {
		return
	}
	decrement := func(c *int64) {
		*c--
		if *c < 0 {
			*c = 0
		}
	}
	decrement(&sw.totalReq)
	if !req.success {
		decrement(&sw.failedReq)
	}
	sw.totalLatency -= req.latency
	if sw.totalLatency < 0 {
		sw.totalLatency = 0
	}
	if req.latency > req.oldLatency*sw.latencyCondition {
		decrement(&sw.latencyExceedCount)
	}
}

// Method to record a request in the window
func (sw *slidingWindow) record(success bool, latency float64) {
	sw.mux.Lock()
	defer sw.mux.Unlock()

	if !success {
		sw.failedReq++
	}
	sw.totalReq++
	sw.totalLatency += latency

	avgLat := sw.averageLatency()
	failRate := sw.failureRate()

	sw.wheel.push(request{
		success:    success,
		latency:    latency,
		oldLatency: avgLat,
		isExist:    true,
	})

	if sw.latencyAlertOn {
		sw.updateCounterAndAlert(
			latency > sw.latencyCond && latency > avgLat*sw.latencyCondition,
			&sw.latencyExceedCount,
			sw.latencyAlertThreshold,
		)
	}

	sw.updateCounterAndAlert(
		failRate > sw.failedCondition,
		&sw.failedExceedCount,
		sw.failedAlertThreshold,
	)

}

// Method to update the counter and potentially trigger an alert
func (sw *slidingWindow) updateCounterAndAlert(condition bool, count *int64, thresh int64) {
	if sw.totalReq < 100 {
		return
	}
	if condition {
		*count++
		if sw.alertFunc != nil && *count >= thresh {
			now := time.Now().Unix()
			if now-sw.alertCooldownTime >= 60*60 {
				sw.alertFunc(sw.path, sw.totalReq, sw.failureRate(), sw.averageLatency())
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

// Method to compute the failure rate
func (sw *slidingWindow) failureRate() float64 {
	if sw.totalReq == 0 {
		return 0
	}
	return float64(sw.failedReq) / float64(sw.totalReq)
}

// Method to compute the average latency
func (sw *slidingWindow) averageLatency() float64 {
	if sw.totalReq == 0 {
		return 0
	}
	return sw.totalLatency / float64(sw.totalReq)
}

// Structure for the sliding window
type slidingWindow struct {
	path  string       // Path
	wheel *timeWheel   // Time wheel
	mux   sync.RWMutex // Read-write lock

	timeLimit int64 // Time limit

	totalReq  int64 // Total number of requests
	failedReq int64 // Number of failed requests

	totalLatency float64 // Total latency

	latencyAlertOn        bool // Whether to alert on latency
	latencyCond           float64
	latencyExceedCount    int64   // Count of latency exceedances
	latencyAlertThreshold int64   // Latency alert threshold
	latencyCondition      float64 // Condition for latency exceedance

	failedExceedCount    int64   // Count of failures
	failedAlertThreshold int64   // Failure alert threshold
	failedCondition      float64 // Condition for failure

	alertFunc         func(path string, total int64, failRate float64, avg float64) // Alert function
	alertCooldownTime int64
}

// HttpCodeMonitor Structure for the HTTP code monitor
type HttpCodeMonitor struct {
	windows sync.Map
}

// Structure for a request
type request struct {
	success    bool
	latency    float64
	oldLatency float64
	isExist    bool
}

// Structure for a time wheel
type timeWheel struct {
	data  []request
	index int
}

// Method to push a request onto the time wheel
func (tw *timeWheel) push(req request) {
	tw.data[tw.index] = req
	tw.index = (tw.index + 1) % WheelSize
}

// Method to rotate the wheel
func (sw *slidingWindow) rotateWheel() {
	ticker := time.NewTicker(time.Second * time.Duration(sw.timeLimit) / WheelSize)
	defer ticker.Stop()
	for range ticker.C {
		sw.mux.Lock()
		oldReq := sw.wheel.data[sw.wheel.index]
		sw.updateLatExceedCount(oldReq)
		sw.wheel.push(request{})
		sw.mux.Unlock()
	}
}
