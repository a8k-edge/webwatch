package main

import (
	"webwatch/db"
	"webwatch/server"
	"webwatch/task"
)

func main() {
	db.Init()

	go task.TaskManagerEventLoop()

	server.Serve()
}
