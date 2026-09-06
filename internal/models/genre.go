package models

type Genre struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"size:50;not null" json:"name"`
}

func (g *Genre) TableName() string {
	return "genres"
}
