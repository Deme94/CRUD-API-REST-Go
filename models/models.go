package models

import (
	"database/sql"
	"time"
)

// Models is the wrapper for database
type Models struct {
	DB DBModels
}

// NewModels returns models with db pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModels{DB: db},
	}
}

// Game is the type for game
type Game struct {
	ID          int            `json:"id"`
	Title       string         `json:"title"`
	ImageUrl    string         `json:"image_url"`
	Genres      map[int]string `json:"genres"`
	Modes       map[int]string `json:"modes"`
	Developers  []string       `json:"developers"`
	Publishers  []string       `json:"publishers"`
	ReleaseDate time.Time      `json:"release_date"`
	Storage     int            `json:"storage"`
	Likes       int            `json:"likes"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
}

// Genre is the type for genre
type Genre struct {
	ID        int       `json:"id"`
	GenreName string    `json:"genre_name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// Mode is the type for gamemode
type Mode struct {
	ID        int       `json:"id"`
	ModeName  string    `json:"mode_name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
