package handlers

import (
	"net/http"

	"nausea-admin/internal/db"
	"nausea-admin/internal/models"
	"nausea-admin/internal/server/template"
	"nausea-admin/internal/server/utils"
)

type AboutHandler struct {
	Template template.Template
	DB       db.DB
}

func NewAboutHandler(db db.DB, t template.Template) AboutHandler {
	return AboutHandler{
		DB:       db,
		Template: t,
	}
}

func (h AboutHandler) GetAboutPage(w http.ResponseWriter, r *http.Request) {
	about, err := h.DB.GetAbout()
	if err != nil {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	_, _, pageData := utils.WithPageData(w, r, map[string]interface{}{
		"About": about,
	})
	h.Template.ExecuteTemplate(w, "/about", pageData)
}

func (h AboutHandler) PatchAbout(w http.ResponseWriter, r *http.Request) {
	var patch models.About
	patch.Bio = r.FormValue("bio")
	patch.Update()
	err := h.DB.SetAbout(patch)
	if err != nil {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
