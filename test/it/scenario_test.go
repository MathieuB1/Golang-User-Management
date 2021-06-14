package itegration_scenario_test

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

var host = "localhost:8000"
var url_host = "http://" + host

func TestHomePageIntegration(t *testing.T) {

	t.Log("Go to Home Page...")

	client := &http.Client{}
	req, err := http.NewRequest("GET", url_host, nil)

	if err != nil {
		t.Log(err)
	}
	req.Header.Add("Host", host)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err == nil {
		defer res.Body.Close()
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("url returned wrong status code: got %v want %v",
			res.StatusCode, http.StatusOK)
	}
	t.Log("Home Page catched!")
}

// Format User
type UserAssert struct {
	ID         string
	First_name string
	Last_name  string
	Email      string
}

func userAssertion(t *testing.T, reference *string, userCreated *UserAssert) {
	var userCompare UserAssert
	assert_with := reference
	json.Unmarshal([]byte(*assert_with), &userCompare)

	if userCreated.First_name != userCompare.First_name {
		t.Errorf("First_name issue!")
	}
	if userCreated.Last_name != userCompare.Last_name {
		t.Errorf("Last_name issue!")
	}
	if userCreated.Email != userCompare.Email {
		t.Errorf("Email issue!")
	}

}

// -- Create One User
// -- Read the User (Basic Auth)
// -- Update User (Basic Auth)
// -- Delete User (Basic Auth)
func TestBasicUserIntegeration(t *testing.T) {

	//*********************************************************
	// Creates One User and assing the ID received in IdCreated
	//**********************************************************
	t.Log("Creating One User...")

	var IdCreated = ""

	url := url_host + "/users/?login=titi&first_name=titi&last_name=titi&password=tuti&email=vfvdfv@vfvf.vf"

	client := &http.Client{}
	reqCreate, err := http.NewRequest("POST", url, nil)
	if err != nil {
		t.Log(err)
	}
	reqCreate.Header.Add("Host", host)
	reqCreate.Header.Set("Content-Type", "application/json")

	res, err := client.Do(reqCreate)
	if err == nil {
		defer res.Body.Close()
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("url returned wrong status code: got %v want %v",
			res.StatusCode, http.StatusOK)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Cannot read Body!")
	}

	var userCreated UserAssert
	json.Unmarshal([]byte(body), &userCreated)

	// Assert User Creation
	refCreate := `{"First_name":"titi","Last_name":"titi","Email":"vfvdfv@vfvf.vf"}`
	userAssertion(t, &refCreate, &userCreated)

	IdCreated = userCreated.ID
	t.Log("User " + IdCreated + " Created!")

	//*********************************************************
	// Read One User with Basic Auth
	//**********************************************************
	t.Log("Display One User...")

	url = url_host + "/users/" + IdCreated

	reqGet, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Log(err)
	}

	// encode in base64
	data := []byte("titi:tuti")
	pass := base64.StdEncoding.EncodeToString(data)

	reqGet.Header.Set("Authorization", "Basic "+pass)

	resGet, err := client.Do(reqGet)
	if err == nil {
		defer resGet.Body.Close()
	}

	if resGet.StatusCode != http.StatusOK {
		t.Errorf("url returned wrong status code: got %v want %v",
			resGet.StatusCode, http.StatusOK)
	}

	bodyGet, err := ioutil.ReadAll(resGet.Body)
	if err != nil {
		t.Errorf("Cannot read Body!")
	}

	var userRead UserAssert
	json.Unmarshal([]byte(bodyGet), &userRead)

	// Assert User Read
	refRead := `{"First_name":"titi","Last_name":"titi","Email":"vfvdfv@vfvf.vf"}`
	userAssertion(t, &refRead, &userRead)

	t.Log("User " + userRead.ID + " displayed!")

	//*********************************************************
	// Update User with Basic Auth
	//*********************************************************
	t.Log("Update One User...")

	url = url_host + "/users/" + IdCreated + "?login=toto&first_name=toto&last_name=toto&password=toto&email=toto@toto.fr"

	reqUpdate, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		t.Log(err)
	}

	// encode in base64
	reqUpdate.Header.Add("Authorization", "Basic "+pass)

	resUpdate, err := client.Do(reqUpdate)
	if err == nil {
		defer resUpdate.Body.Close()
	}

	bodyUpdate, err := ioutil.ReadAll(resUpdate.Body)
	if err != nil {
		t.Errorf("Cannot read Body!")
	}

	var userUpdate UserAssert
	json.Unmarshal([]byte(bodyUpdate), &userUpdate)

	// Assert User Read
	refUpdate := `{"First_name":"toto","Last_name":"toto","Email":"toto@toto.fr"}`
	userAssertion(t, &refUpdate, &userUpdate)

	t.Log("User " + userUpdate.ID + " updated!")

	//*********************************************************
	// Delete User with Basic Auth
	//**********************************************************
	t.Log("Delete One User...")

	url = url_host + "/users/" + IdCreated

	reqDelete, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		t.Log(err)
	}

	// encode in base64
	data = []byte("toto:toto")
	pass = base64.StdEncoding.EncodeToString(data)

	reqDelete.Header.Add("Authorization", "Basic "+pass)

	resDelete, err := client.Do(reqDelete)
	if err == nil {
		defer resDelete.Body.Close()
	}

	if resDelete.StatusCode != http.StatusOK {
		t.Errorf("url returned wrong status code: got %v want %v",
			resDelete.StatusCode, http.StatusOK)
	}

	// Assert User Missing
	t.Log("End User Deletion...")
}
