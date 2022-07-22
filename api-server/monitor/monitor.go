package monitor

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
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

func handleTaskMsg(m *nats.Msg) {
	mu.Lock()
	defer mu.Unlock()
	log.Printf("Received task event: %s\n", string(m.Data))
	var e workflow.TaskEvent
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
	exec.UpdateTaskStatus(e)
	// store into database
	if err := storage.Client.Save(exec).Error; err != nil {
		log.Printf("failed to update execution, %s", err.Error())
		return
	}
}
