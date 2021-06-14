package controllers_tools

import (
	"encoding/json"
	"strconv"
	"user_rest/user_rest/models"
)

func GetUserFromBytes(userByte *[]byte) (*models.User, error) {
	var user models.User
	err := json.Unmarshal(*userByte, &user)
	if err != nil {
		return &models.User{}, err
	}
	return &user, nil
}

func GetUsersFromBytes(userByte *[]byte) (*[]models.User, error) {
	var users []models.User
	err := json.Unmarshal(*userByte, &users)
	if err != nil {
		return &[]models.User{}, err
	}
	return &users, nil
}

func FilterKeysUser(user *models.User) *models.UserUpdate {
	return &models.UserUpdate{ID: strconv.Itoa(user.ID), Login: user.Login,
		Password:   user.Password,
		First_name: user.First_name,
		Last_name:  user.Last_name,
		Email:      user.Email}
}
