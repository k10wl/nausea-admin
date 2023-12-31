package server

import (
	"fmt"
	"net/http"

	"nausea-admin/internal/models"
)

type LazyLoad struct {
	URL string
}

type FormTemplate struct {
	Name          string
	FormSubmitURL string
	Value         string
}

type PageData struct {
	PageMeta
	Forms []FormTemplate
	Lazy  []LazyLoad
}

type PageMeta struct {
	ActiveRoute string
	Title       string
}

func handleAboutPage(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != "/" {
			s.t.ExecuteTemplate(w, "/404", PageData{PageMeta: pageMeta(r)})
			return
		}
		pd := PageData{
			PageMeta: pageMeta(r),
			Lazy:     []LazyLoad{{"/about/bio"}},
		}
		s.executeTemplate(w, r.URL.Path, pd)
	}
}

func handleAboutBio(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info, err := s.db.GetInfo()
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError)
			return
		}
		s.executeTemplate(w, "form", FormTemplate{Name: "Bio", FormSubmitURL: "/about/update", Value: info.Bio})
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
func handleAboutUpdate(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest)
			return
		}
		info := models.Info{Bio: r.FormValue("Bio")}
		err = s.db.WriteInfo(info)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError)
			return
		}
		fmt.Printf("info.Bio: %v\n", info.Bio)
		s.executeTemplate(
			w,
			"form",
			FormTemplate{Name: "Bio", FormSubmitURL: r.URL.Path, Value: info.Bio},
		)
	}
}
