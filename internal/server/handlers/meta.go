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

type MetaHandler struct {
	DB       db.DB
	Storage  *storage.Storage
	Template template.Template
}

func NewMetaHandler(
	DB db.DB,
	Storage *storage.Storage,
	Template template.Template,
) *MetaHandler {
	return &MetaHandler{
		DB:       DB,
		Storage:  Storage,
		Template: Template,
	}
}

func (h *MetaHandler) GetPage(w http.ResponseWriter, r *http.Request) {
	meta, err := h.DB.GetMeta()
	if err != nil {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	_, _, pageData := utils.WithPageData(
		w,
		r,
		map[string]interface{}{"Meta": meta},
	)
	h.Template.ExecuteTemplate(w, "/meta", pageData)
}

func (h *MetaHandler) PutMeta(w http.ResponseWriter, r *http.Request) {
	err := parseMultipartForm(r)
	if err != nil {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	file := r.MultipartForm.File["background-image"]
	urls, errs := filesIntoBucket(file, h.Storage.AddObject)
	if len(errs) > 0 {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusInternalServerError, errs[0])
		return
	}
	media, err := models.NewMedia(urls[0].URL, urls[0].MediaSize)
	if err != nil {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	prev := r.MultipartForm.Value["prev-image-url"][0]
	if prev != "" {
		err := h.Storage.RemoveObject(h.Storage.ParseURLKey(prev))
		if err != nil {
			fmt.Println("> Failed to remove prev image", err)
		}
	}
	meta := models.Meta{Background: media}
	meta.Update()
	h.DB.SetMeta(meta)
	w.WriteHeader(http.StatusOK)
	h.Template.ExecuteTemplate(w, "prev-meta-image-url", map[string]interface{}{
		"Props": map[string]interface{}{
			"Meta": meta,
		},
	})
}
