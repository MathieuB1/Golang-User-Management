package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"user_rest/user_rest/common"
	"user_rest/user_rest/controllers"
	"user_rest/user_rest/models"
	"user_rest/user_rest/repositories"
	"user_rest/user_rest/sqldb"

	"github.com/gorilla/mux"
)

// Format User
type UserAssert struct {
	ID         string
	First_name string
	Login      string
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

func createNewUserRepo() (*controllers.BaseHandler, *repositories.UserRepo) {
	// Init Database Socket
	db := sqldb.ConnectDB()
	// Create User Repo
	userRepo := repositories.NewUserRepo(db)
	// Init Handlers
	return controllers.NewBaseHandler(userRepo), userRepo
}

// Test Home Page
func TestCheckStatus(t *testing.T) {

	h, _ := createNewUserRepo()

	// Prepare Request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add Headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Host", "localhost:8000")

	// Test the Handler
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Status)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// Test User Creation Handler
func TestShouldCreateUser(t *testing.T) {

	var randomUser = "toto"

	t.Log("Test User Creation")

	req, err := http.NewRequest("POST", "/users/?login="+randomUser+"&first_name="+randomUser+"&last_name="+randomUser+"&password="+randomUser+"&email=toto@toto.fr", nil)
	if err != nil {
		t.Fatal(err)
	}

	h, userRepo := createNewUserRepo()

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.CreateUser)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	//Assert User
	assertWith := `{"First_name":"toto","Last_name":"toto","Email":"toto@toto.fr"}`

	var userCreated UserAssert
	json.Unmarshal(rr.Body.Bytes(), &userCreated)

	userAssertion(t, &assertWith, &userCreated)

	t.Log("User " + userCreated.ID + " has been created!")

	// Delete the user
	userRepo.Delete(userCreated.ID)
}

// Test Read User with Basic Auth
func TestShouldReadUserWithAuth(t *testing.T) {
	randomUser := "tutu"

	// Create a New User
	user := models.User{Password: common.CreateHash(&randomUser), Login: randomUser, First_name: "john", Last_name: "wick", Email: "john@wick.com"}

	// Init User Repo
	h, userRepo := createNewUserRepo()

	// Creates one User
	userByte, _ := userRepo.Save(&user)

	var userRead models.User
	json.Unmarshal(*userByte, &userRead)

	// Prepare Request
	req, err := http.NewRequest("GET", "/users/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(userRead.ID)})

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Host", "localhost:8000")
	req.SetBasicAuth(randomUser, randomUser)

	// Test the Handler
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.FindUser)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	//Assert User
	assertWith := `{"First_name":"john","Last_name":"wick","Email":"john@wick.com"}`

	var userReadRes UserAssert
	json.Unmarshal(rr.Body.Bytes(), &userReadRes)

	userAssertion(t, &assertWith, &userReadRes)

	// Delete the user
	userRepo.Delete(userReadRes.ID)
}

func TestShouldDeleteWithAuth(t *testing.T) {

	randomUser := "tutu"

	// Create a New User
	user := models.User{Password: common.CreateHash(&randomUser), Login: randomUser, First_name: "john", Last_name: "wick", Email: "john@wick.com"}

	// Init User Repo
	h, userRepo := createNewUserRepo()

	userByte, _ := userRepo.Save(&user)

	var userRead models.User
	json.Unmarshal(*userByte, &userRead)

	// Prepare Request
	req, err := http.NewRequest("DELETE", "/users/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(userRead.ID)})

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Host", "localhost:8000")
	req.SetBasicAuth(randomUser, randomUser)

	// Test the Handler
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.DeleteUser)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func TestShouldUpdateWithAuth(t *testing.T) {

	randomUser := "tutu"

	// Create a New User
	user := models.User{Password: common.CreateHash(&randomUser), Login: randomUser, First_name: "john", Last_name: "wick", Email: "john@wick.com"}

	// Init User Repo
	h, userRepo := createNewUserRepo()

	// Creates One User
	userByte, _ := userRepo.Save(&user)

	var userRead models.User
	json.Unmarshal(*userByte, &userRead)

	// Prepare Request
	req, err := http.NewRequest("PUT", "/users/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(userRead.ID), "first_name": "toto", "last_name": "toto"})

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Host", "localhost:8000")
	req.SetBasicAuth(randomUser, randomUser)

	// Test the Handler
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.UpdateUser)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var userUpdatedRes UserAssert
	json.Unmarshal(rr.Body.Bytes(), &userUpdatedRes)

	//Assert User
	assertWith := `{"First_name":"toto","Last_name":"toto","Email":"john@wick.com"}`

	userAssertion(t, &assertWith, &userUpdatedRes)

	// Delete the user
	userRepo.Delete(userUpdatedRes.ID)
}

func TestUserAuthEmpty(t *testing.T) {

	randomUser := "tutu"

	// Create a New User
	user := models.User{Password: common.CreateHash(&randomUser), Login: randomUser, First_name: "john", Last_name: "wick", Email: "john@wick.com"}

	// Init User Repo
	h, userRepo := createNewUserRepo()

	userByte, _ := userRepo.Save(&user)

	var userRead models.User
	json.Unmarshal(*userByte, &userRead)

	// Prepare Request
	req, err := http.NewRequest("GET", "/users/", nil)
	if err != nil {
		t.Fatal(err)
	}

	var userTest, userPass string

	// Empty User/Pass
	req.SetBasicAuth(userTest, userPass)

	// Add Headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Host", "localhost:8000")

	// Test the Handler
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.FindUser)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	userRepo.Delete(strconv.Itoa(userRead.ID))
}

