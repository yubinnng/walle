package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"walle/api-server/storage"
	"walle/api-server/workflow"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var ENGINE_URL = os.Getenv("WALLE_ENGINE_URL")

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
	// fmt.Println(yamlSpec)
	wf := workflow.New(yamlSpec)
	// fmt.Println(wf)
	if err := storage.Client.Create(&wf).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
}

func UpdateWorkflow(c *gin.Context) {
	name := c.Param("name")
	// retrieve workflow from storage
	var md workflow.Metadata
	err := storage.Client.First(&md, "name = ?", name).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("workflow %q is not exist", name)})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	// read spec
	bytes, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	yamlSpec := string(bytes)
	// update
	err = md.Update(yamlSpec)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	// save
	if err := storage.Client.Save(md).Error; err != nil {
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
	var results []workflow.Metadata
	storage.Client.Order("created_at asc").Find(&results)
	c.JSON(http.StatusOK, results)
}

func RemoveWorkflow(c *gin.Context) {
	name := c.Param("name")
	// remove from db
	storage.Client.Where("name = ?", name).Delete(&workflow.Metadata{})
}

type ExecuteRequest struct {
	Spec workflow.WorkflowSpec `json:"spec"`
}

func ExecuteWorkflow(c *gin.Context) {
	name := c.Param("name")
	var md workflow.Metadata
	err := storage.Client.First(&md, "name = ?", name).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("workflow %q is not exist", name)})
		return
	}
	// parse spec
	request := ExecuteRequest{
		Spec: workflow.ParseYamlSpec(md.Spec),
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	// send execution request
	resp, err := http.Post(ENGINE_URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)
		resp := string(respBody)
		c.JSON(http.StatusInternalServerError, gin.H{"error": resp})
		c.Error(errors.New(resp))
		return
	}
}

func ListExecution(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func GetExecution(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
