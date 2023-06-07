package myHttp

import (
	"net/http"
	"sync"
	"time"

	"github.com/abursavich/nett"
)

var cp *ClientPool

type ClientPool struct {
	pool     *sync.Pool
	MaxConns int
	MaxIdle  int
	mu       sync.Mutex
	count    int
}

func init() {
	newClientPool()
}

func newClientPool() {
	cp = &ClientPool{
		pool: &sync.Pool{
			New: func() interface{} {
				dialer := &nett.Dialer{
					Resolver: &nett.CacheResolver{TTL: 5 * time.Minute},
					IPFilter: nett.DualStack,
					Timeout:  2 * time.Second,
				}
				client := &http.Client{
					Transport: &http.Transport{
						Dial:                dialer.Dial,
						MaxIdleConnsPerHost: 1000,
						IdleConnTimeout:     time.Second * 10,
					},
					Timeout: time.Second * 10,
				}
				return client
			},
		},
		MaxConns: 2000,
		MaxIdle:  100,
		count:    0,
	}
}

func GetClient() *http.Client {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if cp.pool == nil {
		return nil
	}

	client := cp.pool.Get().(*http.Client)
	cp.count++
	return client
}

func ReleaseClient(client *http.Client) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if cp.pool == nil {
		return
	}

	if cp.count >= cp.MaxIdle {
		client.CloseIdleConnections()
		cp.count--
		return
	}

	cp.pool.Put(client)
	cp.count--
}
