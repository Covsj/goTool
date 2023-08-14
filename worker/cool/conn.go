package cool

import (
	"net"
	"sync"
	"time"
)

var _ net.Conn = (*Conn)(nil)

type Conn struct {
	net.Conn
	mu       sync.RWMutex
	unusable bool
	c        *cool
	t        time.Time
}

// Close overrides the net.Conn Close method
// put the usable connection back to the pool instead of closing it
func (cc *Conn) Close() error {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	if cc.unusable {
		if cc.Conn != nil {
			return cc.Conn.Close()
		}
		return nil
	}
	return cc.c.put(cc.Conn)
}

func (cc *Conn) MarkUnusable() {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.unusable = true
}

func (cc *Conn) IsUnusable() bool {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	return cc.unusable
}
