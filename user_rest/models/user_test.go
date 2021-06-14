package models

import "testing"

func TestUserModel(t *testing.T) {

	t.Log("Validate UserModel")

	user := &User{ID: 1, Login: "toto", Password: "toto", First_name: "toto", Email: "toto@toto.fr", Last_name: "toto"}

	if !(user.ID == 1 &&
		user.Login == "toto" &&
		user.First_name == "toto" &&
		user.Last_name == "toto" &&
		user.Email == "toto@toto.fr") {

		t.Errorf("ControllerTools UserModel assert error!")
	}

}

func TestUserUpdate(t *testing.T) {

	t.Log("Validate UserUpdate")

	user := &UserUpdate{ID: "1", Login: "toto", Password: "toto", First_name: "toto", Email: "toto@toto.fr", Last_name: "toto"}

	if !(user.ID == "1" &&
		user.Login == "toto" &&
		user.First_name == "toto" &&
		user.Last_name == "toto" &&
		user.Email == "toto@toto.fr") {

		t.Errorf("ControllerTools UserUpdate assert error!")
	}

}
