package monitor

// Config  contains the configuration for a monitor.
type Config struct {
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

// Request represents a request.
type Request struct {
	// Whether the request was successful.
	Status bool

	// The latency of the request.
	Latency float64

	// The previous average latency.
	OldAvgLatency float64

	// Whether the request exists.
	IsExist bool
}
