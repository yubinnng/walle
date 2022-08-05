package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()
	// Cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	r.Use(cors.New(corsConfig))
	// Workflow API
	r.POST("/api/workflow", CreateWorkflow)
	r.GET("/api/workflow/:name", GetWorkflow)
	r.GET("/api/workflow/list", ListWorkflows)
	r.DELETE("/api/workflow/:name", RemoveWorkflow)
	// Workflow execution API
	r.POST("/api/workflow/:name/exec", ExecuteWorkflow)
	r.GET("/api/execution/:id", GetExecution)
	r.GET("/api/execution/list", ListExecution)
	return r
}
