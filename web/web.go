package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/ProjectOrangeJuice/mail-filter/filter"
)

var templates *template.Template

func Start() {
	var err error
	templates, err = template.ParseGlob("web/src/*.html")

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	router := mux.NewRouter()

	router.HandleFunc("/", indexPage)
	router.HandleFunc("/add", addPage)
	router.HandleFunc("/delete", deletePage)
	log.Fatal(http.ListenAndServe(":8000", router))
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Hosts: %v\n", filter.Hosts())
	templates.ExecuteTemplate(w, "index", filter.Hosts())
}

func addPage(w http.ResponseWriter, r *http.Request) {
	h := r.FormValue("host")
	filter.Add(h)

	fmt.Printf("Hosts: %v\n", filter.Hosts())
	templates.ExecuteTemplate(w, "index", filter.Hosts())
}

func deletePage(w http.ResponseWriter, r *http.Request) {
	h := r.FormValue("host")
	filter.Remove(h)

	fmt.Printf("Hosts: %v\n", filter.Hosts())
	templates.ExecuteTemplate(w, "index", filter.Hosts())
}
