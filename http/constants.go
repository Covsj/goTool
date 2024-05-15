package http

const (
	DefaultUserAgent   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"
	DefaultContentType = "application/json"
	DefaultRetries     = 3
)

const (
	BodyTypeJSON BodyType = iota
	BodyTypeForm
	BodyTypeMultipartForm
)
