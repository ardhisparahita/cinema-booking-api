package models

type MovieGenre struct {
	MovieID uint `gorm:"primaryKey" json:"movie_id"`
	GenreID uint `gorm:"primaryKey" json:"genre_id"`

	Movie Movie `gorm:"foreignKey:MovieID,references:ID" json:"movie,omitempty"`
	Genre Genre `gorm:"foreignKey:GenreID,references:ID" json:"genre,omitempty"`
}

func (mg *MovieGenre) TableName() string {
	return "movie_genres"
}
