package main

import (
	"HUObjStorageAPI/api"
	"HUObjStorageAPI/config"
	"HUObjStorageAPI/heartbeat"
)

func main() {
	config.Init()

	go heartbeat.ListenHeartbeat()

	api.InitRouter()
}
