package ws

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Covsj/goTool/log"
	"github.com/gorilla/websocket"
)

const (
	pingPeriod = 60 * time.Second
)

// WebSocketPool 是一个 WebSocket 连接池
type WebSocketPool struct {
	URL     string
	Header  http.Header
	conn    *websocket.Conn
	connect bool
}

// NewWebSocketPool 创建一个新的 WebSocket 连接池
func NewWebSocketPool(url string, header http.Header) *WebSocketPool {
	return &WebSocketPool{
		URL:    url,
		Header: header,
	}
}

// Start 启动 WebSocket 连接池
func (pool *WebSocketPool) Start() chan []byte {
	messageChan := make(chan []byte)

	go func() {
		for {
			// 尝试建立连接
			conn, _, err := websocket.DefaultDialer.Dial(pool.URL, pool.Header)
			if err != nil {
				log.ErrorF("连接错误：%s", err)
				// 连接失败，等待一段时间后尝试重新连接
				time.Sleep(5 * time.Second)
				continue
			}
			log.InfoF("WebSocket 连接已建立")
			pool.conn = conn
			pool.connect = true

			// 启动 Ping 定时任务
			go pool.ping()

			// 循环接收消息
			for pool.connect {
				_, message, err := conn.ReadMessage()
				if err != nil {
					log.ErrorF("读取消息错误：%s", err)
					break
				}
				if fmt.Sprintf("%s", message) == "\"ping\"" {
					err := conn.WriteMessage(websocket.TextMessage, []byte("pong"))
					if err != nil {
						log.ErrorF("发送 pong 消息错误：%s", err)
						break
					}
					//log.InfoF("已发送 pong 消息")
				} else {
					messageChan <- message
				}
			}
			log.ErrorF("WebSocket 连接已关闭")
			// 关闭连接
			_ = conn.Close()
		}
	}()

	return messageChan
}

// ping 定时发送 ping 消息以保持连接
func (pool *WebSocketPool) ping() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		<-ticker.C
		if pool.conn != nil && pool.connect {
			err := pool.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.InfoF("发送 ping 消息错误：%s", err)
				break
			}
		}
	}
}
