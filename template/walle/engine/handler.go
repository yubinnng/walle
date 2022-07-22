package engine

import (
	"encoding/json"
	"net/http"

	handler "github.com/openfaas/templates-sdk/go-http"
)

type ExecuteRequest struct {
	ID   string       `json:"id"`
	Spec WorkflowSpec `json:"spec"`
}

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {
	var requestBody ExecuteRequest
	err := json.Unmarshal(req.Body, &requestBody)
	if err != nil {
		return handler.Response{
			Body:       []byte(err.Error()),
			StatusCode: http.StatusBadRequest,
		}, err
	}
	exec := NewExecution(requestBody.ID, requestBody.Spec)
	exec.Start()

	return handler.Response{
		StatusCode: http.StatusOK,
	}, nil
}
