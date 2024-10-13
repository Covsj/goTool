package capsolver

import (
	"fmt"

	"github.com/Covsj/goTool/http"
)

func request(uri string, solverRequest *Request) (*Response, error) {
	capResponse := &Response{}
	_, err := http.DoRequest(&http.ReqOpt{
		Url:     fmt.Sprintf("%s%s", ApiHost, uri),
		Body:    solverRequest,
		Headers: map[string]string{"Content-Type": "application/json"},
		RespOut: capResponse,
	})

	if err != nil {
		return nil, err
	}
	return capResponse, nil
}
