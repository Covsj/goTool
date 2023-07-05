package localUtil

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

var cp *ClientPool
var DefaultUa = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"
var DefaultType = "application/json"

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

func NewHttpRequest(url, method, body string, header map[string]string) (*http.Request, error) {
	if method == "" {
		method = http.MethodPost
	}
	req, err := http.NewRequest(method, url, ioutil.NopCloser(bytes.NewBuffer([]byte(body))))
	req.ContentLength = int64(len(body))
	if len(header) > 0 {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	} else {
		req.Header.Set("User-Agent", DefaultUa)
		req.Header.Set("Content-Type", DefaultType)
	}
	return req, err
}

func DoHttpRequest(req *http.Request, client *http.Client) (*http.Response, error) {
	if client == nil {
		client = GetClient()
		defer ReleaseClient(client)
	}

	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	if resp.Body != nil {
		bodyRaw, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close() // copy body to release socket
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyRaw))
	}
	return resp, err
}

func ReadRespBody(resp *http.Response) (body []byte, err error) {
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()
	if strings.Index(strings.ToLower(resp.Header.Get("Content-Type")), "gzip") >= 0 {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			return body, err
		}
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return body, err
	}
	return
}

func CallHttp(url, method, body string, header map[string]string) (*http.Response, []byte, error) {
	req, err := NewHttpRequest(url, method, body, header)
	if err != nil {
		return nil, nil, err
	}
	response, err := DoHttpRequest(req, nil)
	if response == nil || err != nil {
		return nil, nil, err
	}
	res, err := ReadRespBody(response)
	if err != nil {
		return nil, nil, err
	}
	return response, res, err
}

func newClientPool() {
	cp = &ClientPool{
		pool: &sync.Pool{
			New: func() interface{} {
				client := &http.Client{
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
