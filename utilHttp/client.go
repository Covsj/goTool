package utilHttp

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/abursavich/nett"
)

var dialer *nett.Dialer

var Timeout = 5 * time.Second
var DefaultClient *http.Client

func init() {
	DefaultClient = http.DefaultClient
}

func NewHTTPClient(timeout time.Duration) *http.Client {
	client := &http.Client{
		Timeout: timeout,
	}
	return client
}

func DoHttpRequest(req *http.Request, client *http.Client) (*http.Response, error) {
	if client == nil {
		client = DefaultClient
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

func CallHttp(url, method, body string, header map[string][]string) ([]byte, error) {
	req, err := NewHttpRequest(url, method, body, header)
	if err != nil {
		return nil, err
	}
	response, err := DoHttpRequest(req, DefaultClient)
	if response == nil || err != nil {
		return nil, err
	}
	res, err := ReadRespBody(response)
	if err != nil {
		return nil, err
	}
	return res, err
}
