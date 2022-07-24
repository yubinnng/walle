package api

import (
	"errors"
	"fmt"
	"net/http"
	"walle/api-server/storage"
	"walle/api-server/workflow"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateWorkflowRequest struct {
	Spec string `binding:"required"`
}

func CreateWorkflow(c *gin.Context) {
	// read spec
	bytes, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	yamlSpec := string(bytes)
	// parse spec
	newWf, err := workflow.New(yamlSpec)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	// retrieve workflow from storage
	wf := &workflow.Metadata{}
	err = storage.Client.First(&wf, "name = ?", wf.Name).Error
	// unknown error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	if err == nil {
		// workflow exists, update metadata
		wf.Update(yamlSpec)
	} else {
		// workflow not exists, create new workflow
		wf = newWf
	}
	// save
	if err := storage.Client.Save(wf).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
}

type GetWorkflowResponse struct {
	workflow.Metadata
	ExecutionCount int
}

func GetWorkflow(c *gin.Context) {
	name := c.Param("name")
	var md workflow.Metadata
	err := storage.Client.First(&md, "name = ?", name).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("workflow %q is not exist", name)})
		return
	}
	// TODO count executions
	c.JSON(http.StatusOK, GetWorkflowResponse{
		Metadata:       md,
		ExecutionCount: 0,
	})
}

func ListWorkflows(c *gin.Context) {
	var list []workflow.Metadata
	storage.Client.Select("name").Order("created_at asc").Find(&list)
	names := make([]string, len(list))
	for i, wf := range list {
		names[i] = wf.Name
	}
	c.JSON(http.StatusOK, names)
}

func RemoveWorkflow(c *gin.Context) {
	name := c.Param("name")
	// remove from db
	storage.Client.Where("name = ?", name).Delete(&workflow.Metadata{})
}

func ExecuteWorkflow(c *gin.Context) {
	name := c.Param("name")
	// retrieve spec
	var md workflow.Metadata
	err := storage.Client.First(&md, "name = ?", name).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("workflow %q is not exist", name)})
		return
	}
	// new execution
	exec, err := workflow.NewExecution(md.Spec)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	// store execution into database
	if err := storage.Client.Create(&exec).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	if err := exec.Start(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
}

func GetExecution(c *gin.Context) {
	id := c.Param("id")
	var exec workflow.Execution
	err := storage.Client.First(&exec, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("execution %q is not exist", id)})
		return
	}
	c.JSON(http.StatusOK, exec)
}

func ListExecution(c *gin.Context) {
	workflowName := c.Query("workflow_name")
	var list []workflow.Execution
	storage.Client.Select("id", "status", "start_at", "end_at").Where("workflow_name = ?", workflowName).Order("start_at asc").Find(&list)
	c.JSON(http.StatusOK, list)
}
