package server

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"webwatch/db"

	"github.com/aryann/difflib"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	var targets []db.Target
	db.GetDB().Find(&targets)

	renderTemplate(w, []string{"templates/home.html"}, struct{ Rows []db.Target }{targets})
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

	db.GetDB().Create(&db.Target{Name: name, URL: url})

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

	var target db.Target
	if err := db.GetDB().Preload("History", func(db *gorm.DB) *gorm.DB {
		return db.Order("histories.created_at DESC")
	}).First(&target, id).Error; err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error querying target:", err)
		return
	}

	data := struct {
		Target  db.Target
		History []db.History
	}{
		Target:  target,
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

	if err := db.GetDB().Delete(&db.Target{}, id).Error; err != nil {
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

	var target db.Target
	if err := db.GetDB().Preload("History").First(&target, id).Error; err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error querying target:", err)
		return
	}

	target.IsActive = !target.IsActive
	db.GetDB().Save(&target)
}

func handleHistoryView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["hid"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println("Invalid task ID:", err)
		return
	}

	var history db.History
	if err := db.GetDB().First(&history, id).Error; err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error querying target:", err)
		return
	}

	diff := make([]difflib.DiffRecord, 0)
	if err := json.Unmarshal([]byte(history.Diff), &diff); err != nil {
		log.Println("Error unmarshal prev history diff:", err)
		return
	}
	data := struct {
		History db.History
		Diff    []difflib.DiffRecord
	}{
		History: history,
		Diff:    diff,
	}

	renderTemplate(w, []string{"templates/history_view.html"}, data)
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
