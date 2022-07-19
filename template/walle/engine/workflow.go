package engine

import (
	"encoding/json"
	"fmt"
	"handler/event"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

type WorkflowSpec struct {
	Version float32    `json:"version"`
	Name    string     `json:"name"`
	Desc    string     `json:"desc"`
	Tasks   []TaskSpec `json:"tasks"`
}

type Workflow struct {
	ExecutionID string
	Name        string
	Description string
	g           errgroup.Group
	taskMap     map[string]Task
}

type WorkflowEvent struct {
	ExecutionID    string
	WorkflowName   string
	WorkflowStatus string
	CreatedAt      time.Time
}

const (
	WorkflowStatusRunning = "RUNNING"
	WorkflowStatusSuccess = "SUCCESS"
	WorkflowStatusError   = "ERROR"
)

func NewWorkflow(wfSpec WorkflowSpec, executionID string) *Workflow {
	wf := &Workflow{
		ExecutionID: executionID,
		Name:        wfSpec.Name,
		Description: wfSpec.Desc,
		taskMap:     make(map[string]Task),
	}
	for _, taskSpec := range wfSpec.Tasks {
		if taskSpec.Type == "http" {
			wf.taskMap[taskSpec.Name] = NewHttpTask(taskSpec, executionID)
		}
	}
	for _, taskSpec := range wfSpec.Tasks {
		for _, dep := range taskSpec.Depends {
			task := wf.taskMap[taskSpec.Name]
			wf.taskMap[dep].AddSubscriber(task.GetNotifyChan())
		}
	}
	return wf
}

func (wf *Workflow) Start() error {
	fmt.Printf("Workflow %q start \n", wf.Name)
	wf.PublishStatus(WorkflowStatusRunning)
	for _, task := range wf.taskMap {
		wf.g.Go(task.Run)
	}
	if err := wf.g.Wait(); err != nil {
		fmt.Printf("workflow %q occurred error\n\n", wf.Name)
		wf.PublishStatus(WorkflowStatusError)
		return err
	}
	fmt.Printf("Workflow %q end \n\n", wf.Name)
	wf.PublishStatus(WorkflowStatusSuccess)
	return nil
}

func (wf *Workflow) Interupt() {

}

func (wf *Workflow) PublishStatus(status string) {
	e := WorkflowEvent{
		ExecutionID:    wf.ExecutionID,
		WorkflowName:   wf.Name,
		WorkflowStatus: status,
		CreatedAt:      time.Now(),
	}
	// publish to messaging queue
	data, _ := json.Marshal(e)
	log.Printf("publish workflow event %v\n", e)
	event.Publish("workflow-status", data)
}
