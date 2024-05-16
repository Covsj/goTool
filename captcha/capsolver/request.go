package capsolver

import (
	"fmt"

	"github.com/Covsj/goTool/http"
)

func request(uri string, solverRequest *Request) (*Response, error) {
	capResponse := &Response{}
	_, _, err := http.SendWithRetries(&http.RequestOptions{
		Retries:     2,
		URL:         fmt.Sprintf("%s%s", ApiHost, uri),
		Body:        solverRequest,
		Headers:     map[string]string{"Content-Type": "application/json"},
		ResponseOut: capResponse,
	})

	if err != nil {
		return nil, err
	}
	return capResponse, nil
}
