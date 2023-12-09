package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type PageData struct {
	Info Info
}

func main() {
	projectID := os.Getenv("PROJECT_ID")
	port := os.Getenv("PORT")

	db := NewDB(projectID)

	server := http.Server{
		Addr: ":" + port,
	}
	t := template.Must(template.ParseFiles("view.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		info := db.GetInfo()
		t.ExecuteTemplate(w, "page", PageData{Info: info})
	})

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			fmt.Fprint(w, err.Error())
		}

		info := Info{Bio: r.FormValue("bio")}
		fmt.Printf("info: %v\n", info)

		db.WriteInfo(info)
		t.ExecuteTemplate(w, "page", PageData{Info: info})
	})

	fmt.Printf("Listening and serving on %s\n", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("FATAL SERVER ERROR: %v\n", err)
	}
}
