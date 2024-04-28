package handlers

import (
	"net/http"

	"nausea-admin/internal/server/template"
	"nausea-admin/internal/server/utils"
)

type HomeHandler struct {
	Template template.Template
}

func NewHomeHandler(t template.Template) HomeHandler {
	return HomeHandler{
		Template: t,
	}
}

func (hh HomeHandler) GetHomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		_, _, data := utils.WithPageData(w, r, map[string]interface{}{})
		w.WriteHeader(http.StatusNotFound)
		hh.Template.ExecuteTemplate(w, "/404", data)
		return
	}
	hh.Template.ExecuteTemplate(utils.WithPageData(w, r, map[string]interface{}{}))
}
