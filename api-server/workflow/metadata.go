package workflow

import (
	"errors"
	"time"
)

type Metadata struct {
	Name      string    `json:"name" gorm:"primaryKey"`
	Desc      string    `json:"desc"`
	Spec      string    `json:"spec"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Metadata) TableName() string {
	return "metadata"
}

func New(yamlSpec string) *Metadata {
	// parse workflow specification
	spec := ParseYamlSpec(yamlSpec)
	return &Metadata{
		Name:      spec.Name,
		Desc:      spec.Desc,
		Spec:      yamlSpec,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (md *Metadata) Update(yamlSpec string) error {
	// parse workflow specification
	spec := ParseYamlSpec(yamlSpec)
	if md.Name != spec.Name {
		return errors.New("cannot modify workflow name")
	}
	md.Desc = spec.Desc
	md.Spec = yamlSpec
	md.UpdatedAt = time.Now()
	return nil
}
