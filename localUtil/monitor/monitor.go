package monitor

import (
	"sync/atomic"
	"time"
)

func processNewParam(conditionByLatency, conditionByFailed *float64,
	alertThresholdByLatency, alertThresholdByFailed *int64) {
	if *conditionByLatency <= 0 {
		*conditionByLatency = 5.00
	}
	if *conditionByFailed <= 0 {
		*conditionByLatency = 0.1
	}
	if *alertThresholdByLatency <= 0 {
		*alertThresholdByFailed = 5
	}
	if *alertThresholdByFailed <= 0 {
		*alertThresholdByFailed = 5
	}
}
func newBaseWindow(path string, conditionByLatency, conditionByFailed float64,
	alertThresholdByLatency, alertThresholdByFailed int64) *SlidingWindow {
	processNewParam(&conditionByLatency, &conditionByFailed, &alertThresholdByLatency, &alertThresholdByFailed)
	return &SlidingWindow{
		path: path,
		//data:                    list.New(),
		buffer:                  NewRingBuffer(defaultBufferSize),
		conditionByFailed:       conditionByFailed,
		conditionByLatency:      conditionByLatency,
		alertThresholdByLatency: alertThresholdByLatency,
		alertThresholdByFailed:  alertThresholdByFailed,
	}
}

func NewSlidingWindowBySize(path string, sizeLimit int64,
	conditionByLatency, conditionByFailed float64,
	alertThresholdByLatency, alertThresholdByFailed int64) *SlidingWindow {

	base := newBaseWindow(path, conditionByLatency, conditionByFailed, alertThresholdByLatency, alertThresholdByFailed)

	if sizeLimit <= 0 {
		sizeLimit = defaultSizeLimit
	}
	base.windowType = WindowSizeBased
	base.sizeLimit = sizeLimit

	return base
}

func NewSlidingWindowByTime(path string, timeLimit int64, conditionByLatency, conditionByFailed float64,
	alertThresholdByLatency, alertThresholdByFailed int64) *SlidingWindow {

	base := newBaseWindow(path, conditionByLatency, conditionByFailed, alertThresholdByLatency, alertThresholdByFailed)
	if timeLimit <= 0 {
		timeLimit = defaultTimeLimit
	}
	base.windowType = WindowTimeBased
	base.timeLimit = timeLimit
	return base
}

func decrementCount(count *int64) {
	*count--
	if *count < 0 {
		*count = 0
	}
}

func (sw *SlidingWindow) updateExceedThresholdCountByLatency(req Request) {
	if !req.isExist {
		return
	}
	decrementCount(&sw.totalRequests)
	if !req.success {
		decrementCount(&sw.failedRequests)
	}
	sw.totalLatency -= req.latency
	if sw.totalLatency < 0 {
		sw.totalLatency = 0
	}
	if req.latency > req.oldLatency*sw.conditionByLatency {
		decrementCount(&sw.exceedThresholdCountByLatency)
	}
}

func (sw *SlidingWindow) Record(success bool, latency float64) {
	sw.mux.Lock()
	defer sw.mux.Unlock()
	now := time.Now().Unix()
	if sw.windowType == WindowTimeBased {
		for !sw.buffer.Empty() {
			req, err := sw.buffer.Peek()
			if err != nil {
				break
			}
			if now-req.timestamp > sw.timeLimit {
				req, ok := sw.buffer.Pop()
				if !ok {
					break
				}
				sw.updateExceedThresholdCountByLatency(req)
			} else {
				break
			}
		}
	}
	if sw.windowType == WindowSizeBased && sw.buffer.Full() {
		if req, ok := sw.buffer.Pop(); ok {
			sw.updateExceedThresholdCountByLatency(req)
		}
	}

	if !success {
		atomic.AddInt64(&sw.failedRequests, 1)
	}
	atomic.AddInt64(&sw.totalRequests, 1)
	sw.totalLatency += latency

	avgLatency := sw.averageLatency()
	failedRate := sw.failureRate()

	sw.buffer.Push(Request{
		timestamp:  now,
		success:    success,
		latency:    latency,
		oldLatency: avgLatency,
		isExist:    true,
	})
	if sw.totalRequests > 100 {
		if sw.latencyTurnOnAlert {
			sw.updateCounterAndAlert(
				latency > avgLatency*sw.conditionByLatency,
				&sw.exceedThresholdCountByLatency,
				sw.alertThresholdByLatency,
				sw.alertFunc,
			)
		}

		sw.updateCounterAndAlert(
			failedRate > sw.conditionByFailed,
			&sw.exceedThresholdCountByFailed,
			sw.alertThresholdByFailed,
			sw.alertFunc,
		)
	}

}

func (sw *SlidingWindow) updateCounterAndAlert(condition bool, counter *int64, threshold int64, alertFunc func(string, any, any, any)) {
	if condition {
		*counter++
	} else if *counter > 0 {
		*counter--
	}

	if *counter >= threshold {
		if alertFunc != nil {
			alertFunc(sw.path, sw.totalRequests, sw.failureRate(), sw.averageLatency())
		}
		*counter = 0
	}
}
func (sw *SlidingWindow) failureRate() float64 {
	if sw.totalRequests == 0 {
		return 0
	}
	return float64(sw.failedRequests) / float64(sw.totalRequests)
}

func (sw *SlidingWindow) FailureRate() float64 {
	sw.mux.RLock()
	defer sw.mux.RUnlock()
	return sw.failureRate()
}

func (sw *SlidingWindow) AverageLatency() float64 {
	sw.mux.RLock()
	defer sw.mux.RUnlock()
	return sw.averageLatency()
}

func (sw *SlidingWindow) averageLatency() float64 {
	if sw.totalRequests == 0 {
		return 0
	}
	return sw.totalLatency / float64(sw.totalRequests)
}
