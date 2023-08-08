package storage

import (
	"fmt"
	"gomatri/models"
	"log"

	"gorm.io/driver/sqlite"
	//"github.com/glebarez/sqlite"
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
	//as.DB.AutoMigrate(ad)
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

func (as *AdSqlliteStore) GetSection(section string) ([]models.Ad, error) {
	var ads []models.Ad
	log.Println("Get Ads section ", section)
	// result := as.DB.Find(&ads)
	//"name LIKE ?", "%jin%"
	result := as.DB.Where("section = ?", section).Find(&ads)

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
