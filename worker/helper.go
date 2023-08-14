package worker

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

func (v *Pool) submit(wait bool, task func()) {
	if task == nil {
		return
	}
	if wait {
		ctx, done := context.WithCancel(context.Background())
		v.taskC <- func() {
			task()
			done()
		}
		_ = atomic.AddUint32(&v.taskNum, 1)
		<-ctx.Done()
	} else {
		v.taskC <- task
		_ = atomic.AddUint32(&v.taskNum, 1)
	}
}

func (v *Pool) consume(wait bool, taskC chan func()) {
	if taskC == nil {
		return
	}
	nums := len(taskC)
	if wait {
		wg := new(sync.WaitGroup)
		wg.Add(nums)
		for i := 0; i < nums; i++ {
			task := <-taskC
			v.taskC <- func() {
				task()
				wg.Done()
			}
			_ = atomic.AddUint32(&v.taskNum, 1)
		}
		wg.Wait()
	} else {
		for i := 0; i < nums; i++ {
			task := <-taskC
			v.taskC <- task
			_ = atomic.AddUint32(&v.taskNum, 1)
		}
	}
}

func (v *Pool) pause(ctx context.Context) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.IsShutdown() {
		return
	}
	wg := new(sync.WaitGroup)
	wg.Add(v.MaxWorkerNum())
	for i := 0; i < v.MaxWorkerNum(); i++ {
		v.Submit(func() {
			wg.Done()
			select {
			case <-ctx.Done():
			case <-v.pauseC:
			}
		})
	}
	wg.Wait()
}

func (v *Pool) shutdown(wait bool) {
	v.once.Do(func() {
		close(v.taskC)
		close(v.pauseC)
		if wait {
			_ = atomic.CompareAndSwapUint32(&v.status, statusPlaying, statusCleaning)
			v.clean()
			v.waitClose()
		} else {
			v.waitClose()
			_ = atomic.CompareAndSwapUint32(&v.status, statusPlaying, statusShutdown)
		}
	})
}

func (v *Pool) play() {
	defer func() {
		close(v.shutdownC)
		if v.IsCleaning() {
			_ = atomic.CompareAndSwapUint32(&v.status, statusCleaning, statusShutdown)
		}
	}()
	_ = atomic.CompareAndSwapUint32(&v.status, statusInitialized, statusPlaying)
	wg := new(sync.WaitGroup)
	timer := time.NewTimer(v.options.workerIdleTimeout)
	defer timer.Stop()
LOOP:
	for {
		if int(v.WorkerNum()) < v.MaxWorkerNum() {
			_ = atomic.AddUint32(&v.workerNum, 1)
			go v.recruit(wg)
		}
		select {
		case task, ok := <-v.waitingC:
			if !ok {
				break LOOP
			}
			select {
			case v.workerC <- task:
				_ = atomic.AddUint32(&v.waitingTaskNum, ^uint32(0))
			default:
				v.waitingC <- task
			}
		default:
			select {
			case task, ok := <-v.taskC:
				if !ok {
					break LOOP
				}
				select {
				case v.workerC <- task:
				default:
					v.waitingC <- task
					_ = atomic.AddUint32(&v.waitingTaskNum, 1)
				}
			case <-timer.C:
				if v.WorkerNum() > 0 {
					v.tryDismiss()
				}
				timer.Reset(v.options.workerIdleTimeout)
			}
		}
	}
	v.dismissAll()
	wg.Wait()
}

func (v *Pool) clean() {
	for v.WaitingTaskNum() > 0 {
		task, ok := <-v.waitingC
		if !ok {
			break
		}
		_ = atomic.AddUint32(&v.waitingTaskNum, ^uint32(0))
		v.workerC <- task
	}
}

func (v *Pool) recruit(wg *sync.WaitGroup) {
	wg.Add(1)
	defer func() {
		_ = atomic.AddUint32(&v.workerNum, ^uint32(0))
		wg.Done()
	}()
LOOP:
	for {
		if atomic.LoadUint32(&v.status) == statusShutdown {
			break
		}
		select {
		case <-v.dismissC:
			break LOOP
		case task, ok := <-v.workerC:
			if !ok {
				break LOOP
			}
			task()
			_ = atomic.AddUint32(&v.taskNum, ^uint32(0))
		}
	}
}

func (v *Pool) tryDismiss() {
	select {
	case v.dismissC <- struct{}{}:
	default:
	}
}

func (v *Pool) dismissAll() {
	for v.WorkerNum() > 0 {
		v.tryDismiss()
	}
	close(v.dismissC)
}

func (v *Pool) waitClose() {
	<-v.shutdownC
	close(v.waitingC)
	close(v.workerC)
}
