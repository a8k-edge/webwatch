package server

import (
	"embed"
	"fmt"
	"html/template"
	"io"

	"webwatch/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed templates
var tmplFS embed.FS

type Renderer struct{}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, err := template.ParseFS(tmplFS, "templates/base.html", name)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}

func Serve() {
	e := echo.New()

	e.Renderer = &Renderer{}

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", handleHome)
	e.GET("/add", handleAddTargetForm)
	e.POST("/add", handleAddTarget)
	e.GET("/target/:id", handleTargetHistory)
	e.DELETE("/target/:id", handleTargetDelete)
	e.POST("/target/:id/toggle", handleToggleActive)
	e.GET("/target/:tid/history/:hid", handleHistoryView)

	addr := fmt.Sprintf("localhost:%s", config.Cfg.Server.Port)
	fmt.Println("Starting UI server", addr)
	e.Logger.Fatal(e.Start(addr))
}
