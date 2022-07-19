package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()
	// Workflow API
	r.POST("/workflow", CreateWorkflow)
	r.PUT("/workflow/:name", UpdateWorkflow)
	r.GET("/workflow/:name", GetWorkflow)
	r.GET("/workflow/list", ListWorkflows)
	r.DELETE("/workflow/:name", RemoveWorkflow)
	// Workflow execution API
	r.POST("/workflow/:name/exec", ExecuteWorkflow)
	r.GET("/workflow/:name/exec", GetExecution)
	r.GET("/workflow/:name/exec/list", ListExecution)
	return r
}
