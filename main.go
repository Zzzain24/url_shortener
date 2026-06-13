package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"
	"net/http"
    "github.com/joho/godotenv"
    _ "github.com/jackc/pgx/v5/stdlib"
)

type Application struct {
    db *sql.DB
}

func main() {
	// load .env file and stop the program if there's an error
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	connStr := os.Getenv("DATABASE_URL")
    db, err := sql.Open("pgx", connStr)
    if err != nil {
        log.Fatalf("Failed to open database %v", err)
    }
	defer db.Close()

	// check if db is connected 
	pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }
    fmt.Println("Connected!")

	// stores pointer to db connection 
	app := &Application{db: db}

	// create router
	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", app.handleShorten)
	log.Println("Server starting on :8080...")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}