package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	server := http.Server{
		Addr: ":8081",
	}
	t := template.Must(template.ParseFiles("view.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		t.ExecuteTemplate(w, "page", map[string]interface{}{})
	})

	fmt.Printf("Listening and serving on %s", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("FATAL SERVER ERROR: %v", err)
	}
}
