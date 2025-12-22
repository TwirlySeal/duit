package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/BurntSushi/toml"
)

type Config struct {
	PostgresUrl string `toml:"postgres_url"`
}

func main() {
	var config Config
	_, err := toml.DecodeFile(".env.toml", &config)
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	dbPool, err := pgxpool.New(context.Background(), config.PostgresUrl)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}
	defer dbPool.Close()

	mux := http.NewServeMux()

	html(mux, dbPool)
	mux.Handle("/api/", http.StripPrefix("/api", api(dbPool)))

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func html(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	tmpl := template.Must(template.ParseFiles("../client/home.html"))

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

		w.Header().Set("Content-Type", "text/html")

		encoder := selectEncoding(w, r)
		if encoder == nil {
			tmpl.Execute(w, pageData)
		} else {
			tmpl.Execute(encoder, pageData)
			encoder.Close()
		}
	})
}
