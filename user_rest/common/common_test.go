package common

import (
	"encoding/json"
	"testing"
	"user_rest/user_rest/models"
)

func TestCommonHash(t *testing.T) {

	t.Log("Validate Hash")
	var stringTest = ""
	res := CreateHash(&stringTest)

	if !(len(res) > 0) {
		t.Errorf("MD5 error!")
	}

	stringTest = "toto"
	res = CreateHash(&stringTest)

	if !(len(res) > 0) {
		t.Errorf("MD5 error!")
	}

	stringTest = "salt" + "toto"
	res = CreateHash(&stringTest)

	if !(len(res) > 0) {
		t.Errorf("MD5 error!")
	}

}

func TestSerializeAndSendResponse(t *testing.T) {

	t.Log("Validate SerializeSender")
	// Input
	user := &models.User{ID: 1, Login: "toto"}

	res, err := SerializeSender(&user)
	if err != nil {
		t.Errorf("SerializeSender encoder error!")
	}

	var decodedUser = &models.User{}
	json.Unmarshal(res, decodedUser)

	//Assert Result
	if decodedUser.ID != 1 && decodedUser.Login != "toto" {
		t.Errorf("SerializeSender decoder error!")
	}

	// Input
	user = &models.User{}

	res, err = SerializeSender(&user)
	if err != nil {
		t.Errorf("SerializeSender encoding empty Obj error!")
	}

	decodedUser = &models.User{}
	json.Unmarshal(res, decodedUser)

	//Assert Result
	if decodedUser.ID != 0 {
		t.Errorf("SerializeSender decoder error!")
	}

}
