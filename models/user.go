package models

import (
	// "time"

	//"github.com/google/uuid"
	"gorm.io/gorm"
)

// User "Object
type User struct {
	gorm.Model
	//UID      uuid.UUID `json:"uuid"`
	Username string    `json:"username" binding:"required"`
	Password string    `json:"password"`
	Message  string    `json:"message"`
	Email    string    `gorm:"uniqueIndex" json:"email"`

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

