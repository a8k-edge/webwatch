package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"webwatch/db"

	"github.com/aryann/difflib"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func handleHome(c echo.Context) error {
	var targets []db.Target
	db.GetDB().Find(&targets)

	data := struct{ Rows []db.Target }{targets}
	return c.Render(http.StatusOK, "templates/home.html", data)
}

func handleAddTargetForm(c echo.Context) error {
	return c.Render(http.StatusOK, "templates/add.html", nil)
}

func handleAddTarget(c echo.Context) error {
	name := c.FormValue("name")
	url := c.FormValue("url")

	db.GetDB().Create(&db.Target{Name: name, URL: url})

	return c.Redirect(http.StatusFound, "/")
}

func handleTargetHistory(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	var target db.Target
	if err := db.GetDB().Preload("History", func(db *gorm.DB) *gorm.DB {
		return db.Order("histories.created_at DESC")
	}).First(&target, id).Error; err != nil {
		return err
	}

	data := struct {
		Target  db.Target
		History []db.History
	}{
		Target:  target,
		History: target.History,
	}
	return c.Render(http.StatusOK, "templates/history.html", data)
}

func handleTargetDelete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	if err := db.GetDB().Delete(&db.Target{}, id).Error; err != nil {
		return err
	}
	return c.String(http.StatusOK, "")
}

func handleToggleActive(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	var target db.Target
	if err := db.GetDB().Preload("History").First(&target, id).Error; err != nil {
		return err
	}

	target.IsActive = !target.IsActive
	db.GetDB().Save(&target)
	return c.String(http.StatusOK, "")
}

func handleHistoryView(c echo.Context) error {
	idStr := c.Param("hid")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	var history db.History
	if err := db.GetDB().First(&history, id).Error; err != nil {
		return err
	}

	diff := make([]difflib.DiffRecord, 0)
	if err := json.Unmarshal([]byte(history.Diff), &diff); err != nil {
		return err
	}
	data := struct {
		History db.History
		Diff    []difflib.DiffRecord
	}{
		History: history,
		Diff:    diff,
	}

	return c.Render(http.StatusOK, "templates/history_view.html", data)
}
