package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func internalErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Println(err)
}

func writeJson(w http.ResponseWriter, r *http.Request, v any) {
	response, err := json.Marshal(v)
	if err != nil {
		internalErr(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	encoder := selectEncoding(w, r)
	if encoder == nil {
		w.Write(response)
	} else {
		encoder.Write(response)
		encoder.Close()
	}
}

func api(dbPool *pgxpool.Pool) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /projects/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		tasks, err := getTasks(dbPool, r.Context(), id)
		if err != nil {
			internalErr(w, err)
			return
		}

		writeJson(w, r, tasks)
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
		writeJson(w, r, struct {
			Id int `json:"id"`
		}{id})
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
