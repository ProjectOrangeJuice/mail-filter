package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/ProjectOrangeJuice/mail-filter/filter"
	"github.com/ProjectOrangeJuice/mail-filter/mail"
)

var templates *template.Template

type data struct {
	Whitelist []string
	Blacklist []string
}

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
	router.HandleFunc("/addBlacklist", addBlacklistPage)
	router.HandleFunc("/deleteBlacklist", deleteBlacklistPage)
	router.HandleFunc("/check", check)
	log.Fatal(http.ListenAndServe(":9090", router))
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	i := data{
		filter.Hosts(),
		filter.Blacklist(),
	}
	fmt.Printf("Hosts: %v\n", i)
	templates.ExecuteTemplate(w, "index", i)
}

func addPage(w http.ResponseWriter, r *http.Request) {
	h := r.FormValue("host")
	filter.Add(h)
	i := data{
		filter.Hosts(),
		filter.Blacklist(),
	}
	templates.ExecuteTemplate(w, "index", i)
}

func check(w http.ResponseWriter, r *http.Request) {

	mail.CheckInbox()
	i := data{
		filter.Hosts(),
		filter.Blacklist(),
	}
	templates.ExecuteTemplate(w, "index", i)
}

func deletePage(w http.ResponseWriter, r *http.Request) {
	h := r.FormValue("host")
	filter.Remove(h)
	i := data{
		filter.Hosts(),
		filter.Blacklist(),
	}
	templates.ExecuteTemplate(w, "index", i)
}

func addBlacklistPage(w http.ResponseWriter, r *http.Request) {
	h := r.FormValue("host")
	filter.AddBlacklist(h)
	i := data{
		filter.Hosts(),
		filter.Blacklist(),
	}
	templates.ExecuteTemplate(w, "index", i)
}

func deleteBlacklistPage(w http.ResponseWriter, r *http.Request) {
	h := r.FormValue("host")
	filter.RemoveBlacklist(h)
	i := data{
		filter.Hosts(),
		filter.Blacklist(),
	}
	templates.ExecuteTemplate(w, "index", i)
}
