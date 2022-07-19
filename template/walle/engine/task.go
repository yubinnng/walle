package engine

import (
	"encoding/json"
	"fmt"
	"handler/event"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type TaskSpec struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Url         string            `json:"url"`
	Method      string            `json:"method"`
	RequestBody map[string]string `json:"requestBody"`
	Retry       int               `json:"retry"`
	Timeout     time.Duration     `json:"timeout"`
	Depends     []string          `json:"depends"`
}

const (
	StatusWaiting = "WAITING"
	StatusRunning = "RUNNING"
	StatusSuccess = "SUCCESS"
	StatusError   = "ERROR"
)

type TaskEvent struct {
	ExecutionID string
	TaskName    string
	TaskStatus  string
	CreatedAt   time.Time
}

type Task interface {
	GetName() string
	GetType() string
	Run() error
	GetStatus() string
	Ready() bool
	AddSubscriber(notifyChan chan TaskEvent)
	GetNotifyChan() chan TaskEvent
	String() string
}

type HttpTask struct {
	ExecID      string
	Name        string
	Url         string
	Method      string
	RequestBody map[string]string
	Retry       int
	Timeout     time.Duration
	Dependency  map[string]bool
	Subscribers []chan TaskEvent
	NotifyChan  chan TaskEvent
	Status      string
}

func NewHttpTask(taskSpec TaskSpec, execID string) Task {
	// set default value
	if taskSpec.Timeout <= 0 {
		taskSpec.Timeout = 10 * time.Second
	}
	if taskSpec.Retry < 0 {
		taskSpec.Retry = 3
	}
	task := &HttpTask{
		ExecID:      execID,
		Name:        taskSpec.Name,
		Url:         taskSpec.Url,
		Method:      taskSpec.Method,
		RequestBody: taskSpec.RequestBody,
		Retry:       taskSpec.Retry,
		Timeout:     taskSpec.Timeout,
		Status:      StatusWaiting,
		NotifyChan:  make(chan TaskEvent),
	}
	// build dependency set
	task.Dependency = make(map[string]bool)
	for _, dep := range taskSpec.Depends {
		task.Dependency[dep] = false
	}
	return task
}

func (task *HttpTask) GetName() string {
	return task.Name
}

func (task *HttpTask) GetType() string {
	return "http"
}

func (task *HttpTask) Run() error {
	task.SetStatus(StatusWaiting)
	for len(task.Dependency) > 0 {
		// fmt.Printf("%q waiting for %v tasks\n", task.Name, task.Dependency)
		event := <-task.NotifyChan
		if event.TaskStatus == StatusSuccess {
			delete(task.Dependency, event.TaskName)
		}
		if event.TaskStatus == StatusError {
			fmt.Printf("task %q interupted\n", task.Name)
			task.SetStatus(StatusError)
			return fmt.Errorf("task %q interupted", task.Name)
		}
	}
	task.SetStatus(StatusRunning)
	client := &http.Client{
		Timeout: task.Timeout,
	}
	req, err := http.NewRequest(task.Method, task.Url, nil)
	if err != nil {
		task.SetStatus(StatusError)
		fmt.Printf("Task %q occurred error: %q \n", task.Name, err.Error())
		return fmt.Errorf("error occurred in task %q, error: %q", task.GetName(), err.Error())
	}
	resp, err := client.Do(req)
	if err != nil {
		task.SetStatus(StatusError)
		fmt.Printf("Task %q occurred error: %q \n", task.Name, err.Error())
		return fmt.Errorf("error occurred in task %q, error: %q", task.GetName(), err.Error())
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		task.SetStatus(StatusError)
		fmt.Printf("task %q occurred error: %q \n", task.Name, err.Error())
		return fmt.Errorf("error occurred in task %q, error: %q", task.GetName(), err.Error())
	}
	if resp.StatusCode != 200 {
		task.SetStatus(StatusError)
		fmt.Printf("Task %q occurred error: %q \n", task.Name, responseBody)
		return fmt.Errorf("error occurred in task %q, error: %q", task.GetName(), responseBody)
	}
	task.SetStatus(StatusSuccess)
	fmt.Printf("Task %q successfully done, resp: %q \n", task.Name, responseBody)
	return nil
}

func (task *HttpTask) GetStatus() string {
	return task.Status
}

func (task *HttpTask) SetStatus(status string) {
	task.Status = status
	e := TaskEvent{
		ExecutionID: task.ExecID,
		TaskName:    task.Name,
		TaskStatus:  status,
		CreatedAt:   time.Now(),
	}
	// notify local subscriber task
	task.notifyAll(e)
	// publish to messaging queue
	data, _ := json.Marshal(e)
	log.Printf("publish task event %v\n", e)
	event.Publish("task-status", data)
}

// notify all subscriber
func (task *HttpTask) notifyAll(event TaskEvent) {
	for _, subscriber := range task.Subscribers {
		// non-blocking send notification in case of the subscriber error
		select {
		case subscriber <- event:
			// fmt.Printf("%q send notification\n", task.Name)
		default:
		}
	}
}

func (task *HttpTask) Ready() bool {
	return task.Status == StatusWaiting && len(task.Dependency) == 0
}

func (task *HttpTask) TasksDone(tasks []Task) {
	for _, t := range tasks {
		if t.GetStatus() == StatusSuccess {
			delete(task.Dependency, t.GetName())
		}
	}
}

func (task *HttpTask) AddSubscriber(notifyChan chan TaskEvent) {
	task.Subscribers = append(task.Subscribers, notifyChan)
}

func (task *HttpTask) GetNotifyChan() chan TaskEvent {
	return task.NotifyChan
}

func (task *HttpTask) String() string {
	return task.Name
}
