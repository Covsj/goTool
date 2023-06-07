package myHttp

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"strings"
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

func DoHttpRequest(req *http.Request) (*http.Response, error) {
	client := GetClient()
	defer ReleaseClient(client)

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

func CallHttp(url, method, body string, header map[string][]string) ([]byte, error) {
	req, err := NewHttpRequest(url, method, body, header)
	if err != nil {
		return nil, err
	}
	response, err := DoHttpRequest(req)
	if response == nil || err != nil {
		return nil, err
	}
	res, err := ReadRespBody(response)
	if err != nil {
		return nil, err
	}
	return res, err
}
