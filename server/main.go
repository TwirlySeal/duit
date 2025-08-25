package main

import (
	"context"
	"encoding/json"
	"fmt"
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
	fmt.Println("POSTGRES_URL:", os.Getenv("POSTGRES_URL"))
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("POSTGRES_URL"))
	//dbpool, err := pgxpool.New(context.Background(), "postgres://postgres:admin@localhost:5432/duit?sslmode=disable")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	mux := http.NewServeMux()
	tmpl := template.Must(template.ParseFiles("../client/index.html"))

	mux.HandleFunc("GET /{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		projects, tasks, err := homeData(dbpool, r.Context(), id) // todo: handle err
		if err != nil {
			internalErr(w, err)
			return
		}

		tmpl.Execute(w, PageData{
			projects,
			tasks,
			id,
		}) // todo: handle err
	})

	mux.HandleFunc("GET /data/{id}", func(w http.ResponseWriter, r *http.Request) {
		tasks, err := getTasks(dbpool, r.Context(), r.PathValue("id")) // todo: handle err
		if err != nil {
			internalErr(w, err)
			return
		}

		response, err := json.Marshal(tasks) // todo: handle err
		if err != nil {
			internalErr(w, err)
			return
		}
		w.Write(response) // todo: handle err
	})

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../client"))))

	mux.HandleFunc("POST /addTask", func(w http.ResponseWriter, r *http.Request) {
		var req map[string]interface{}
		//err =json.Unmarshal(r.Body, req)
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		task := req["task"].(string)
		project := req["project"].(string)

		fmt.Println("task", task)
		fmt.Println("project", project)

		w.Write([]byte("Task is" + task))

	})
	log.Fatal(http.ListenAndServe(":8080", mux))
}