func TestUserAuthNotExist(t *testing.T) {

	randomUser := "tutu"

	// Create a New User
	user := models.User{Password: common.CreateHash(&randomUser), Login: randomUser, First_name: "john", Last_name: "wick", Email: "john@wick.com"}

	// Init User Repo
	h, userRepo := createNewUserRepo()

	userByte, _ := userRepo.Save(&user)

	var userRead models.User
	json.Unmarshal(*userByte, &userRead)

	// Prepare Request
	req, err := http.NewRequest("GET", "/users/", nil)
	if err != nil {
		t.Fatal(err)
	}

	userTest := "None"
	userPass := "None"

	// Empty User/Pass
	req.SetBasicAuth(userTest, userPass)

	// Add Headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Host", "localhost:8000")

	// Test the Handler
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.FindUser)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	userRepo.Delete(strconv.Itoa(userRead.ID))
}

func TestUserAuthUser(t *testing.T) {

	randomUser := "tutu"

	// Create a New User
	user := models.User{Password: common.CreateHash(&randomUser), Login: randomUser, First_name: "john", Last_name: "wick", Email: "john@wick.com"}

	// Init User Repo
	h, userRepo := createNewUserRepo()

	userByte, _ := userRepo.Save(&user)

	var userRead models.User
	json.Unmarshal(*userByte, &userRead)

	// Prepare Request
	req, err := http.NewRequest("GET", "/users/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set User ID
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(userRead.ID)})

	// Set User
	req.SetBasicAuth(randomUser, randomUser)

	// Add Headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Host", "localhost:8000")

	// Test the Handler
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.FindUser)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	//Assert User
	assertWith := `{"First_name":"john","Last_name":"wick","Email":"john@wick.com"}`

	var userCreated UserAssert
	json.Unmarshal(rr.Body.Bytes(), &userCreated)

	userAssertion(t, &assertWith, &userCreated)

	userRepo.Delete(strconv.Itoa(userRead.ID))
}

func TestIsAuthorised(t *testing.T) {

	randomUser := "tutu"

	// Create a New User
	user := models.User{Password: common.CreateHash(&randomUser), Login: randomUser, First_name: "john", Last_name: "wick", Email: "john@wick.com"}

	// Init User Repo
	h, userRepo := createNewUserRepo()

	userByte, _ := userRepo.Save(&user)

	var userRead models.User
	json.Unmarshal(*userByte, &userRead)

	// Check User exist
	state, _ := h.IsAuthorised(&randomUser, &randomUser)

	if state != true {
		t.Errorf("handler returned wrong status code: got %v want %v",
			state, true)
	}

	// Check User not exist
	randomUser = "titi"
	state, _ = h.IsAuthorised(&randomUser, &randomUser)

	if state != false {
		t.Errorf("handler returned wrong status code: got %v want %v",
			state, false)
	}

	userRepo.Delete(strconv.Itoa(userRead.ID))

}

func TestListUsers(t *testing.T) {

	randomUser := "tutu"

	// Create a New User
	user := models.User{Password: common.CreateHash(&randomUser), Login: randomUser, First_name: "john", Last_name: "wick", Email: "john@wick.com"}

	// Init User Repo
	h, userRepo := createNewUserRepo()

	userByte, _ := userRepo.Save(&user)

	var userRead models.User
	json.Unmarshal(*userByte, &userRead)

	// Prepare Request
	req, err := http.NewRequest("GET", "/users/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add Headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Host", "localhost:8000")

	// Test the Handler
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.ListUsers)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	//Assert User
	assertWith := `{"First_name":"john","Last_name":"wick","Email":"john@wick.com"}`

	var userCreated []UserAssert
	json.Unmarshal(rr.Body.Bytes(), &userCreated)

	userAssertion(t, &assertWith, &userCreated[0])

	userRepo.Delete(strconv.Itoa(userRead.ID))

}
