package storage

import (
	"fmt"
	"gomatri/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UserSqlliteStore struct {
	DB *gorm.DB
}

type AdSqlliteStore struct {
	DB *gorm.DB
}

func NewSqliteAdsStore() *AdSqlliteStore {

	db, err := gorm.Open(sqlite.Open("matri.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return &AdSqlliteStore{
		DB: db,
	}
}

func NewSqliteUserStore() *UserSqlliteStore {

	db, err := gorm.Open(sqlite.Open("matri.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return &UserSqlliteStore{
		DB: db,
	}
}

func (us *UserSqlliteStore) Create(mu *models.User) (string, error) {
	log.Println("Before migrate......")
	us.DB.AutoMigrate(mu)
	log.Println("...... After migrate")
	result := us.DB.Create(mu)
	if result.Error != nil {
		return "", result.Error
	}
	log.Println("...... After Create ..... ID: ", mu.ID)
	return fmt.Sprint(mu.ID), nil
}

func (us *UserSqlliteStore) Get() ([]models.User, error) {

	var users []models.User
	log.Println("Get Users")
	result := us.DB.Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}
	log.Println("...... Total records : ", result.RowsAffected)
	return users, nil
}
