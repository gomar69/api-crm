package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// User struct
type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment")
	}

	// Ambil config dari environment
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	appPort := os.Getenv("APP_PORT")

	if appPort == "" {
		appPort = "8080"
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to DB:", err)
	}
	defer db.Close()

	// Test koneksi DB
	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot ping DB:", err)
	}

	// Route /users
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, email, created_at FROM users")
		if err != nil {
			log.Println("DB query error:", err)
			http.Error(w, "DB query error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		users := []User{}
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
				log.Println("Row scan error:", err)
				http.Error(w, "Row scan error", http.StatusInternalServerError)
				return
			}
			users = append(users, u)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})

	// Optional: route root /
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API CRM running. Endpoint: /users"))
	})

	log.Printf("Server running on port %s", appPort)
	log.Fatal(http.ListenAndServe(":"+appPort, nil))
}
