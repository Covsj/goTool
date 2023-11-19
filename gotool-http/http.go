package gotool_http

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func init() {
	defaultClient = &Client{
		httpClient: &http.Client{},
	}
}

func NewClient(httpClient *http.Client) *Client {
	defaultClient = &Client{
		httpClient: httpClient,
	}
	return defaultClient
}

func Get(url string, out interface{}) (*http.Response, []byte, error) {
	opts := &RequestOptions{
		URL:         url,
		Method:      http.MethodGet,
		ResponseOut: out,
	}
	return Send(opts)
}

func Post(url string, body interface{}, out interface{}) (*http.Response, []byte, error) {
	opts := &RequestOptions{
		URL:         url,
		Method:      http.MethodPost,
		Body:        body,
		ResponseOut: out,
	}
	return Send(opts)
}

func SendWithRetries(opts *RequestOptions) (*http.Response, []byte, error) {
	if opts.Retries == 0 {
		opts.Retries = DefaultRetries
	}
	var resp *http.Response
	var body []byte
	var err error
	for i := 0; i < opts.Retries; i++ {
		resp, body, err = Send(opts)
		if err == nil {
			return resp, body, nil
		}
		time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
	}
	return resp, body, err
}

func NewRequest(opts *RequestOptions) (*http.Request, error) {
	if opts.Method == "" {
		opts.Method = http.MethodPost
	}
	var bodyBuffer io.Reader // 使用 io.Reader 接口，这样可以直接传递 nil
	var err error

	if opts.Body != nil {
		switch opts.BodyType {
		case BodyTypeForm:
			formData := make(url.Values)
			bodyMap, ok := opts.Body.(map[string]string)
			if !ok {
				return nil, fmt.Errorf("body is not a map[string]string")
			}
			for key, value := range bodyMap {
				formData.Set(key, value)
			}
			bodyBuffer = strings.NewReader(formData.Encode()) // 直接使用 strings.NewReader
		case BodyTypeMultipartForm:
			var buf bytes.Buffer
			writer := multipart.NewWriter(&buf)
			for fieldName, fileContent := range opts.Files {
				var part io.Writer
				part, err = writer.CreateFormFile(fieldName, fieldName)
				if err != nil {
					return nil, fmt.Errorf("failed to create form file for field '%s': %w", fieldName, err)
				}
				if _, err = part.Write(fileContent); err != nil {
					return nil, fmt.Errorf("failed to write content for field '%s': %w", fieldName, err)
				}
			}
			for key, value := range opts.Body.(map[string]string) {
				if err = writer.WriteField(key, value); err != nil {
					return nil, fmt.Errorf("failed to write field '%s': %w", key, err)
				}
			}
			err = writer.Close()
			if err != nil {
				return nil, fmt.Errorf("failed to close multipart writer: %w", err)
			}
			bodyBuffer = &buf
		default: // Default to JSON
			var data []byte
			data, err = json.Marshal(opts.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal body to JSON: %w", err)
			}
			bodyBuffer = bytes.NewReader(data)
		}
	}

	req, err := http.NewRequest(opts.Method, opts.URL, bodyBuffer) // 可以直接传递 nil 或有效的 io.Reader
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for k, v := range opts.Headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

func Execute(req *http.Request, cli *http.Client) (*http.Response, error) {
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", DefaultUserAgent)
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", DefaultContentType)
	}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	if resp.StatusCode >= 400 {
		return resp, fmt.Errorf("received HTTP error: %s", resp.Status)
	}
	return resp, nil
}

func Decode(resp *http.Response) ([]byte, error) {
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()
	if strings.Contains(strings.ToLower(resp.Header.Get("Content-Encoding")), "gzip") {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer func() {
			if reader != nil {
				_ = reader.Close()
			}
		}()
		return ioutil.ReadAll(reader)
	}
	return ioutil.ReadAll(resp.Body)
}

func Send(opts *RequestOptions) (*http.Response, []byte, error) {
	req, err := NewRequest(opts)
	if err != nil {
		return nil, nil, err
	}
	cli := opts.HttpClient
	if cli == nil {
		cli = defaultClient.httpClient
	}
	resp, err := Execute(req, cli)
	if err != nil {
		return resp, nil, err
	}
	body, err := Decode(resp)
	if err != nil {
		return resp, nil, err
	}
	if opts.ResponseOut != nil {
		if err := json.Unmarshal(body, opts.ResponseOut); err != nil {
			return resp, body, fmt.Errorf("failed to unmarshal response body: %w", err)
		}
	}
	return resp, body, nil
}
