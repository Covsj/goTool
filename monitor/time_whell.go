package monitor

// TimeWheel represents a time wheel.
type TimeWheel struct {
	// The data in the time wheel.
	Data []Request

	// The current index in the time wheel.
	Index int
}

// pushes a request onto the time wheel.
func (tw *TimeWheel) push(req Request) {
	tw.Data[tw.Index] = req
	tw.Index = (tw.Index + 1) % len(tw.Data)
}
