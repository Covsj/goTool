package monitor

import "errors"

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		data: make([]Request, size),
	}
}

func (rb *RingBuffer) Empty() bool {
	rb.mux.RLock()
	defer rb.mux.RUnlock()
	return !rb.full && (rb.start == rb.end)
}

func (rb *RingBuffer) Full() bool {
	rb.mux.RLock()
	defer rb.mux.RUnlock()
	return rb.full
}

func (rb *RingBuffer) Peek() (Request, error) {
	rb.mux.RLock()
	defer rb.mux.RUnlock()
	if rb.Empty() {
		return Request{}, errors.New("ring buffer is empty")
	}
	return rb.data[rb.start], nil
}

func (rb *RingBuffer) Push(req Request) {
	rb.mux.Lock()
	defer rb.mux.Unlock()
	rb.data[rb.end] = req
	rb.end = (rb.end + 1) % len(rb.data)
	if rb.full {
		rb.start = (rb.start + 1) % len(rb.data)
	}
	rb.full = rb.full || rb.end == rb.start
}

func (rb *RingBuffer) Pop() (Request, bool) {
	rb.mux.Lock()
	defer rb.mux.Unlock()
	if rb.start == rb.end && !rb.full {
		return Request{}, false
	}
	req := rb.data[rb.start]
	rb.start = (rb.start + 1) % len(rb.data)
	rb.full = false
	return req, true
}
