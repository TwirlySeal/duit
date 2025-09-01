package main

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"html/template"
	"log"
	"net/http"
	"os"
)

type PageData struct {
	Projects []Project
	Tasks    []Task
	ActiveId string
}

type NewTask struct {
	Name      string
	ProjectId string
}

type NewProject struct {
	Name string
}

func internalErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Println(err)
}

func main() {
	dbPool, err := pgxpool.New(context.Background(), os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbPool.Close()

	mux := http.NewServeMux()
	tmpl := template.Must(template.ParseFiles("../client/index.html"))

	mux.HandleFunc("GET /{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		projects, tasks, err := homeData(dbPool, r.Context(), id)
		if err != nil {
			internalErr(w, err)
			return
		}

		err = tmpl.Execute(w, PageData{
			projects,
			tasks,
			id,
		})
		if err != nil {
			log.Println(err)
		}
	})

	mux.HandleFunc("GET /data/{id}", func(w http.ResponseWriter, r *http.Request) {
		tasks, err := getTasks(dbPool, r.Context(), r.PathValue("id"))
		if err != nil {
			internalErr(w, err)
			return
		}

		response, err := json.Marshal(tasks)
		if err != nil {
			internalErr(w, err)
			return
		}

		_, err = w.Write(response)
		if err != nil {
			log.Println(err)
			return
		}
	})

	mux.HandleFunc("POST /addTask", func(w http.ResponseWriter, r *http.Request) {
		var task NewTask
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			internalErr(w, err)
			return
		}

		err = addTask(dbPool, r.Context(), task)
		if err != nil {
			internalErr(w, err)
			return
		}

		w.Header().Set("Location", "/"+task.ProjectId+"/"+task.Name)
		w.WriteHeader(http.StatusCreated)
	})

	mux.HandleFunc("POST /addProject", func(w http.ResponseWriter, r *http.Request) {
		var project NewProject
		err := json.NewDecoder(r.Body).Decode(&project)
		if err != nil {
			internalErr(w, err)
			return
		}

		id, err := addProject(dbPool, r.Context(), project.Name)
		if err != nil {
			internalErr(w, err)
			return
		}

		w.Header().Set("Location", "/"+id)
		w.WriteHeader(http.StatusCreated)
	})

	fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir("../client")))
	mux.Handle("GET /static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache")
		fileServer.ServeHTTP(w, r)
	}))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
