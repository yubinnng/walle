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
	r.POST("/workflow", CreateWorkflow)
	r.GET("/workflow/:name", GetWorkflow)
	r.GET("/workflow/list", ListWorkflows)
	r.DELETE("/workflow/:name", RemoveWorkflow)
	// Workflow execution API
	r.POST("/workflow/:name/exec", ExecuteWorkflow)
	r.GET("/execution/:id", GetExecution)
	r.GET("/execution/list", ListExecution)
	return r
}
