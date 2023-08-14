package monitor

import (
	"errors"
	"time"
)

// SlidingWindow represents a sliding window for monitoring HTTP requests.
type SlidingWindow struct {
	// The path of the HTTP endpoint.
	Path string

	// The time wheel.
	Wheel *TimeWheel

	// The channel for sending requests to the processing goroutine.
	RequestChannel chan Request

	// The size for request channel.
	RequestChannelSize int64

	// The time limit for the sliding window.
	TimeLimit int64

	// The size of the wheel.
	WheelSize int64

	// The total number of requests.
	TotalRequests int64

	// The number of failed requests.
	FailedRequests int64

	// The total latency.
	TotalLatency float64

	// Whether to alert on latency.
	AlertOnLatency bool

	// Latency threshold.
	LatencyThreshold float64

	// The count of latency exceed latencyThreshold.
	LatencyExceedCount int64

	// The latency alert threshold.
	LatencyAlertThreshold int64

	// The multiplier for the latency threshold.
	LatencyMultiplier float64

	// The count of failures.
	FailureExceedCount int64

	// The failure alert threshold.
	FailureAlertThreshold int64

	// The failure rate threshold.
	FailureThreshold float64

	// The minimum number of requests needed to perform alert checks.
	RequestThreshold int64

	// The function to call when an alert is triggered.
	AlertFunc func(path string, total int64, failRate float64, avg float64)

	// The alert cooldown time.
	AlertCooldownTime int64
}

// creates a new time-based window with the given configuration.
func newTimeBasedWindow(config *Config) (*SlidingWindow, error) {
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
		Path:                  config.Path,
		FailureThreshold:      config.FailureThreshold,
		LatencyMultiplier:     config.LatencyMultiplier,
		LatencyAlertThreshold: config.LatencyCount,
		FailureAlertThreshold: config.FailureCount,
		AlertOnLatency:        config.AlertOnLatency,
		LatencyThreshold:      config.LatencyThreshold,
		TimeLimit:             config.TimeLimit,
		WheelSize:             config.WheelSize,
		Wheel: &TimeWheel{
			Data:  make([]Request, config.WheelSize),
			Index: 0,
		},
		RequestThreshold: config.RequestThreshold,
		AlertFunc:        config.AlertFunc,
		RequestChannel:   make(chan Request, config.RequestChannelSize),
	}
	go base.processRequests()
	go base.rotateWheel()
	return base, nil
}

// updates the count of latency exceed count.
func (sw *SlidingWindow) updateLatencyExceedCount(req Request) {
	if !req.IsExist {
		return
	}
	if sw.TotalRequests = sw.TotalRequests - 1; sw.TotalRequests < 0 {
		sw.TotalRequests = 0
	}
	if !req.Status {
		if sw.FailedRequests = sw.FailedRequests - 1; sw.FailedRequests < 0 {
			sw.FailedRequests = 0
		}
	}
	if sw.TotalLatency = sw.TotalLatency - req.Latency; sw.TotalLatency < 0 {
		sw.TotalLatency = 0
	}

	if sw.isNeedExceedLatencyCount(req.Latency, req.OldAvgLatency) {
		if sw.LatencyExceedCount = sw.LatencyExceedCount - 1; sw.LatencyExceedCount < 0 {
			sw.LatencyExceedCount = 0
		}
	}
}

func (sw *SlidingWindow) isNeedExceedLatencyCount(latency, avg float64) bool {
	if latency > sw.LatencyThreshold && latency > avg*sw.LatencyMultiplier {
		return true
	}
	return false
}

// processes requests sent through the request channel.
func (sw *SlidingWindow) processRequests() {
	for req := range sw.RequestChannel {
		if req.IsExist {
			sw.doRecord(req.Status, req.Latency)
		}
	}
}

// doRecord records a request in the window.
func (sw *SlidingWindow) doRecord(status bool, latency float64) {
	if !status {
		sw.FailedRequests++
	}
	sw.TotalRequests++
	sw.TotalLatency += latency
	avgLat := sw.averageLatency()
	failRate := sw.failureRate()
	sw.Wheel.push(Request{
		Status:        status,
		Latency:       latency,
		OldAvgLatency: avgLat,
		IsExist:       true,
	})
	if sw.AlertOnLatency {
		sw.updateCounterAndAlert(
			latency > sw.LatencyThreshold && latency > avgLat*sw.LatencyMultiplier,
			&sw.LatencyExceedCount,
			sw.LatencyAlertThreshold,
		)
	}
	sw.updateCounterAndAlert(
		failRate > sw.FailureThreshold,
		&sw.FailureExceedCount,
		sw.FailureAlertThreshold,
	)
}

// record records a request in the window.
func (sw *SlidingWindow) record(status bool, latency float64) {
	sw.RequestChannel <- Request{Status: status, Latency: latency, IsExist: true}
}

// updates the counter and potentially triggers an alert.
func (sw *SlidingWindow) updateCounterAndAlert(condition bool, count *int64, threshold int64) {
	if sw.RequestThreshold > 0 && sw.TotalRequests < sw.RequestThreshold {
		return
	}

	if condition {
		*count++
		if sw.AlertFunc != nil && *count >= threshold {
			now := time.Now().Unix()
			if now-sw.AlertCooldownTime >= 60*60 {
				sw.AlertFunc(sw.Path, sw.TotalRequests, sw.failureRate(), sw.averageLatency())
				sw.AlertCooldownTime = now
			}
			*count = 0
		}
	} else if *count > 0 {
		*count--
	} else {
		*count = 0
	}
}

// computes the failure rate.
func (sw *SlidingWindow) failureRate() float64 {
	if sw.TotalRequests == 0 {
		return 0
	}
	return float64(sw.FailedRequests) / float64(sw.TotalRequests)
}

// computes the average latency.
func (sw *SlidingWindow) averageLatency() float64 {
	if sw.TotalRequests == 0 {
		return 0
	}
	return sw.TotalLatency / float64(sw.TotalRequests)
}

// rotateWheel rotates the wheel.
func (sw *SlidingWindow) rotateWheel() {
	ticker := time.NewTicker(time.Duration(sw.TimeLimit) * time.Second / time.Duration(sw.WheelSize))
	defer ticker.Stop()
	for range ticker.C {
		oldReq := sw.Wheel.Data[sw.Wheel.Index]
		sw.updateLatencyExceedCount(oldReq)
		sw.Wheel.push(Request{})
	}
}
