package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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
	file, header, err := r.FormFile("file")

	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	out, err := os.Create("./public/tmp/" + header.Filename)
	if err != nil {
		log.Println(err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		log.Println(err)
	}

	log.Println(header.Filename)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
