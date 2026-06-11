package main 

import (
	"time"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type LinkResponse struct {
	ShortCode string `json:"short_code"`
	ShortUrl string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	Clicks int `json:"clicks"`
	CreatedAt time.Time `json:"created_at"`
}

type Link struct {
	ShortCode: string
	OriginalURL: string
	Clicks: int
	CreatedAt: time.Time
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Application struct {
    db *sql.DB
}



