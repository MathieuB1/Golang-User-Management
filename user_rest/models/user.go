package models

import "gorm.io/gorm"

// User
type User struct {
	gorm.Model `json:"-"`
	ID         int    `gorm:"primary_key"`
	Login      string `gorm:"unique_index:idx_login"`
	Password   string
	First_name string
	Last_name  string
	Email      string `gorm:"unique_index:idx_email"`
}

// Format User
type UserUpdate struct {
	ID         string
	Login      string `json:"-"`
	Password   string `json:"-"`
	First_name string
	Last_name  string
	Email      string
}

// UserRepository
type UserRepository interface {
	FindAllRecords() (*[]byte, error)
	FindOneRecord(column *string, value *string) (*[]byte, error)
	Delete(id string) error
	Save(objToSave interface{}) (*[]byte, error)
	Update(id string, objToUpdate interface{}) (*[]byte, error)
}
