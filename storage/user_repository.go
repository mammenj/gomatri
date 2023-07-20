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

func NewSqliteUserStore() *UserSqlliteStore {
	db, err := gorm.Open(sqlite.Open("matri.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(models.User{})
	return &UserSqlliteStore{
		DB: db,
	}
}

func (us *UserSqlliteStore) Create(mu *models.User) (string, error) {
	log.Println("Before migrate......")
	//us.DB.AutoMigrate(mu)
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

func (us *UserSqlliteStore) Delete(id string) (string, error) {
	log.Println("Delete Users ID: ", id)
	result := us.DB.Delete(&models.User{}, id)

	if result.Error != nil {
		return "", result.Error
	}
	log.Println("...... Total deleted records : ", result.RowsAffected)
	return id, nil
}

func (us *UserSqlliteStore) Update(mu *models.User) (uint, error) {
	log.Println("Update Users ID: ", mu.ID)
	result := us.DB.Updates(mu)

	if result.Error != nil {
		return 0, result.Error
	}
	log.Println("...... Total Updated records : ", result.RowsAffected)
	return mu.ID, nil
}

func (us *UserSqlliteStore) GetOne(id string) (*models.User, error) {
	log.Println("Get Users ID: ", id)
	var user *models.User
	result := us.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
