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

type ExecutionTasks []TaskEvent

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
	WorkflowName string         `json:"workflow_name"`
	Tasks        ExecutionTasks `json:"tasks"`
	Spec         string         `json:"-" gorm:"-"`
	StartAt      time.Time      `json:"start_at"`
	EndAt        time.Time      `json:"end_at"`
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
	tasks := make([]TaskEvent, len(spec.Tasks))
	for i, taskSpec := range spec.Tasks {
		tasks[i] = TaskEvent{
			ExecutionID: id,
			TaskName:    taskSpec.Name,
			TaskStatus:  StatusWaiting,
			UpdatedAt:   now,
		}
	}
	return &Execution{
		ID:           id,
		WorkflowName: spec.Name,
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

func (exec *Execution) UpdateTaskStatus(e TaskEvent) {
	allDone := true
	for i := 0; i < len(exec.Tasks); i++ {
		task := &exec.Tasks[i]
		// update tasks
		if task.TaskName == e.TaskName {
			task.TaskStatus = e.TaskStatus
			task.TaskLog = e.TaskLog
			task.UpdatedAt = e.UpdatedAt
		}
		if task.TaskStatus != StatusSuccess {
			allDone = false
		}
	}
	if allDone {
		exec.EndAt = e.UpdatedAt
	}
}

type TaskEvent struct {
	ExecutionID string    `json:"execution_id"`
	TaskName    string    `json:"task_name"`
	TaskStatus  string    `json:"task_status"`
	TaskLog     string    `json:"task_log"`
	UpdatedAt   time.Time `json:"updated_at"`
}
