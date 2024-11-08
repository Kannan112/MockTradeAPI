package db

import (
	"fmt"
	"log"

	"github.com/kannan112/mock-trading-platform-api/pkg/config"
	"github.com/kannan112/mock-trading-platform-api/pkg/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// func to connect data base using config(database config) and return address of a new instnce of gorm DB
func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DBPassword)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		return nil, err
	}

	// migrate the database tables
	err = db.AutoMigrate(&domain.User{}, &domain.Account{}, &domain.Order{}, &domain.Trade{}, &domain.Position{})

	if err != nil {
		log.Printf("failed to migrate database models")
		return nil, err
	}

	return db, err
}
