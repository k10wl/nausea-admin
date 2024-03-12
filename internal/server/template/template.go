package template

import (
	"fmt"
	"html/template"
	"net/http"

	"nausea-admin/internal/server/logger"
)

type Template struct {
	t *template.Template
	l logger.ServerLogger
}

type AsideLink struct {
	Active bool
	URL    string
	Name   string
}

type PageData struct {
	Props      interface{}
	AsideLinks []AsideLink
	Title      string
}

func NewTemplate(l logger.ServerLogger) Template {
	t := template.Must(template.ParseGlob("views/**"))
	return Template{t: t, l: l}
}

func (t Template) ExecuteTemplate(w http.ResponseWriter, tmpl string, data any) {
	err := t.t.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		t.sendError(w, tmpl, err)
	}
}

func (t Template) sendError(w http.ResponseWriter, tmpl string, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	t.l.Logln(fmt.Sprintf("Error upon executing template \"%s\": %+v", tmpl, err))
	t.t.ExecuteTemplate(w, "_error", map[string]interface{}{"Error": err})
}
