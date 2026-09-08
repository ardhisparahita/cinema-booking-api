package models

import "time"

type Movie struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `gorm:"size:200'not null" json:"title"`
	Synopsis    string    `gorm:"type:text" json:"synopsis"`
	DurationMin uint16    `gorm:"not null" json:"duration_min"`
	Rating      string    `gorm:"type:enum('su','13+','17+','21+');not null" json:"rating"`
	PosterUrl   string    `gorm:"size:255" json:"poster_url"`
	ReleaseDate time.Time `gorm:"type:date" json:"release_date"`
	IsActive    bool      `gorm:"not null;default=true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	MovieGenres []MovieGenre `gorm:"foreignKey:MovieID" json:"movie_genres,omitempty"`
}

func (m *Movie) TableName() string {
	return "movies"
}
