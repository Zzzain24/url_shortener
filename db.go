package main

import (
	"fmt"
	"log"
	"os"
	"math/rand/v2"
    "database/sql"
	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// character set for short_code
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateShortCode () string {
	var code []byte
	for i := 0; i < 6; i++ {
		randomNum := rand.IntN(len(charset))
		code = append(code, charset[randomNum]) 
	}

	return string(code)
}

func InsertLink (db *sql.DB, originalURL string) (string, error) {
	// flag to check if generated short_code exists in db or not 
	var exists bool
	var shortCode string
	for {
		shortCode = GenerateShortCode()
		err := db.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM url_data WHERE short_code = $1)",
			shortCode,
		).Scan(&exists)
		// ensure query is clean
		if err != nil {
			return "", err
		}
		// exit loop if valid short_code is generated
		if !exists {
			break
		}
	}
	// insert shortcode and url data
	query := "INSERT INTO url_data (short_code, original_url) VALUES ($1, $2)"
	_, err := db.Exec(query, shortCode, originalURL)
	if err != nil{
		return "", err
	}
	return shortCode, nil
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





