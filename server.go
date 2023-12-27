package main

import (
	"fmt"
	"html/template"
	"net/http"

	"nausea-admin/internal/storage"
)

type Server struct {
	addr    string
	db      *DB
	t       *template.Template
	storage *storage.Storage
}

type PageData struct {
	PageMeta
	Info Info
}

type PageMeta struct {
	ActiveRoute string
	Title       string
}

func NewServer(addr string, db *DB, t *template.Template, storage *storage.Storage) *Server {
	return &Server{
		addr:    addr,
		db:      db,
		t:       t,
		storage: storage,
	}
}

func pageMeta(r *http.Request) PageMeta {
	return PageMeta{
		ActiveRoute: r.RequestURI,
	}
}

func (s *Server) Run() error {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/", handleHome(s))
	http.HandleFunc("/gallery", handleGallery(s))
	http.HandleFunc("/gallery/upload", handleGalleryUpload(s))

	http.HandleFunc("/update", handleBioUpdate(s))

	fmt.Printf("Listening and serving on %s\n", s.addr)
	return http.ListenAndServe(s.addr, nil)
}

func handleHome(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != "/" {
			s.t.ExecuteTemplate(w, "fallback", PageData{PageMeta: pageMeta(r)})
			return
		}
		info := s.db.GetInfo()
		pd := PageData{
			Info:     info,
			PageMeta: pageMeta(r),
		}
		s.t.ExecuteTemplate(w, "home", pd)
	}
}

func handleGallery(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("landed in gallery!")
		pd := PageData{
			PageMeta: pageMeta(r),
		}
		s.t.ExecuteTemplate(w, "gallery", pd)
	}
}

func handleGalleryUpload(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			fmt.Println("not a post")
			return
		}
		fmt.Println("landed in upload~")
		file, header, err := r.FormFile("file")
		if err != nil {
			panic(err)
		}
		fmt.Println("received file")
		defer file.Close()
		url, _ := s.storage.AddObject(file, header.Filename)

		fmt.Println("created new file")
		fmt.Fprintf(w, "<a href=\"%s\"><img src=\"%s\"></a>", url, url)
	}
}

func handleBioUpdate(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			fmt.Fprint(w, err.Error())
		}
		info := Info{Bio: r.FormValue("bio")}
		s.db.WriteInfo(info)
		s.t.ExecuteTemplate(w, "home", PageData{Info: info})
	}
}

func handleFileUpload(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
