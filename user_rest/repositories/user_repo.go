package repositories

import (
	"encoding/json"
	"user_rest/user_rest/models"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) FindOneRecord(column *string, value *string) (*[]byte, error) {
	var user models.User
	condition := *column + " = ?"
	r.db.Where(condition, value).Find(&user)
	byteData, _ := json.Marshal(user)
	return &byteData, nil
}

func (r *UserRepo) FindAllRecords() (*[]byte, error) {
	var users []models.User
	r.db.Find(&users)
	byteData, _ := json.Marshal(users)
	return &byteData, nil
}

func (r *UserRepo) Save(user interface{}) (*[]byte, error) {
	r.db.Create(user)
	byteData, _ := json.Marshal(user)
	return &byteData, nil
}

func (r *UserRepo) Delete(id string) error {
	var users []models.User
	r.db.Where("id = ?", id).Delete(&users)
	return nil
}

func (r *UserRepo) Update(id string, userUpdate interface{}) (*[]byte, error) {
	r.db.Save(userUpdate)
	byteData, _ := json.Marshal(userUpdate)
	return &byteData, nil
}
