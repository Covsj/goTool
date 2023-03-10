package utilHttp

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func NewHttpRequest(url, method, body string, header map[string][]string) (*http.Request, error) {

	req, err := http.NewRequest(method, url, ioutil.NopCloser(bytes.NewBuffer([]byte(body))))
	req.ContentLength = int64(len(body))

	if header != nil {
		for k, vv := range header {
			for _, v := range vv {
				req.Header.Add(k, v)
			}
		}
	}

	return req, err
}
