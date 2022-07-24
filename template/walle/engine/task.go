package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"handler/event"
	"log"
	"net/http"
	"time"
)

type TaskSpec struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Url         string            `json:"url"`
	Header      map[string]string `json:"header"`
	Method      string            `json:"method"`
	RequestBody map[string]string `json:"requestBody"`
	Retry       int               `json:"retry"`
	Timeout     time.Duration     `json:"timeout"`
	Depends     []string          `json:"depends"`
}

const (
	StatusWaiting = "WAITING"
	StatusRunning = "RUNNING"
	StatusAbort   = "ABORT"
	StatusSuccess = "SUCCESS"
	StatusError   = "ERROR"
)

type TaskEvent struct {
	ExecutionID string    `json:"executionId"`
	TaskName    string    `json:"taskName"`
	TaskStatus  string    `json:"taskStatus"`
	TaskLog     string    `json:"taskLog"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Task interface {
	Name() string
	Type() string
	Run() error
	EventChan() chan TaskEvent
	AddSubscriber(subscriber Task)
}

type HttpTask struct {
	executionID string
	name        string
	url         string
	method      string
	header      map[string]string
	requestBody map[string]string
	retry       int
	timeout     time.Duration
	status      string
	dependency  map[string]bool
	subscribers []chan TaskEvent
	eventChan   chan TaskEvent
}

func NewHttpTask(taskSpec TaskSpec, executionID string, totalTask int) Task {
	// set default value
	if taskSpec.Timeout <= 0 {
		taskSpec.Timeout = 10 * time.Second
	}
	if taskSpec.Retry < 0 {
		taskSpec.Retry = 3
	}
	task := &HttpTask{
		executionID: executionID,
		name:        taskSpec.Name,
		url:         taskSpec.Url,
		method:      taskSpec.Method,
		header:      taskSpec.Header,
		requestBody: taskSpec.RequestBody,
		retry:       taskSpec.Retry,
		timeout:     taskSpec.Timeout,
		status:      StatusWaiting,
		eventChan:   make(chan TaskEvent, totalTask),
	}
	// build dependency set
	task.dependency = make(map[string]bool)
	for _, dep := range taskSpec.Depends {
		task.dependency[dep] = false
	}
	return task
}

func (task *HttpTask) Name() string {
	return task.name
}

func (task *HttpTask) Type() string {
	return "http"
}

func (task *HttpTask) request() error {
	client := &http.Client{
		Timeout: task.timeout,
	}
	jsonBody, err := json.Marshal(task.requestBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(task.method, task.url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	for k, v := range task.header {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 == 2 {
		return nil
	}
	return fmt.Errorf("error request, response status: %s", resp.Status)
}

func (task *HttpTask) Run() error {
	for len(task.dependency) > 0 {
		event := <-task.eventChan
		if event.TaskStatus == StatusSuccess {
			delete(task.dependency, event.TaskName)
		}
		if event.TaskStatus == StatusError {
			log.Printf("Execution %q: task %q abort\n", task.executionID, task.name)
			task.status = StatusAbort
			task.PublishEvent(TaskEvent{
				ExecutionID: task.executionID,
				TaskName:    task.name,
				TaskStatus:  task.status,
				UpdatedAt:   time.Now(),
			})
			return fmt.Errorf("Execution %q: task %q abort", task.executionID, task.name)
		}
	}
	// all dependencies satisfied, close event channel
	close(task.eventChan)
	task.status = StatusRunning
	task.PublishEvent(TaskEvent{
		ExecutionID: task.executionID,
		TaskName:    task.name,
		TaskStatus:  task.status,
		UpdatedAt:   time.Now(),
	})
	log.Printf("Execution %q: task %q start\n", task.executionID, task.name)

	if err := task.request(); err != nil {
		task.status = StatusError
		task.PublishEvent(TaskEvent{
			ExecutionID: task.executionID,
			TaskName:    task.name,
			TaskStatus:  task.status,
			TaskLog:     err.Error(),
			UpdatedAt:   time.Now(),
		})
		log.Printf("Execution %q: http task %q request error\n", task.executionID, task.name)
		return fmt.Errorf("Execution %q: http task %q request error", task.executionID, task.name)
	}
	task.status = StatusSuccess
	task.PublishEvent(TaskEvent{
		ExecutionID: task.executionID,
		TaskName:    task.name,
		TaskStatus:  task.status,
		UpdatedAt:   time.Now(),
	})
	log.Printf("Execution %q: task %q done successfully\n", task.executionID, task.name)
	return nil
}

func (task *HttpTask) PublishEvent(e TaskEvent) {
	// notify local subscriber task
	task.notifyAll(e)
	// publish to messaging queue
	data, _ := json.Marshal(e)
	event.Publish("task-status", data)
}

// notify all subscriber
func (task *HttpTask) notifyAll(event TaskEvent) {
	for _, subscriber := range task.subscribers {
		subscriber <- event
	}
}

func (task *HttpTask) AddSubscriber(subscriber Task) {
	task.subscribers = append(task.subscribers, subscriber.EventChan())
}

func (task *HttpTask) EventChan() chan TaskEvent {
	return task.eventChan
}
