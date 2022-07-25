package workflow

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

var ENGINE_URL = os.Getenv("WALLE_ENGINE_URL")

const (
	StatusWaiting = "WAITING"
	StatusRunning = "RUNNING"
	StatusAbort   = "ABORT"
	StatusSuccess = "SUCCESS"
	StatusError   = "ERROR"
)

type ExecutionTask struct {
	Name        string    `json:"name"`
	ExecutionID string    `json:"executionId"`
	Status      string    `json:"status"`
	Log         string    `json:"log"`
	StartedAt   time.Time `json:"startedAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ExecutionTasks []ExecutionTask

func (tasks ExecutionTasks) Value() (driver.Value, error) {
	bytes, err := json.Marshal(tasks)
	return bytes, err
}

func (tasks *ExecutionTasks) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), tasks)
}

// func (tasks ExecutionTasks) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(tasks)
// }

// func (tasks *ExecutionTasks) UnmarshalJSON(b []byte) error {
// 	return json.Unmarshal(b, tasks)
// }

type Execution struct {
	ID           string         `json:"id"`
	WorkflowName string         `json:"workflowName"`
	Status       string         `json:"status"`
	Tasks        ExecutionTasks `json:"tasks"`
	Spec         string         `json:"-" gorm:"-"`
	StartAt      time.Time      `json:"startAt"`
	EndAt        time.Time      `json:"endAt"`
}

func (Execution) TableName() string {
	return "execution"
}

func NewExecution(yamlSpec string) (*Execution, error) {
	now := time.Now()
	// parse spec
	spec, err := ParseYamlSpec(yamlSpec)
	if err != nil {
		return nil, err
	}
	id := uuid.New().String()
	tasks := make([]ExecutionTask, len(spec.Tasks))
	for i, taskSpec := range spec.Tasks {
		tasks[i] = ExecutionTask{
			ExecutionID: id,
			Name:        taskSpec.Name,
			Status:      StatusWaiting,
			UpdatedAt:   now,
		}
	}
	return &Execution{
		ID:           id,
		WorkflowName: spec.Name,
		Status:       StatusRunning,
		Tasks:        tasks,
		Spec:         yamlSpec,
		StartAt:      now,
	}, nil
}

type ExecuteRequest struct {
	ID   string       `json:"id"`
	Spec WorkflowSpec `json:"spec"`
}

func (exec *Execution) Start() error {
	spec, err := ParseYamlSpec(exec.Spec)
	if err != nil {
		return err
	}
	request := ExecuteRequest{
		ID:   exec.ID,
		Spec: spec,
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		return err
	}
	// send execution request
	resp, err := http.Post(ENGINE_URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("error request, response status: %s", resp.Status)
	}
	return nil
}

func (exec *Execution) UpdateTaskStatus(e ExecutionTask) {
	if exec.Status == StatusRunning {
		// assume execution is success
		exec.Status = StatusSuccess
	}
	for i := 0; i < len(exec.Tasks); i++ {
		task := &exec.Tasks[i]
		// find the task
		if task.Name == e.Name {
			task.Status = e.Status
			task.Log = e.Log
			task.UpdatedAt = e.UpdatedAt
			if task.Status == StatusRunning {
				task.StartedAt = e.UpdatedAt
				exec.Status = StatusRunning
			}
			if task.Status == StatusError || task.Status == StatusAbort {
				exec.Status = StatusError
			}
			break
		}
	}
	if exec.Status != StatusRunning {
		exec.EndAt = e.UpdatedAt
	}
}
