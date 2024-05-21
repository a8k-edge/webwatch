package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Serve() {
	router := mux.NewRouter()

	router.HandleFunc("/", handleHome).Methods(http.MethodGet)
	router.HandleFunc("/add", handleAddTargetForm).Methods(http.MethodGet)
	router.HandleFunc("/add", handleAddTarget).Methods(http.MethodPost)
	router.HandleFunc("/target/{id:[0-9]+}", handleTargetHistory).Methods(http.MethodGet)
	router.HandleFunc("/target/{id:[0-9]+}", handleTargetDelete).Methods(http.MethodDelete)
	router.HandleFunc("/target/{id:[0-9]+}/toggle", handleToggleActive).Methods(http.MethodPost)
	router.HandleFunc("/target/{tid:[0-9]+}/history/{hid:[0-9]+}", handleHistoryView).Methods(http.MethodGet)

	// not found
	http.Handle("/", router)

	fmt.Println("Starting UI server on port 8080...")
	if err := http.ListenAndServe("localhost:8080", router); err != nil {
		fmt.Println("Error starting UI server:", err)
	}
}
