package storage

import (
	"fmt"
	"log"

	// "gorm.io/driver/sqlite"
	"gomatri/models"

	_ "github.com/glebarez/go-sqlite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AdSqlliteStore struct {
	DB *gorm.DB
}

func NewSqliteAdsStore() *AdSqlliteStore {
	db, err := gorm.Open(sqlite.Open("matri.db"), &gorm.Config{})
	db.AutoMigrate(models.Ad{})
	if err != nil {
		panic("failed to connect database")
	}
	return &AdSqlliteStore{
		DB: db,
	}
}

func (as *AdSqlliteStore) Create(ad *models.Ad) (string, error) {
	log.Println("Before migrate......ADS")
	log.Println("...... After migrate ADS")
	result := as.DB.Create(ad)
	if result.Error != nil {
		return "", result.Error
	}
	log.Println("...... After Create AD..... ID: ", ad.ID)
	return fmt.Sprint(ad.ID), nil
}

func (as *AdSqlliteStore) Get() ([]models.Ad, error) {
	var ads []models.Ad
	log.Println("Get Ads")
	result := as.DB.Find(&ads)

	if result.Error != nil {
		return nil, result.Error
	}
	log.Println("...... Total records ADS : ", result.RowsAffected)
	return ads, nil
}

func (as *AdSqlliteStore) GetSection(
	section string,
	offset int,
) ([]models.Ad, error) {
	var ads []models.Ad
	log.Println("Get Ads section ", section)
	// db.Limit(10).Offset(5).Find(&users)
	result := as.DB.Limit(10).
		Offset(offset).
		Where("section = ?", section).
		Find(&ads)

	if result.Error != nil {
		return nil, result.Error
	}
	log.Println("...... Total records ADS : ", result.RowsAffected)
	return ads, nil
}

func (as *AdSqlliteStore) Delete(id string) (string, error) {
	log.Println("Delete AD ID: ", id)
	result := as.DB.Delete(&models.Ad{}, id)

	if result.Error != nil {
		return "", result.Error
	}
	log.Println("...... Total deleted records : ", result.RowsAffected)
	return id, nil
}

func (as *AdSqlliteStore) Update(ad *models.Ad) (uint, error) {
	log.Println("Update Users ID: ", ad.ID)
	result := as.DB.Updates(ad)

	if result.Error != nil {
		return 0, result.Error
	}
	log.Println("...... Total Updated records : ", result.RowsAffected)
	return ad.ID, nil
}

func (as *AdSqlliteStore) GetOne(id string) (*models.Ad, error) {
	log.Println("Get Users ID: ", id)
	var ad *models.Ad
	result := as.DB.First(&ad, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return ad, nil
}
