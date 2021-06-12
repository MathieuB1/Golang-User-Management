package repositories

import (
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

func (r *UserRepo) FindOneRecord(column *string, value *string) (*models.User, error) {
	var user models.User
	condition := *column + " = ?"
	r.db.Where(condition, value).Find(&user)
	return &user, nil
}

func (r *UserRepo) FindAllRecords() (*[]models.User, error) {
	var users []models.User
	r.db.Find(&users)
	return &users, nil
}

func (r *UserRepo) Save(user *models.User) (*models.User, error) {
	r.db.Create(user)
	return user, nil
}

func (r *UserRepo) Delete(id string) error {
	var users []models.User
	r.db.Where("id = ?", id).Delete(&users)
	return nil
}

func (r *UserRepo) Update(id string, userUpdate *models.User) (*models.User, error) {
	r.db.Save(&userUpdate)
	return userUpdate, nil
}
