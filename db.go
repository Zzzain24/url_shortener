package main

import (
    "database/sql"
    "errors"
    "math/rand/v2"
	"time"
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

func InsertLink (db *sql.DB, originalURL string) (string, time.Time, error) {
	// flag to check if generated short_code exists in db or not 
	var exists bool
	var shortCode string
	var createdAt time.Time 
	for {
		shortCode = GenerateShortCode()
		query := "SELECT EXISTS(SELECT 1 FROM url_data WHERE short_code = $1)"
		err := db.QueryRow(query, shortCode).Scan(&exists)
		// ensure query is clean
		if err != nil {
			return "", time.Time{}, err
		}
		// exit loop if valid short_code is generated
		if !exists {
			break
		}
	}
	// insert shortcode and url data
	query := "INSERT INTO url_data (short_code, original_url) VALUES ($1, $2) RETURNING created_at"
	err := db.QueryRow(query, shortCode, originalURL).Scan(&createdAt)
	if err != nil{
		return "", time.Time{}, err
	}
	return shortCode, createdAt, nil
}

func GetOriginalURL (db *sql.DB, shortCode string) (string, error) {
	var originalURL string
	// use shortCode to query db for originalURL 
	query := "SELECT original_url FROM url_data WHERE short_code = $1"
	err := db.QueryRow(query, shortCode).Scan(&originalURL)

	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNotFound
	}

	if err != nil {
		return "", err
	}

	return originalURL, nil
}

func IncrementClicks (db *sql.DB, shortCode string) error {
	query := "UPDATE url_data SET clicks = clicks + 1 WHERE short_code = $1"
	_, err := db.Exec(query, shortCode)
	if err != nil {
		return err
	}
	return nil	
}

func GetStats (db *sql.DB, shortCode string) (Link, error) {
	// return Link object with short_code, original_url, short_url, clicks, created_at
	var linkResponse Link
	linkResponse.ShortCode = shortCode
	query := "SELECT original_url, clicks, created_at FROM url_data WHERE short_code = $1"
	err := db.QueryRow(query, shortCode).Scan(&linkResponse.OriginalURL, &linkResponse.Clicks, &linkResponse.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return Link{}, ErrNotFound
	}

	if err != nil {
		return Link{}, err
	}

	return linkResponse, nil
}





