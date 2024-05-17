package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func Serve() {
	router := mux.NewRouter()

	router.HandleFunc("/", handleHome).Methods(http.MethodGet)
	router.HandleFunc("/add", handleAddTargetForm).Methods(http.MethodGet)
	router.HandleFunc("/add", handleAddTarget).Methods(http.MethodPost)
	router.HandleFunc("/target/{id:[0-9]+}", handleTargetHistory).Methods(http.MethodGet)
	router.HandleFunc("/target/{id:[0-9]+}", handleTargetDelete).Methods(http.MethodDelete)
	router.HandleFunc("/target/{id:[0-9]+}/toggle", handleToggleActive).Methods(http.MethodPost)

	fmt.Println("Starting UI server on port 8080...")
	if err := http.ListenAndServe("localhost:8080", router); err != nil {
		fmt.Println("Error starting UI server:", err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	var targets []Target
	db.Find(&targets)

	renderTemplate(w, []string{"templates/home.html"}, struct{ Rows []Target }{targets})
}

func handleAddTargetForm(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, []string{"templates/add.html"}, nil)
}

func handleAddTarget(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println("Error parsing form:", err)
		return
	}
	url := r.Form.Get("url")
	if url == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	name := r.Form.Get("name")
	if url == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	db.Create(&Target{Name: name, URL: url})

	http.Redirect(w, r, "/", http.StatusFound)
}

func handleTargetHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println("Invalid task ID:", err)
		return
	}

	var target Target
	if err := db.Preload("History", func(db *gorm.DB) *gorm.DB {
		return db.Order("histories.created_at DESC")
	}).First(&target, id).Error; err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error querying target:", err)
		return
	}

	data := struct {
		Name    string
		URL     string
		History []History
	}{
		Name:    target.Name,
		URL:     target.URL,
		History: target.History,
	}

	renderTemplate(w, []string{"templates/history.html"}, data)
}

func handleTargetDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println("Invalid task ID:", err)
		return
	}

	if err := db.Delete(&Target{}, id).Error; err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error querying target:", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleToggleActive(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println("Invalid task ID:", err)
		return
	}

	var target Target
	if err := db.Preload("History").First(&target, id).Error; err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error querying target:", err)
		return
	}

	target.IsActive = !target.IsActive
	db.Save(&target)
}

func renderTemplate(w http.ResponseWriter, templateFiles []string, data interface{}) {
	templateFiles = append([]string{"templates/__base.html"}, templateFiles...)

	tmpl, err := template.ParseFiles(templateFiles...)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error parsing template:", err)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing template:", err)
		return
	}
}
