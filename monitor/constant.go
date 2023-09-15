package monitor

// Default configurations for HTTPCodeMonitor.
const (
	DefaultLatencyMultiplier = 3
	DefaultLatencyCount      = 10

	DefaultFailureThreshold = 0.5
	DefaultFailureCount     = 10

	DefaultWheelSize          = 600
	DefaultTimeLimit          = 300
	DefaultRequestChannelSize = 10000
	DefaultRequestThreshold   = 100

	DefaultLatencyMax       = 8000
	DefaultCoolDown         = 3600
	DefaultLatencyThreshold = 3000
)
