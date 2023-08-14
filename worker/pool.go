package worker

import (
	"context"
	"sync"
	"sync/atomic"
)

type Pool struct {
	options *options

	mu   sync.Mutex
	once sync.Once

	workerNum      uint32
	taskNum        uint32
	waitingTaskNum uint32
	status         uint32

	workerC   chan func()
	taskC     chan func()
	waitingC  chan func()
	dismissC  chan struct{}
	pauseC    chan struct{}
	shutdownC chan struct{}
}

const (
	_ uint32 = iota
	statusInitialized
	statusPlaying
	statusCleaning
	statusShutdown
)

func NewPoolWhitOption(opts ...Option) *Pool {
	options := newOptions(opts...)
	v := &Pool{
		options:   options,
		workerC:   make(chan func()),
		taskC:     make(chan func()),
		waitingC:  make(chan func(), options.waitingQueueSize),
		dismissC:  make(chan struct{}),
		pauseC:    make(chan struct{}),
		shutdownC: make(chan struct{}),
	}
	_ = atomic.CompareAndSwapUint32(&v.status, 0, statusInitialized)
	go v.play()
	return v
}

// Submit a task to the worker pool
func (v *Pool) Submit(task func()) {
	v.submit(false, task)
}

// SubmitWait submit a task to the worker pool and wait for it to complete
func (v *Pool) SubmitWait(task func()) {
	v.submit(true, task)
}

// Consume the current tasks in the taskC
// Note: Consume will not execute the tasks which put into the taskC after calling Consume
func (v *Pool) Consume(taskC chan func()) {
	v.consume(false, taskC)
}

// ConsumeWait consume tasks in the taskC and wait for them to complete
// Note: ConsumeWait will not execute the tasks which put into the taskC after calling ConsumeWait
func (v *Pool) ConsumeWait(taskC chan func()) {
	v.consume(true, taskC)
}

// Pause the worker pool
func (v *Pool) Pause(ctx context.Context) {
	v.pause(ctx)
}

// Shutdown the worker pool
func (v *Pool) Shutdown() {
	v.shutdown(false)
}

// ShutdownWait graceful shutdown, wait for all tasks to complete
func (v *Pool) ShutdownWait() {
	v.shutdown(true)
}

// IsPlaying returns true if the worker pool is playing (running)
func (v *Pool) IsPlaying() bool {
	return atomic.LoadUint32(&v.status) == statusPlaying
}

// IsCleaning returns true if the worker pool is cleaning (graceful shutdown)
func (v *Pool) IsCleaning() bool {
	return atomic.LoadUint32(&v.status) == statusCleaning
}

// IsShutdown returns true if the worker pool is shutdown
func (v *Pool) IsShutdown() bool {
	return atomic.LoadUint32(&v.status) == statusShutdown
}

// MaxWorkerNum returns the maximum number of workers
func (v *Pool) MaxWorkerNum() int {
	return v.options.maxWorkers
}

// WaitingQueueSize returns the size of the waiting queue
func (v *Pool) WaitingQueueSize() int {
	return v.options.waitingQueueSize
}

// TaskNum returns the number of tasks
func (v *Pool) TaskNum() uint32 {
	return atomic.LoadUint32(&v.taskNum)
}

// WaitingTaskNum returns the number of waiting tasks
func (v *Pool) WaitingTaskNum() uint32 {
	return atomic.LoadUint32(&v.waitingTaskNum)
}

// WorkerNum returns the number of workers
func (v *Pool) WorkerNum() uint32 {
	return atomic.LoadUint32(&v.workerNum)
}
