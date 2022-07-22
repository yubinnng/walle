package workflow

import (
	"time"

	"gopkg.in/yaml.v3"
)

type WorkflowSpec struct {
	Name     string
	Desc     string
	Tasks    []TaskSpec
	Triggers []TriggerSpec
}

type TaskSpec struct {
	Name        string
	Type        string
	Url         string
	Method      string
	RequestBody map[string]string
	Retry       int
	Timeout     time.Duration
	Depends     []string
}

type TriggerSpec struct {
	Name string
	Type string
}

func ParseYamlSpec(yamlSpec string) (WorkflowSpec, error) {
	var spec WorkflowSpec
	err := yaml.Unmarshal([]byte(yamlSpec), &spec)
	return spec, err
}
