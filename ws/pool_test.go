package ws

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNewPool(t *testing.T) {
	pool := NewWebSocketPool("wss://zws5-1.web.telegram.org/apiws",
		http.Header(map[string][]string{
			//"Sec-Websocket-Key": {"aNxM+aTf7x/mC9fiFFHP1A=="},
			//"Connection": {"Upgrade"},
			//"Upgrade":    {"websocket"},
			"Sec-WebSocket-Protocol": {"binary"},
			"User-Agent":             {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"},
			//"Host":       {"zws5-1.web.telegram.org"}}
		}))
	ch := pool.Start()
	for data := range ch {
		fmt.Println(string(data))
	}
}
