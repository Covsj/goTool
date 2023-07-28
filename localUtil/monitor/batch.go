package monitor

func NewHttpCodeMonitor(defaultType string) *HttpCodeMonitor {
	if defaultType == "" {
		defaultType = WindowTimeBased
	}
	return &HttpCodeMonitor{
		defaultType: defaultType,
		windows:     make(map[string]*SlidingWindow),
	}
}

func (m *HttpCodeMonitor) Record(endpoint string, success bool, latency float64) {
	m.mux.Lock()
	window, exists := m.windows[endpoint]
	if !exists {
		if m.defaultType == WindowSizeBased {
			window = NewSlidingWindowBySize(endpoint, defaultSizeLimit, 5,
				0.1, 5, 5)
		} else {
			window = NewSlidingWindowByTime(endpoint, defaultTimeLimit, 5,
				0.1, 5, 5)
		}
		m.windows[endpoint] = window
	}
	m.mux.Unlock()
	window.Record(success, latency)
}

func (m *HttpCodeMonitor) FailureRate(endpoint string) (float64, bool) {
	m.mux.Lock()
	window, exists := m.windows[endpoint]
	m.mux.Unlock()
	if !exists {
		return 0, false
	}
	return window.FailureRate(), true
}

func (m *HttpCodeMonitor) AverageLatency(endpoint string) (float64, bool) {
	m.mux.Lock()
	window, exists := m.windows[endpoint]
	m.mux.Unlock()
	if !exists {
		return 0, false
	}
	return window.AverageLatency(), true
}
