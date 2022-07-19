package engine

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	handler "github.com/openfaas/templates-sdk/go-http"
)

type ExecuteRequest struct {
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
	log.Println(requestBody)

	wf := NewWorkflow(requestBody.Spec, uuid.New().String())
	// spec := WorkflowSpec{}
	// file, err := ioutil.ReadFile("./workflow.yaml")
	// yaml.Unmarshal(file, &spec)
	// wf := NewWorkflow(spec, "1")
	wf.Start()

	return handler.Response{
		StatusCode: http.StatusOK,
	}, nil
}
