package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func internalErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Println(err)
}

func main() {
	dbPool, err := pgxpool.New(context.Background(), os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}
	defer dbPool.Close()

	mux := http.NewServeMux()

	html(mux, dbPool)
	mux.Handle("/api/", http.StripPrefix("/api", api(dbPool)))

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../client"))))

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func html(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	tmpl := template.Must(template.ParseFiles("../client/index.html"))

	mux.HandleFunc("GET /{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			internalErr(w, err)
			return
		}

		pageData, err := homeData(dbPool, r.Context(), id)
		if err != nil {
			internalErr(w, err)
			return
		}

		err = tmpl.Execute(w, pageData)
		if err != nil {
			log.Println(err)
		}
	})
}
