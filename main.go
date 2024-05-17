package main

func main() {
	Start()
}

func Start() {
	initializeDB("db.db")

	go TaskManagerEventLoop()

	Serve()
}
