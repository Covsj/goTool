package monitor

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

var defaultHTTPMonitor *HTTPCodeMonitor

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
			config.LatencyMultiplier = DefaultLatencyMultiplier
		}
		if config.LatencyCount <= 0 {
			config.LatencyCount = DefaultLatencyCount
		}
		if config.FailureThreshold <= 0 {
			config.FailureThreshold = DefaultFailureThreshold
		}
		if config.FailureCount <= 0 {
			config.FailureCount = DefaultFailureCount
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
		if config.RequestThreshold <= 0 {
			config.RequestThreshold = DefaultRequestThreshold
		}
		//if config.AlertFunc == nil {
		//	config.AlertFunc = AlertHttp
		//}
		if config.LatencyThreshold <= 0 {
			config.LatencyThreshold = DefaultLatencyThreshold
		}
		window, err := newTimeBasedWindow(config)
		if err != nil {
			return err
		}
		defaultHTTPMonitor.windows.Store(config.Path, window)
	}
	return nil
}

// creates a new time-based window with the given configuration.
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
	now := time.Now()
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
		lastRotationTime: now,
		sampleRate:       10, // 设置采样率为10，你可以根据实际需要进行调整
		requestCount:     0,
	}
	rand.Seed(now.Unix())
	go base.processAndRotate()
	return base, nil
}

// HTTPCodeMonitor get endpoint status
func (m *HTTPCodeMonitor) HTTPCodeMonitor(endpoint string) (error, int64, int64, float64, int64, int64) {
	value, exists := m.windows.Load(endpoint)
	if !exists {
		return errors.New("endpoint not found: " + endpoint), 0, 0, 0, 0, 0
	}
	sw, ok := value.(*SlidingWindow)
	if !ok {
		return errors.New("type assertion failed for endpoint: " + endpoint), 0, 0, 0, 0, 0
	}
	return nil, sw.totalRequests, sw.failedRequests, sw.totalLatency, sw.failureExceedCount, sw.latencyExceedCount
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

// records a request in the window.
func (sw *SlidingWindow) record(success bool, latency float64) {
	// Randomly drop 20% of the requests
	if rand.Float32() > 0.8 {
		return
	}
	req := Request{success: success, latency: latency, isValid: true}
	// Try to send to channel, if not possible, just return
	select {
	case sw.requestChannel <- req:
	default:
		return
	}
}

// update the count of latency exceed count.
func (sw *SlidingWindow) updateLatencyExceedCount(req Request) {
	if !req.isValid {
		return
	}
	if sw.totalRequests = sw.totalRequests - 1; sw.totalRequests < 0 {
		sw.totalRequests = 0
	}
	if !req.success {
		if sw.failedRequests = sw.failedRequests - 1; sw.failedRequests < 0 {
			sw.failedRequests = 0
		}
	}
	if sw.totalLatency = sw.totalLatency - req.latency; sw.totalLatency < 0 {
		sw.totalLatency = 0
	}

	if sw.isNeedExceedLatencyCount(req.latency, req.oldAvgLatency) {
		if sw.latencyExceedCount = sw.latencyExceedCount - 1; sw.latencyExceedCount < 0 {
			sw.latencyExceedCount = 0
		}
	}
}

func (sw *SlidingWindow) isNeedExceedLatencyCount(latency, avg float64) bool {
	if (latency > sw.latencyThreshold && latency > avg*sw.latencyMultiplier) || latency > DefaultLatencyMax {
		return true
	}
	return false
}

// records a request in the window.
func (sw *SlidingWindow) doRecord(success bool, latency float64, isValid bool) {
	if !success {
		sw.failedRequests++
	}
	sw.totalRequests++
	sw.totalLatency += latency
	avgLat := sw.averageLatency()
	failRate := sw.failureRate()
	sw.push(Request{
		success:       success,
		latency:       latency,
		oldAvgLatency: avgLat,
		isValid:       isValid,
	})

	totalRequests := sw.totalRequests
	now := time.Now().Unix()

	if totalRequests >= sw.requestThreshold {
		sw.checkLatency(latency, avgLat, now, totalRequests, failRate)
		sw.checkFailure(failRate, avgLat, now, totalRequests)
	}
}

func (sw *SlidingWindow) checkLatency(latency, avgLat float64, now int64, totalRequests int64, failRate float64) {
	if sw.alertOnLatency && sw.isNeedExceedLatencyCount(latency, avgLat) {
		sw.latencyExceedCount++
		if sw.alertFunc != nil && sw.latencyExceedCount >= sw.latencyAlertThreshold {
			if now-sw.alertCoolDownTimeByLatency >= DefaultCoolDown {
				sw.alertFunc(sw.path, totalRequests, failRate, avgLat)
				sw.alertCoolDownTimeByLatency = now
			}
			sw.latencyExceedCount = 0
		}
	} else if sw.latencyExceedCount > 0 {
		sw.latencyExceedCount--
	} else {
		sw.latencyExceedCount = 0
	}
}

func (sw *SlidingWindow) checkFailure(failRate, avgLat float64, now int64, totalRequests int64) {
	if failRate >= sw.failureThreshold {
		sw.failureExceedCount++
		if sw.alertFunc != nil && sw.failureExceedCount >= sw.failureAlertThreshold {
			if now-sw.alertCoolDownTimeByFailed >= DefaultCoolDown {
				sw.alertFunc(sw.path, totalRequests, failRate, avgLat)
				sw.alertCoolDownTimeByFailed = now
			}
			sw.failureExceedCount = 0
		}
	} else if sw.failureExceedCount > 0 {
		sw.failureExceedCount--
	} else {
		sw.failureExceedCount = 0
	}
}

// computes the failure rate.
func (sw *SlidingWindow) failureRate() float64 {
	if sw.totalRequests == 0 {
		return 0
	}
	return float64(sw.failedRequests) / float64(sw.totalRequests)
}

// computes the average latency.
func (sw *SlidingWindow) averageLatency() float64 {
	if sw.totalRequests == 0 {
		return 0
	}
	return sw.totalLatency / float64(sw.totalRequests)
}

// pushes a request onto the time wheel.
func (sw *SlidingWindow) push(req Request) {
	oldReq := sw.wheel.data[sw.wheel.index]
	sw.updateLatencyExceedCount(oldReq)
	sw.wheel.data[sw.wheel.index] = req
	sw.wheel.index = (sw.wheel.index + 1) % len(sw.wheel.data)
}

func (sw *SlidingWindow) processAndRotate() {
	for req := range sw.requestChannel {
		sw.requestCount++
		if req.isValid && sw.requestCount%sw.sampleRate == 0 {
			sw.doRecord(req.success, req.latency, req.isValid)
		}
		// Reset requestCount if it becomes too large
		if sw.requestCount >= sw.sampleRate*100 {
			sw.requestCount = 0
			//if rand.Intn(10) == 1 {
			//	log.Info("key", "http monitor processAndRotate",
			//		"path", sw.path,
			//		"total", sw.totalRequests,
			//		"failed", sw.failedRequests,
			//		"failed_rate", sw.failureRate(),
			//		"avg", sw.averageLatency(),
			//		"failure_exceed_count", sw.failureExceedCount,
			//		"latency_exceed_count", sw.latencyExceedCount)
			//}
		}
	}
}
