package utilHttp

import (
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"strings"
)

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
