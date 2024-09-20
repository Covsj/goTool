package http

import (
	"github.com/imroc/req/v3"
)

// https://req.cool/zh/docs/prologue/introduction/

var defaultClient = req.C().
	ImpersonateChrome().
	DisableDumpAll()

type ReqOpt struct {
	Method               string
	Url                  string
	Headers              map[string]string
	Result               interface{}
	RetryCount           int
	Body                 interface{}
	FormData             map[string]string
	EnableForceMultipart bool
}

func DoRequest(opt ReqOpt) (*req.Response, error) {
	request := defaultClient.R()
	if opt.RetryCount != 0 {
		request = request.SetRetryCount(opt.RetryCount)
	}
	if opt.Headers != nil {
		request = request.SetHeaders(opt.Headers)
	}
	if opt.Result != nil {
		request = request.SetSuccessResult(opt.Result)
	}
	if opt.Body != nil {
		request = request.SetBody(opt.Body)
	}
	if opt.EnableForceMultipart {
		request = request.EnableForceMultipart()
	}
	if opt.FormData != nil {
		request = request.SetFormData(opt.FormData)
	}
	if opt.Method == "GET" {
		return request.Get(opt.Url)
	} else if opt.Method == "POST" {
		return request.Post(opt.Url)
	} else if opt.Method == "PUT" {
		return request.Put(opt.Url)
	} else if opt.Method == "DELETE" {
		return request.Delete(opt.Url)
	} else if opt.Method == "OPTIONS" {
		return request.Options(opt.Url)
	}

	return request.Head(opt.Url)
}
