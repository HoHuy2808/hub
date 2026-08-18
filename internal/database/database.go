package database

import (
	"hub/internal/database/entities"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToPostgreSQL() (*gorm.DB, error) {
	var err error
	DB, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	err = DB.AutoMigrate(
		&entities.User{},
		&entities.UserProfile{},
		&entities.RefreshToken{},
		&entities.Post{},
		&entities.PostAttachment{},
		&entities.Comment{},
		&entities.Reaction{},
		&entities.Request{},
		&entities.Contact{},
		&entities.Notification{},
	)
	if err != nil {
		panic("Failed to auto migrate: " + err.Error())
	}
	return DB, nil
}
