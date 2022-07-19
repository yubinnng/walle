package main

import (
	"walle/api-server/api"
	"walle/api-server/monitor"
	"walle/api-server/storage"
)

func main() {
	r := api.SetupRoutes()
	storage.ConnectDB()
	monitor.Start()
	r.Run()
}
