package worker

import "time"

type Option func(*options)

type options struct {
	maxWorkers        int
	waitingQueueSize  int
	workerIdleTimeout time.Duration
}

var defaultOptions = options{
	maxWorkers:        5,
	waitingQueueSize:  64,
	workerIdleTimeout: 3 * time.Second,
}

func newOptions(opts ...Option) *options {
	options := &options{
		maxWorkers:        defaultOptions.maxWorkers,
		waitingQueueSize:  defaultOptions.waitingQueueSize,
		workerIdleTimeout: defaultOptions.workerIdleTimeout,
	}
	options.apply(opts...)
	return options
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithMaxWorkers set the maximum number of workers
func WithMaxWorkers(max int) Option {
	if max < 1 {
		max = 1
	}
	return func(o *options) {
		o.maxWorkers = max
	}
}

// WithWaitingQueueSize set the size of the waiting queue
func WithWaitingQueueSize(size int) Option {
	return func(o *options) {
		o.waitingQueueSize = size
	}
}

// WithWorkerIdleTimeout set the destroyed timeout of idle workers
func WithWorkerIdleTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.workerIdleTimeout = timeout
	}
}
