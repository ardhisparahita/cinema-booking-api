package seeders

import (
	"errors"
	"log"

	"github.com/ardhisparahita/cinema-booking-api/internal/models"
	"github.com/ardhisparahita/cinema-booking-api/pkg/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) error {
	passwordAdmin := config.Get("ADMIN_PASSWORD")
	if passwordAdmin == "" {
		return errors.New("ADMIN_PASSWORD .env not found")
	}

	var count int64
	db.Model(&models.User{}).Where("email = ?", "admin@gmail.com").Count(&count)
	if count > 0 {
		log.Println("Admin already exist")
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(passwordAdmin),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	admin := models.User{
		Name:     "Admin",
		Email:    "admin@gmail.com",
		Password: string(hashedPassword),
		Role:     "admin",
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}
	return nil
}
