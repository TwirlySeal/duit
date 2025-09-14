package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func api(dbPool *pgxpool.Pool) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /projects/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		tasks, err := getTasks(dbPool, r.Context(), id)
		if err != nil {
			internalErr(w, err)
			return
		}

		response, err := json.Marshal(tasks)
		if err != nil {
			internalErr(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(response)
		if err != nil {
			log.Println(err)
		}
	})

	mux.HandleFunc("POST /tasks", func(w http.ResponseWriter, r *http.Request) {
		var task NewTask
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			internalErr(w, err)
			return
		}

		id, err := addTask(dbPool, r.Context(), task)
		if err != nil {
			internalErr(w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		body, err := json.Marshal(struct {
			Id int `json:"id"`
		}{id})
		if err != nil {
			internalErr(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(body)
		if err != nil {
			log.Println(err)
		}
	})

	mux.HandleFunc("PATCH /tasks", func(w http.ResponseWriter, r *http.Request) {
		var patch TaskPatch
		err := json.NewDecoder(r.Body).Decode(&patch)
		if err != nil {
			internalErr(w, err)
			return
		}
		err = patchTask(dbPool, r.Context(), patch)
		if err != nil {
			internalErr(w, err)
			return
		}
	})

	return mux
}
