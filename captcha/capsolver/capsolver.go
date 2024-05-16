package capsolver

import (
	"errors"
	"time"
)

func (c *CapSolver) Solve(task map[string]interface{}) (*CapSolverResponse, error) {
	capRes, err := request(CreateTaskUri, &capSolverRequest{Task: &task, ClientKey: c.getApiKey()})
	if err != nil {
		return nil, err
	}
	if capRes.ErrorId == 1 {
		return nil, errors.New(capRes.ErrorDescription)
	}
	if capRes.Status == StatusReady {
		return capRes, nil
	}
	for i := 0; i < TaskTimeout; i++ {
		capRes, err = request(GetTaskUri, &capSolverRequest{ClientKey: c.getApiKey(), TaskId: capRes.TaskId})
		time.Sleep(time.Second * 1)
		if err != nil {
			return nil, err
		}
		if capRes.ErrorId == 1 {
			return nil, errors.New(capRes.ErrorDescription)
		}
		if capRes.Status == StatusReady {
			break
		}
	}
	return capRes, err
}

func (c *CapSolver) Balance() (*CapSolverResponse, error) {
	capRes, err := request(BalanceUri, &capSolverRequest{ClientKey: c.getApiKey()})
	if err != nil {
		return nil, err
	}
	return capRes, nil
}

func (c *CapSolver) getApiKey() string {
	if c.ApiKey != "" {
		return c.ApiKey
	}
	return ApiKey
}
