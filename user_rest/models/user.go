package models

import "gorm.io/gorm"

// User
type User struct {
	gorm.Model `json:"-"`
	ID         int    `gorm:"primary_key"`
	Login      string `gorm:"unique_index:idx_login" json:"-"`
	Password   string `json:"-"`
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
	FindAllRecords() (*[]User, error)
	FindOneRecord(column *string, value *string) (*User, error)
	Save(user *User) (*User, error)
	Update(id string, userUpdate *UserUpdate) (*User, error)
	Delete(id string) error
}
