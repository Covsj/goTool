package http

import "net/http"

type BodyType int

type Middleware func(*http.Request) (*http.Request, error)

type Client struct {
	httpClient *http.Client
}

var defaultClient *Client

type RequestOptions struct {
	URL         string
	Method      string
	Body        interface{}
	Headers     map[string]string
	Middlewares []Middleware
	Retries     int
	BodyType    BodyType
	Files       map[string][]byte // for multipart/form-data
	HttpClient  *http.Client
	ResponseOut interface{}
	CheckStatus bool
}
