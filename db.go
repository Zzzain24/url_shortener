package main

import (
	"fmt"
	"log"
	"os"
	"time"
    "database/sql"
	"net/http"
	"github.com/joho/godotenv"
    _ "github.com/jackc/pgx/v5/stdlib"
)

// POST function to insert new data into db
func inserLink error () {
	shortUrl := "https://localhost:8080"


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
        panic(err)
    }
	defer db.Close()

	// check if db is connected 
	pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }
    fmt.Println("Connected!")
}





