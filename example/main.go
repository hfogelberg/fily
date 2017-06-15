package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hfogelberg/fily"
	"github.com/urfave/negroni"
)

var tpl *template.Template

func main() {
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/", uploadHandler).Methods("POST")

	mux := http.NewServeMux()
	mux.Handle("/", router)

	static := http.StripPrefix("/public/", http.FileServer(http.Dir("public")))
	router.PathPrefix("/public").Handler(static)

	n := negroni.Classic()
	n.UseHandler(mux)
	http.ListenAndServe(":8080", n)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.New("").ParseFiles("templates/index.html", "templates/layout.html")
	err = tpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Fatalln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	name, err := fily.New(r, 700)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("File save OK! " + name)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
