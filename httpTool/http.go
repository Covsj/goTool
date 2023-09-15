package httpTool

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultUserAgent   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"
	DefaultContentType = "application/json"
)

var defaultClient = &http.Client{
	Timeout:   10 * time.Second,
	Transport: http.DefaultTransport,
}

type RequestOptions struct {
	URL     string
	Method  string
	Body    string
	Headers map[string]string
}

func NewRequest(opts RequestOptions) (*http.Request, error) {
	if opts.Method == "" {
		opts.Method = http.MethodGet
	}
	req, err := http.NewRequest(opts.Method, opts.URL, bytes.NewBufferString(opts.Body))
	if err != nil {
		return nil, err
	}
	for k, v := range opts.Headers {
		req.Header.Set(k, v)
	}
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", DefaultUserAgent)
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", DefaultContentType)
	}
	return req, nil
}

func Execute(req *http.Request) (*http.Response, error) {
	return defaultClient.Do(req)
}

func DecodeBody(resp *http.Response) ([]byte, error) {
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()
	if strings.Contains(strings.ToLower(resp.Header.Get("Content-Encoding")), "gzip") {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		return ioutil.ReadAll(reader)
	}
	return ioutil.ReadAll(resp.Body)
}

func Send(opts RequestOptions) (*http.Response, []byte, error) {
	req, err := NewRequest(opts)
	if err != nil {
		return nil, nil, err
	}
	response, err := Execute(req)
	if err != nil {
		return nil, nil, err
	}
	body, err := DecodeBody(response)
	if err != nil {
		return response, nil, err
	}
	return response, body, nil
}
