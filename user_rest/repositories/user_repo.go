package repositories

import (
	"encoding/json"
	"log"
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
	log.Println("FIND ONE RECORD")
	var user models.User
	condition := *column + " = ?"
	r.db.Where(condition, value).Find(&user)
	byteData, _ := json.Marshal(user)
	log.Println("END FIND ONE RECORD")
	return &byteData, nil

}

func (r *UserRepo) FindAllRecords() (*[]byte, error) {
	log.Println("FIND RECORDS")
	var users []models.User
	r.db.Find(&users)
	byteData, _ := json.Marshal(users)
	log.Println("END FIND RECORDS")
	return &byteData, nil
}

func (r *UserRepo) Save(user interface{}) (*[]byte, error) {
	log.Println("SAVE RECORD")
	r.db.Create(user)
	byteData, _ := json.Marshal(user)
	log.Println("END SAVE RECORD")
	return &byteData, nil
}

func (r *UserRepo) Delete(id string) error {
	log.Println("DELETE RECORD")
	var users []models.User
	r.db.Where("id = ?", id).Delete(&users)
	log.Println("END DELETE RECORD")
	return nil
}

func (r *UserRepo) Update(id string, userUpdate interface{}) (*[]byte, error) {
	log.Println("UPDATE RECORD")
	r.db.Save(userUpdate)
	byteData, _ := json.Marshal(userUpdate)
	log.Println("END UPDATE RECORD")
	return &byteData, nil
}
