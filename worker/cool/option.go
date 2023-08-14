package cool

import "time"

type Option func(*options)

type options struct {
	connIdleTimeout time.Duration
}

func newOptions(opts ...Option) *options {
	options := &options{}
	options.apply(opts...)
	return options
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithConnIdleTimeout will set the connection idle timeout
func WithConnIdleTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.connIdleTimeout = timeout
	}
}
