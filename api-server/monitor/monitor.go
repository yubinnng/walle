package monitor

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
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
	client.Subscribe("workflow-status", handleWorkflowMsg)
	client.Subscribe("task-status", handleTaskMsg)
}

type WorkflowEvent struct {
	ExecutionID    string
	WorkflowName   string
	WorkflowStatus string
	CreatedAt      time.Time
}

const (
	StatusWaiting = "WAITING"
	StatusRunning = "RUNNING"
	StatusError   = "ERROR"
	StatusSuccess = "SUCCESS"
)

func handleWorkflowMsg(m *nats.Msg) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Printf("Received workflow event: %s\n", string(m.Data))
	var e WorkflowEvent
	if err := json.Unmarshal(m.Data, &e); err != nil {
		log.Println("failed to parse event")
	}
	// if e.WorkflowStatus == workflow.StatusRunning {

	// }
}

type TaskEvent struct {
	ExecutionID string
	TaskName    string
	TaskStatus  string
	CreatedAt   time.Time
}

func handleTaskMsg(m *nats.Msg) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Printf("Received task event: %s\n", string(m.Data))
}
