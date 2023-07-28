package monitor

import (
	"sync"
)

const (
	WindowSizeBased = "size"
	WindowTimeBased = "time"
)

const (
	defaultSizeLimit  int64 = 2000
	defaultTimeLimit  int64 = 5 * 60
	defaultBufferSize       = 30000
)

type Request struct {
	timestamp  int64
	success    bool
	latency    float64
	oldLatency float64
	isExist    bool
}

type RingBuffer struct {
	data  []Request
	start int
	end   int
	full  bool
	mux   sync.RWMutex
}

type SlidingWindow struct {
	path       string
	windowType string
	//data       *list.List
	buffer *RingBuffer
	mux    sync.RWMutex

	sizeLimit int64
	timeLimit int64

	totalRequests  int64
	failedRequests int64

	totalLatency float64

	latencyTurnOnAlert            bool
	exceedThresholdCountByLatency int64   // 超过延迟平均延迟计数
	alertThresholdByLatency       int64   // 报警延迟计数阈值
	conditionByLatency            float64 // 超过平均延迟多少倍

	exceedThresholdCountByFailed int64   // 失败计数
	alertThresholdByFailed       int64   // 报警失败阈值
	conditionByFailed            float64 // 错误失败比例

	alertFunc func(path string, total any, failRate any, avg any)
}

type HttpCodeMonitor struct {
	defaultType string
	windows     map[string]*SlidingWindow
	mux         sync.Mutex
}
