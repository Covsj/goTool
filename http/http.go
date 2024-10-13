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
	RespOut              interface{}
	RetryCount           int
	Body                 interface{}
	FormData             map[string]interface{}
	EnableForceMultipart bool
	NeedSkipTls          bool
}

func DoRequest(opt *ReqOpt) (*req.Response, error) {
	cli := defaultClient
	if opt.NeedSkipTls {
		cli.TLSClientConfig.InsecureSkipVerify = true
	}
	request := cli.R()
	if opt.RetryCount != 0 {
		request = request.SetRetryCount(opt.RetryCount)
	}
	if opt.Headers != nil {
		request = request.SetHeaders(opt.Headers)
	}
	if opt.RespOut != nil {
		request = request.SetSuccessResult(opt.RespOut)
	}
	if opt.Body != nil {
		request = request.SetBody(opt.Body)
	}
	if opt.EnableForceMultipart {
		request = request.EnableForceMultipart()
	}
	if opt.FormData != nil {
		request = request.SetFormDataAnyType(opt.FormData)
	}

	switch opt.Method {
	case "POST":
		return request.Post(opt.Url)
	case "PUT":
		return request.Put(opt.Url)
	case "DELETE":
		return request.Delete(opt.Url)
	case "OPTIONS":
		return request.Options(opt.Url)
	case "HEAD":
		return request.Head(opt.Url)
	default:
		return request.Get(opt.Url)
	}
}
