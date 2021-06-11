package repositories

import (
	"api-test/web_rest/models"

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

func (r *UserRepo) Update(id string, userUpdate *models.UserUpdate) (*models.User, error) {

	var existingUser models.User
	r.db.Where("id = ?", id).Find(&existingUser)

	if userUpdate.Email != "" {
		existingUser.Email = userUpdate.Email
	}
	if userUpdate.Login != "" {
		existingUser.Login = userUpdate.Login
	}
	if userUpdate.Password != "" {
		existingUser.Password = userUpdate.Password
	}
	if userUpdate.First_name != "" {
		existingUser.First_name = userUpdate.First_name
	}
	if userUpdate.Last_name != "" {
		existingUser.Last_name = userUpdate.Last_name
	}

	r.db.Save(&existingUser)

	return &existingUser, nil
}
