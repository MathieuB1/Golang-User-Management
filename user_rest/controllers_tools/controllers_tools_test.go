package controllers_tools

import (
	"encoding/json"
	"testing"
	"user_rest/user_rest/models"
)

func TestGetUserFromBytes(t *testing.T) {

	t.Log("Validate GetUserFromBytes")

	// Input
	var user = &models.User{ID: 1, Login: "toto"}
	byteData, _ := json.Marshal(user)

	user, err := GetUserFromBytes(&byteData)
	if err != nil {
		t.Errorf("ControllerTools GetUserFromBytes error!")
	}

	//Assert
	if user.ID != 1 && user.Login != "toto" {
		t.Errorf("ControllerTools GetUserFromBytes assert error!")
	}

	// Input
	user = &models.User{}
	byteData, _ = json.Marshal(user)

	user, err = GetUserFromBytes(&byteData)
	if err != nil {
		t.Errorf("ControllerTools GetUserFromBytes error!")
	}

	//Assert
	if user.ID != 0 && user.Login != "" {
		t.Errorf("ControllerTools GetUserFromBytes assert error!")
	}

}

func TestGetUsersFromBytes(t *testing.T) {

	t.Log("Validate GetUsersFromBytes")

	// Input
	var user = &[]models.User{{ID: 1, Login: "toto"},
		{ID: 2, Login: "tutu"}}

	byteData, _ := json.Marshal(user)

	users, err := GetUsersFromBytes(&byteData)
	if err != nil {
		t.Errorf("ControllerTools GetUsersFromBytes error!")
	}

	// Assert
	for index, element := range *users {
		if element.Login != (*user)[index].Login {
			t.Errorf("ControllerTools GetUserFromBytes assert error!")
		}
	}

}

func TestFilterKeysUser(t *testing.T) {

	t.Log("Validate FilterKeysUser")

	// Input
	var user = &models.User{ID: 1, Login: "toto", First_name: "toto", Password: "toto"}

	userRes := FilterKeysUser(user)

	// Assert
	if userRes.ID != "toto" && userRes.Login != "toto" {
		t.Errorf("ControllerTools FilterKeysUser assert error!")
	}

}
