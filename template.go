package router

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)


type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func Templates() *Template {
    var files []string

    err := filepath.Walk("pages", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if strings.HasSuffix(path, ".go.html") {
            files = append(files, path)
        }
        return nil
    })
    if err != nil {
        panic(err)
    } 

    tmpl := template.New("")
    _, err = tmpl.ParseFiles(files...)
    if err != nil {
        panic(err)
    }

    return &Template{
        templates: tmpl,
    }
}

