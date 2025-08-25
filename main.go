package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

//-update23

func main() {
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	appPort := os.Getenv("APP_PORT")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	http.HandleFunc("/db-check", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			http.Error(w, "DB not connected", 500)
			return
		}
		w.Write([]byte("DB connected"))
	})

	log.Println("Server running on port " + appPort)
	http.ListenAndServe(":"+appPort, nil)
}
