package models

import (
	// "time"

	//"github.com/google/uuid"
	"gorm.io/gorm"
)

// User "Object
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `gorm:"uniqueIndex" json:"email"`
	Password string `json:"password"`
	Role     string `json:"role" gorm:"default:user"`
	Status   string `json:"status" gorm:"default:inactive"`
	Message  string `json:"message" `
}

// BeforeCreate BeforeCreate
/* func (u *User) BeforeCreate(scope *gorm.DB) error {
	//scope.Statement.SetColumn("UpdatedAt", time.Now())
	scope.Statement.SetColumn("UID", uuid.New().String())
	return nil
} */

// BeforeUpdate BeforeUpdate
/* func (u *User) BeforeUpdate(scope *gorm.DB) error {
	scope.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
} */
