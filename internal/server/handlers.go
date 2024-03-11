package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"nausea-admin/internal/models"
)

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

var links = []struct {
	URL  string
	Name string
}{
	{URL: "/", Name: "Home"},
	{URL: "/folders/", Name: "Folders"},
}

func withPageData(
	w http.ResponseWriter,
	r *http.Request,
	props map[string]interface{},
) (http.ResponseWriter, string, PageData) {
	asideLinks := make([]AsideLink, len(links))
	var title string
	for i, v := range links {
		asideLinks[i] = AsideLink{Name: v.Name, URL: v.URL}
		if (len(v.URL) > 1 && strings.HasPrefix(r.URL.Path, v.URL)) || v.URL == r.URL.Path {
			asideLinks[i].Active = true
			title = v.Name
		}
	}
	return w, r.URL.Path, PageData{Props: props, AsideLinks: asideLinks, Title: title}
}

func GetHomePage(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			_, _, data := withPageData(w, r, map[string]interface{}{})
			w.WriteHeader(http.StatusNotFound)
			s.t.ExecuteTemplate(w, "/404", data)
			return
		}
		s.t.ExecuteTemplate(withPageData(w, r, map[string]interface{}{}))
	}
}

func GetFoldersPage(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		folderID := getFolderID(r)
		folder, err := s.db.GetFolderByID(folderID)
		// TODO interface as struct
		_, _, data := withPageData(
			w,
			r,
			map[string]interface{}{
				"Folder": folder,
				"Error":  err,
			},
		)
		s.t.ExecuteTemplate(w, "/folders", data)
	}
}

func CreateFolder(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parentID := getFolderID(r)
		name := r.FormValue("name")
		if name == "" {
			errorResponse(
				w,
				r,
				http.StatusBadRequest,
				errors.New("insufficient information, missing name"),
			)
			return
		}
		if parentID == "" {
			errorResponse(
				w,
				r,
				http.StatusBadRequest,
				errors.New("insufficient information, missing parent folder id"),
			)
			return
		}
		folder, err := models.NewFolder(parentID, name)
		if err != nil {
			errorResponse(
				w,
				r,
				http.StatusInternalServerError,
				errors.New("failed to create folder"),
			)
			return
		}
		_, asContent, err := s.db.CreateFolder(*folder)
		if err != nil {
			errorResponse(
				w,
				r,
				http.StatusInternalServerError,
				errors.New("failed to create folder"),
			)
			return
		}
		s.t.ExecuteTemplate(w, "folder", asContent)
	}
}

func DeleteFolder(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		folderID := getFolderID(r)
		if folderID == "" {
			w.Header().Set("HX-Reswap", "innerHTML")
			errorResponse(w, r, http.StatusBadRequest, errors.New("no folderID"))
			return
		}
		if folderID == models.RootFolderID {
			w.Header().Set("HX-Reswap", "innerHTML")
			errorResponse(w, r, http.StatusBadRequest, errors.New("can't delete root folder"))
			return
		}
		folder, err := s.db.MarkFolderDeletedByID(folderID)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		asContent, err := folder.AsContent()
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, errors.New("folder is deleted, reload to update UI"))
			return
		}
		fmt.Printf("asContent: %+v\n", asContent.RefID)
		w.WriteHeader(http.StatusOK)
		s.executeTemplate(w, "folder", asContent)
	}
}

func RestoreFolder(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		folderID := getFolderID(r)
		if folderID == "" {
			w.Header().Set("HX-Reswap", "innerHTML")
			errorResponse(w, r, http.StatusBadRequest, errors.New("no folderID"))
			return
		}
		folder, err := s.db.MarkFolderRestoredByID(folderID)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		asContent, err := folder.AsContent()
		fmt.Printf("folder: %+v\n", folder)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, errors.New("folder is restored, reload to update UI"))
			return
		}
		fmt.Printf("asContent: %+v\n", asContent.RefID)
		w.WriteHeader(http.StatusOK)
		s.executeTemplate(w, "folder", asContent)
	}
}

func getFolderID(r *http.Request) string {
	folderID := r.PathValue("id")
	if folderID == "" {
		folderID = models.RootFolderID
	}
	return folderID
}
