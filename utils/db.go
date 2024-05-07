package utils

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDb(cfg *AcmeConfig) (*gorm.DB, error) {
	var err error
	var db *gorm.DB

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host,
		cfg.Username,
		cfg.Password,
		cfg.DBName,
		cfg.Port,
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
}
