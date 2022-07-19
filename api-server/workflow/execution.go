package workflow

import "time"

const (
	StatusWaiting = "WAITING"
	StatusRunning = "RUNNING"
	StatusError   = "ERROR"
	StatusSuccess = "SUCCESS"
)

type WorkflowExec struct {
	ID           string     `json:"id"`
	WorkflowName string     `json:"workflow_name"`
	Tasks        []TaskExec `json:"tasks"`
	StartAt      time.Time  `json:"start_at"`
	EndAt        time.Time  `json:"end_at"`
}

func NewExecution(id string, workflowName string, startAt time.Time) *WorkflowExec {
	return &WorkflowExec{
		ID:           id,
		WorkflowName: workflowName,
		Tasks:        []TaskExec{},
		StartAt:      startAt,
	}
}

func (WorkflowExec) TableName() string {
	return "execution"
}

type TaskExec struct {
	Status   string    `json:"status"`
	Response string    `json:"response"`
	Error    string    `json:"error"`
	StartAt  time.Time `json:"start_at"`
	EndAt    time.Time `json:"end_at"`
}
