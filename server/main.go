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

func internalErr(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
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
			internalErr(w, err)
			return
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
			internalErr(w, err)
			return
		}
	})

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../client"))))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
