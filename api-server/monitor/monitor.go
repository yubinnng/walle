package monitor

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"
	"walle/api-server/storage"
	"walle/api-server/workflow"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

var client *nats.Conn
var mu sync.Mutex

func Start() {
	// Connect to a server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("cannot connect to NATS")
	}
	client = nc
	log.Println("connected to NATS")
	// subscribe workflow status topic
	// client.Subscribe("workflow-status", handleWorkflowMsg)
	client.Subscribe("task-status", handleTaskMsg)
}

const (
	StatusWaiting = "WAITING"
	StatusRunning = "RUNNING"
	StatusError   = "ERROR"
	StatusSuccess = "SUCCESS"
)

type TaskEvent struct {
	ExecutionID string    `json:"executionId"`
	TaskName    string    `json:"taskName"`
	TaskStatus  string    `json:"taskStatus"`
	TaskLog     string    `json:"taskLog"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func handleTaskMsg(m *nats.Msg) {
	mu.Lock()
	defer mu.Unlock()
	log.Printf("Received task event: %s\n", string(m.Data))
	var e TaskEvent
	if err := json.Unmarshal(m.Data, &e); err != nil {
		log.Println("failed to parse event")
	}
	// retrieve execution
	exec := &workflow.Execution{}
	err := storage.Client.First(exec, "id = ?", e.ExecutionID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("execution %q is not exist\n", e.ExecutionID)
		return
	}
	// update task status
	exec.UpdateTaskStatus(workflow.ExecutionTask{
		Name: e.TaskName,
		ExecutionID: e.ExecutionID,
		Status: e.TaskStatus,
		Log: e.TaskLog,
		UpdatedAt: e.UpdatedAt,
	})
	// store into database
	if err := storage.Client.Save(exec).Error; err != nil {
		log.Printf("failed to update execution, %s", err.Error())
		return
	}
}
