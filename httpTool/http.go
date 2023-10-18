package httpTool

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

const (
	DefaultUserAgent   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"
	DefaultContentType = "application/json"
	DefaultRetries     = 3
)

type BodyType int

const (
	BodyTypeJSON BodyType = iota
	BodyTypeForm
	BodyTypeMultipartForm
)

type Middleware func(*http.Request) (*http.Request, error)

type Client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
	}
}

type RequestOptions struct {
	URL         string
	Method      string
	Body        interface{}
	Headers     map[string]string
	Middlewares []Middleware
	Retries     int
	BodyType    BodyType
	Files       map[string][]byte // for multipart/form-data
	ResponseOut interface{}
}

func (c *Client) Get(url string, out interface{}) (*http.Response, []byte, error) {
	opts := &RequestOptions{
		URL:         url,
		Method:      http.MethodGet,
		ResponseOut: out,
	}
	return c.Send(opts)
}

func (c *Client) Post(url string, body interface{}, out interface{}) (*http.Response, []byte, error) {
	opts := &RequestOptions{
		URL:         url,
		Method:      http.MethodPost,
		Body:        body,
		ResponseOut: out,
	}
	return c.Send(opts)
}

func (c *Client) SendWithRetries(opts *RequestOptions) (*http.Response, []byte, error) {
	if opts.Retries == 0 {
		opts.Retries = DefaultRetries
	}
	var resp *http.Response
	var body []byte
	var err error
	for i := 0; i < opts.Retries; i++ {
		resp, body, err = c.Send(opts)
		if err == nil {
			return resp, body, nil
		}
		time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
	}
	return resp, body, err
}
func (c *Client) NewRequest(opts *RequestOptions) (*http.Request, error) {
	if opts.Method == "" {
		opts.Method = http.MethodPost
	}
	var bodyBuffer *bytes.Buffer
	var err error

	switch opts.BodyType {
	case BodyTypeForm:
		formData := make(url.Values)
		for key, value := range opts.Body.(map[string]string) {
			formData.Set(key, value)
		}
		bodyBuffer = bytes.NewBufferString(formData.Encode())
	case BodyTypeMultipartForm:
		bodyBuffer = &bytes.Buffer{}
		writer := multipart.NewWriter(bodyBuffer)
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
	default: // Default to JSON
		var data []byte
		data, err = json.Marshal(opts.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body to JSON: %w", err)
		}
		bodyBuffer = bytes.NewBuffer(data)
	}

	req, err := http.NewRequest(opts.Method, opts.URL, bodyBuffer)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for k, v := range opts.Headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

func (c *Client) Execute(req *http.Request) (*http.Response, error) {
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", DefaultUserAgent)
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", DefaultContentType)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	if resp.StatusCode >= 400 {
		return resp, fmt.Errorf("received HTTP error: %s", resp.Status)
	}
	return resp, nil
}

func (c *Client) Decode(resp *http.Response) ([]byte, error) {
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

func (c *Client) Send(opts *RequestOptions) (*http.Response, []byte, error) {
	req, err := c.NewRequest(opts)
	if err != nil {
		return nil, nil, err
	}
	resp, err := c.Execute(req)
	if err != nil {
		return resp, nil, err
	}
	body, err := c.Decode(resp)
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
