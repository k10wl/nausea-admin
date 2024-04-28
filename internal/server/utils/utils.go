package utils

import (
	"log"
	"net/http"
	"slices"
	"strings"

	"nausea-admin/internal/server/template"
)

type link struct {
	URL  string
	Name string
}

var links = []link{
	{URL: "/", Name: "Home"},
	{URL: "/meta/", Name: "Meta"},
	{URL: "/about/", Name: "About"},
	{URL: "/contacts/", Name: "Contacts"},
	{URL: "/folders/", Name: "Folders"},
	{URL: "/folders/--CAROUSEL--", Name: "Carousel"},
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, code int, e error) {
	log.Printf(
		"Error in request: %s %s --- %d: %s --- Error: %s",
		r.Method,
		r.URL.Path,
		code,
		http.StatusText(code),
		e,
	)
	msg := e.Error()
	if msg == "" {
		msg = http.StatusText(code)
	}
	log.Println(msg)
	http.Error(w, msg, code)
}

func WithPageData(
	w http.ResponseWriter,
	r *http.Request,
	props map[string]interface{},
) (http.ResponseWriter, string, template.PageData) {
	asideLinks := make([]template.AsideLink, len(links))
	title, titleExists := props["Title"].(string)
	for i, v := range links {
		asideLinks[i] = template.AsideLink{Name: v.Name, URL: v.URL}
		if (len(v.URL) > 1 &&
			strings.HasPrefix(r.URL.Path, v.URL) &&
			slices.IndexFunc(links, func(e link) bool {
				return e.URL == r.URL.Path
			}) == -1) ||
			v.URL == r.URL.Path {
			asideLinks[i].Active = true
			if !titleExists {
				title = v.Name
			}
		}
	}
	return w, r.URL.Path, template.PageData{
		Props:      props,
		AsideLinks: asideLinks,
		Title:      title,
	}
}
