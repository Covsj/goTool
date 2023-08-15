package monitor

import (
	"errors"
	"math/rand"
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

	// The alert cool down timestamp.
	AlertCoolDownTimeByFailed  int64
	AlertCoolDownTimeByLatency int64

	LastRotationTime time.Time
	SampleRate       int64
	RequestCount     int64
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
	now := time.Now()
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
		LastRotationTime: now,
		SampleRate:       10,
		RequestCount:     0,
	}
	rand.Seed(now.Unix())
	go base.processAndRotate()
	return base, nil
}

// updates the count of latency exceed count.
func (sw *SlidingWindow) updateLatencyExceedCount(req Request) {
	if !req.IsValid {
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

// doRecord records a request in the window.
func (sw *SlidingWindow) doRecord(status bool, latency float64, IsValid bool) {
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
		IsValid:       IsValid,
	})

	totalRequests := sw.TotalRequests
	now := time.Now().Unix()

	if totalRequests >= sw.RequestThreshold {
		if sw.AlertOnLatency && sw.isNeedExceedLatencyCount(latency, avgLat) {
			sw.LatencyExceedCount++
			if sw.AlertFunc != nil && sw.LatencyExceedCount >= sw.LatencyAlertThreshold {
				if now-sw.AlertCoolDownTimeByLatency >= defaultCoolDown {
					sw.AlertFunc(sw.Path, totalRequests, failRate, avgLat)
					sw.AlertCoolDownTimeByLatency = now
				}
				sw.LatencyExceedCount = 0
			}
		} else if sw.LatencyExceedCount > 0 {
			sw.LatencyExceedCount--
		} else {
			sw.LatencyExceedCount = 0
		}
		if failRate > sw.FailureThreshold {
			sw.FailureExceedCount++
			if sw.AlertFunc != nil && sw.FailureExceedCount >= sw.FailureAlertThreshold {
				if now-sw.AlertCoolDownTimeByFailed >= defaultCoolDown {
					sw.AlertFunc(sw.Path, totalRequests, failRate, avgLat)
					sw.AlertCoolDownTimeByFailed = now
				}
				sw.FailureExceedCount = 0
			}
		} else if sw.FailureExceedCount > 0 {
			sw.FailureExceedCount--
		} else {
			sw.FailureExceedCount = 0
		}
	}
}

// record records a request in the window.
func (sw *SlidingWindow) record(status bool, latency float64) {
	// Randomly drop 20% of the requests
	if rand.Float32() > 0.8 {
		return
	}
	req := Request{Status: status, Latency: latency, IsValid: true}
	// Try to send to channel, if not possible, just return
	select {
	case sw.RequestChannel <- req:
	default:
		return
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

func (sw *SlidingWindow) processAndRotate() {
	rotationInterval := time.Duration(sw.TimeLimit) * time.Second / time.Duration(sw.WheelSize)
	for req := range sw.RequestChannel {
		sw.RequestCount++
		if req.IsValid && sw.RequestCount%sw.SampleRate == 0 {
			sw.doRecord(req.Status, req.Latency, req.IsValid)
		}
		// Reset requestCount if it becomes too large
		if sw.RequestCount >= sw.SampleRate*100 {
			sw.RequestCount = 0
		}
		// Check if it's time to rotate the wheel
		if time.Since(sw.LastRotationTime) >= rotationInterval {
			oldReq := sw.Wheel.Data[sw.Wheel.Index]
			sw.updateLatencyExceedCount(oldReq)
			sw.Wheel.push(Request{IsValid: false})
			sw.LastRotationTime = time.Now()
		}
	}
}
