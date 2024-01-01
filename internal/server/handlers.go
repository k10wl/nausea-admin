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

type EmailTemplate string

type UrlTemplate struct {
	ID   string
	Text string
	URL  string
}

type PageData struct {
	PageMeta
	Forms []FormTemplate
	Links []UrlTemplate
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
		bio, err := s.db.GetBio()
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		s.executeTemplate(
			w,
			"form",
			FormTemplate{Name: "Bio", FormSubmitURL: "/about/update", Value: bio},
		)
	}
}

func handleAboutUpdate(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}
		bio := r.FormValue("Bio")
		err = s.db.SetBio(bio)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		s.executeTemplate(
			w,
			"form",
			FormTemplate{Name: "Bio", FormSubmitURL: r.URL.Path, Value: bio},
		)
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
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}
		defer file.Close()
		url, err := s.storage.AddObject(file, header.Filename)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		// TODO: replace with template after defined how it will look
		fmt.Fprintf(w, "<a href=\"%s\"><img src=\"%s\"></a>", url, url)
	}
}

func handleContactsPage(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pd := PageData{
			PageMeta: pageMeta(r),
			Lazy:     []LazyLoad{{"/contacts/data"}},
		}
		s.executeTemplate(w, r.URL.Path, pd)
	}
}

func handleContactsData(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var email string
		var links []models.Link
		var err error
		emailChan := make(chan string)
		linksChan := make(chan []models.Link)
		errChan := make(chan error)
		go func() {
			e, err := s.db.GetEmail()
			if err != nil {
				errChan <- err
				return
			}
			emailChan <- e
		}()

		go func() {
			l, err := s.db.GetLinks()
			if err != nil {
				errChan <- err
				return
			}
			linksChan <- l
		}()
	outer:
		for i := 0; i < 2; i++ {
			select {
			case e := <-emailChan:
				email = e
			case l := <-linksChan:
				links = l
			case e := <-errChan:
				err = e
				break outer

			}
		}
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
		}
		urlTemplates := []UrlTemplate{}
		for _, v := range links {
			urlTemplates = append(urlTemplates, UrlTemplate{URL: v.URL, Text: v.Text, ID: v.ID})
		}
		pd := PageData{
			Forms: []FormTemplate{
				{
					Name:          "Email",
					Value:         email,
					FormSubmitURL: "/contacts/update/email",
				},
			},
			Links: urlTemplates,
		}
		s.executeTemplate(w, "contacts lazy", pd)
	}
}

func handleEmailUpdate(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}
		email := r.FormValue("Email")
		err = s.db.SetEmail(email)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		pd := FormTemplate{
			Name:          "Email",
			Value:         email,
			FormSubmitURL: "/contacts/update",
		}
		s.executeTemplate(w, "form", pd)
	}
}

func handleLinkUpdate(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}
		link := models.Link{
			ID:   r.FormValue("ID"),
			URL:  r.FormValue("URL"),
			Text: r.FormValue("Text"),
		}
		link, err = s.db.SetLink(link)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		pd := UrlTemplate{
			ID:   link.ID,
			URL:  link.URL,
			Text: link.Text,
		}
		s.executeTemplate(w, "contacts link", pd)
	}
}

func handleLinkDelete(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}
		link := models.Link{
			ID: r.FormValue("ID"),
		}

		err = s.db.DeleteLink(link.ID)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}
