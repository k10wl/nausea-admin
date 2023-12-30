package server

import (
	"fmt"
	"net/http"

	"nausea-admin/internal/models"
)

func handleAboutPage(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != "/" {
			s.t.ExecuteTemplate(w, "/404", PageData{PageMeta: pageMeta(r)})
			return
		}
		info, err := s.db.GetInfo()
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError)
			return
		}
		pd := PageData{
			Info:     *info,
			PageMeta: pageMeta(r),
		}
		s.executeTemplate(w, r.URL.Path, pd)
	}
}

func handleGalleryPage(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pd := PageData{
			PageMeta: pageMeta(r),
		}
		s.executeTemplate(w, r.URL.Path, pd)
	}
}

func handleGalleryUpload(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("file")
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest)
			return
		}
		defer file.Close()
		url, err := s.storage.AddObject(file, header.Filename)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError)
			return
		}
		// TODO: replace with template after defined how it will look
		fmt.Fprintf(w, "<a href=\"%s\"><img src=\"%s\"></a>", url, url)
	}
}

// XXX actually this is about update
func handleBioUpdate(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest)
			return
		}
		info := models.Info{Bio: r.FormValue("bio")}
		err = s.db.WriteInfo(info)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError)
			return
		}
		pd := PageData{Info: info}
		s.executeTemplate(w, r.URL.Path, pd)
	}
}
