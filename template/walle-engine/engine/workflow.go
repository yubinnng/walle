package engine

import (
	"log"

	"golang.org/x/sync/errgroup"
)

type WorkflowSpec struct {
	Version float32    `json:"version"`
	Name    string     `json:"name"`
	Desc    string     `json:"desc"`
	Tasks   []TaskSpec `json:"tasks"`
}

type Execution struct {
	ID           string
	WorkflowName string
	g            errgroup.Group
	taskMap      map[string]Task
}

const (
	WorkflowStatusRunning = "RUNNING"
	WorkflowStatusSuccess = "SUCCESS"
	WorkflowStatusError   = "ERROR"
)

func NewExecution(id string, wfSpec WorkflowSpec) *Execution {
	exec := &Execution{
		ID:           id,
		WorkflowName: wfSpec.Name,
		taskMap:      make(map[string]Task),
	}
	totalTasks := len(wfSpec.Tasks)
	for _, taskSpec := range wfSpec.Tasks {
		// deafult task type is HTTP task
		exec.taskMap[taskSpec.Name] = NewHttpTask(taskSpec, id, totalTasks)
	}
	for _, taskSpec := range wfSpec.Tasks {
		for _, dep := range taskSpec.Depends {
			task := exec.taskMap[taskSpec.Name]
			exec.taskMap[dep].AddSubscriber(task)
		}
	}
	return exec
}

func (exec *Execution) Start() error {
	log.Printf("Execution %q: workflow %q start\n", exec.ID, exec.WorkflowName)
	for _, task := range exec.taskMap {
		exec.g.Go(task.Run)
	}
	if err := exec.g.Wait(); err != nil {
		return err
	}
	log.Printf("Execution %q: workflow %q end\n\n", exec.ID, exec.WorkflowName)
	return nil
}

func (exec *Execution) Interupt() {

}
