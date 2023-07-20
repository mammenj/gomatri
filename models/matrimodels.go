package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"uniqueIndex"`
	//Ads      []*Ad
}

type Ad struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Section     string `gorm:"not null"`
	Gender      string `gorm:"not null"`
	Religiion   string `gorm:"not null"`
	Cast        string `gorm:"not null"`
	Height      int    `gorm:"not null"`
	Job         string
	JobType     string
	Preferences string
	Phone1      string
	Phone2      string
	Description string `gorm:"not null;type:varchar(255)"`
	Age         int    `gorm:"not null"`
	Education   string `gorm:"not null"`
	Other       string
	Status      int   `gorm:"not null"`
	User        *User `gorm:"omit empty; foreignKey:ID"`
}
