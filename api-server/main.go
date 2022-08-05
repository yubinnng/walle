package main

import (
	"walle/api-server/api"
	"walle/api-server/monitor"
	"walle/api-server/storage"
	"walle/api-server/workflow"
)

func main() {
	r := api.SetupRoutes()
	storage.ConnectDB()
	// create tables
	storage.Client.AutoMigrate(&workflow.Metadata{}, &workflow.Execution{})
	monitor.Start()
	r.Run()
}
