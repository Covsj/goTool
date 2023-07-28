package monitor

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestBatch(t *testing.T) {
	monitor := NewHttpCodeMonitor(WindowTimeBased)
	for i := 0; i < 10; i++ {
		go func() {
			for {
				time.Sleep(100 * time.Millisecond)
				endpoint := fmt.Sprintf("endpoint%d", rand.Intn(5))
				monitor.Record(endpoint, rand.Float64() > 0.2, rand.Float64())
			}
		}()
	}
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		for i := 0; i < 5; i++ {
			endpoint := fmt.Sprintf("endpoint%d", i)
			rate, ok := monitor.FailureRate(endpoint)
			if ok {
				fmt.Println(endpoint, "failure rate:", rate)
			} else {
				fmt.Println(endpoint, ok)
			}
			//latency, ok := monitor.AverageLatency(endpoint)
			//if ok {
			//	fmt.Println(endpoint, "avg latency:", latency)
			//} else {
			//	fmt.Println(endpoint, ok)
			//}
		}
	}
}
