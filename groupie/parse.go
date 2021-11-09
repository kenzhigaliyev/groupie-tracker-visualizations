package student

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func Err(Str string, Status int, w http.ResponseWriter, r *http.Request) {

	Info := Error{Str, Status}
	val, err := template.ParseFiles("templates/err.html")

	if err != nil {
		Err("500 Internal Server Error", http.StatusInternalServerError, w, r)
		return
	}
	w.WriteHeader(Status)
	err = val.Execute(w, Info)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
}

func Media(w http.ResponseWriter, r *http.Request) {

	if !Result {
		Err("500 Internal Server Error", http.StatusInternalServerError, w, r)
		return
	}

	if r.URL.Path != "/" && r.Method == "GET" {
		Err("404 Not Found", http.StatusNotFound, w, r)
		return
	}

	if r.URL.Path == "/" && r.Method != "GET" {
		Err("404 Not Found", http.StatusNotFound, w, r)
		return
	}

	val, err := template.ParseFiles("templates/groupie.html")
	if err != nil {
		Err("500 Internal Server Error", http.StatusInternalServerError, w, r)
		return
	}

	err = val.Execute(w, ArtistsNew)
	if err != nil {
		Err("500 Internal Server Error", http.StatusInternalServerError, w, r)
	}

}

func Album(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/artists/" && r.Method == "POST" {
		Err("404 Not Found", http.StatusNotFound, w, r)
		return
	}

	if r.Method != "POST" && r.URL.Path == "/ascii-art/" {
		Err("405 Method Not Allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	val, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		Err("500 Internal Server Error", http.StatusInternalServerError, w, r)
		return
	}
	name := strings.TrimPrefix(r.URL.Path, "/artists/")
	id, smt := strconv.Atoi(name)
	if smt != nil {
		Err("404 Not Found", http.StatusNotFound, w, r)
		return
	}
	if (len(RelationNew.Index) < id) || (id < 1) {
		Err("404 Not Found", http.StatusNotFound, w, r)
		return
	}
	ArtistsNew[id-1].DatesLocations = RelationNew.Index[id-1].DatesLocations
	err = val.Execute(w, ArtistsNew[id-1])
	if err != nil {
		Err("500 Internal Server Error", http.StatusInternalServerError, w, r)
		return
	}
}

func MainFunc() {
	val := http.FileServer(http.Dir("style"))
	http.Handle("/style/", http.StripPrefix("/style", val))
	Func()
	http.HandleFunc("/", Media)
	http.HandleFunc("/artists/", Album)
	http.ListenAndServe(":7770", nil)
}
