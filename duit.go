package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type PageData struct {
	Projects []Project
	Tasks []Task
	ActiveId string
}

func home(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	data := PageData{
		projects,
		getTasks(id),
		id,
	}

	tmpl.Execute(w, data)
}

func taskData(w http.ResponseWriter, r *http.Request) {
	tasks := getTasks(r.PathValue("id"))
	response, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response) // todo: handle error
}

var tmpl *template.Template

func main() {
	mux := http.NewServeMux()
	tmpl = template.Must(template.ParseFiles("client/index.html"))

	mux.HandleFunc("GET /{id}", home)
	mux.HandleFunc("GET /data/{id}", taskData)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("client"))))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
