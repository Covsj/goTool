package httpUtil

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	localhttp "github.com/Covsj/goTool/http"
	"github.com/Covsj/goTool/log"
)

var (
	httpClient = &http.Client{
		Timeout: 50 * time.Second,
	}
)

type Res struct {
	Header http.Header
	Body   []byte
	Code   int
}

type RequestParams struct {
	Header  map[string]string
	Body    []byte
	Timeout time.Duration
}

func Get(baseUrl string, param map[string]string, out interface{}) (err error) {
	urlPath := baseUrl
	if len(param) != 0 {
		params := url.Values{}
		for k, v := range param {
			params.Set(k, v)
		}
		httpUrl, err := url.Parse(baseUrl)
		if err != nil {
			return err
		}
		httpUrl.RawQuery = params.Encode()
		urlPath = httpUrl.String()
	}

	resp, err := localhttp.DoRequest(&localhttp.ReqOpt{Url: urlPath, RespOut: out})
	if err != nil {
		log.Error("请求失败", "url", urlPath, "error", err, "resp", resp)
		return err
	}
	return nil
}

func Post(url string, params RequestParams) ([]byte, error) {
	ctx := context.Background()
	if params.Timeout > 0 {
		ctx2, cancel := context.WithTimeout(ctx, params.Timeout)
		ctx = ctx2
		defer cancel()
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(params.Body))
	if err != nil {
		return nil, err
	}
	for k, v := range params.Header {
		request.Header.Set(k, v)
	}

	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, fmt.Errorf("post %v response code = %v", url, resp.Status)
	}
	res, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return nil, errRead
	}
	return res, nil
}

func Request(method, url string, header map[string]string, body []byte) (*Res, error) {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		request.Header.Set(k, v)
	}

	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	res, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return nil, errRead
	}
	return &Res{
		Header: resp.Header,
		Body:   res,
		Code:   resp.StatusCode,
	}, nil
}

func RequestSync(method, url string, header map[string]string, body []byte, httpWg *sync.WaitGroup, resp *Res, err error) {
	response, err := Request(method, url, header, body)
	if err == nil {
		*resp = *response
	}
	httpWg.Done()
}
