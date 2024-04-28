package handlers

import (
	"fmt"
	"net/http"

	"nausea-admin/internal/db"
	"nausea-admin/internal/models"
	"nausea-admin/internal/server/template"
	"nausea-admin/internal/server/utils"
	"nausea-admin/internal/storage"
)

type AboutHandler struct {
	Template template.Template
	DB       *db.DB
	Storage  *storage.Storage
}

func NewAboutHandler(db *db.DB, t template.Template, s *storage.Storage) AboutHandler {
	return AboutHandler{
		DB:       db,
		Template: t,
		Storage:  s,
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
	err := parseMultipartForm(r)
	if err != nil {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	var patch models.About
	patch.Bio = r.MultipartForm.Value["bio"][0]
	file := r.MultipartForm.File["image"]
	if len(file) == 1 {
		url, errs := filesIntoBucket(file, h.Storage.AddObject)
		if len(errs) != 0 {
			w.Header().Set("HX-Reswap", "innerHTML")
			utils.ErrorResponse(w, r, http.StatusInternalServerError, errs[0])
			return
		}
		img, err := models.NewMedia(url[0].URL, url[0].MediaSize)
		if err != nil {
			w.Header().Set("HX-Reswap", "innerHTML")
			utils.ErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		patch.Image = &img
	}
	patch.Update()
	prevUrl := r.MultipartForm.Value["prev-image-url"][0]
	if prevUrl != "" {
		err := h.Storage.RemoveObject(h.Storage.ParseURLKey(prevUrl))
		if err != nil {
			fmt.Println("Failed to remove prev image", err)
		}
	}
	err = h.DB.SetAbout(patch)
	if err != nil {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	h.Template.ExecuteTemplate(w, "prev-about-image-url", map[string]interface{}{
		"Props": map[string]interface{}{
			"About": patch,
		},
	})
}
