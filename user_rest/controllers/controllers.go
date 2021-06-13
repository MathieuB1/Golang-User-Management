package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"user_rest/user_rest/common"
	"user_rest/user_rest/models"

	"github.com/gorilla/mux"
)

// BaseHandler will hold everything that controller needs
type BaseHandler struct {
	userRepo models.UserRepository
}

// NewBaseHandler returns a new BaseHandler
func NewBaseHandler(userRepo models.UserRepository) *BaseHandler {
	return &BaseHandler{
		userRepo: userRepo,
	}
}

//**********************************
// Tools
//**********************************
func getUserFromBytes(w http.ResponseWriter, userByte *[]byte) *models.User {
	var user models.User
	err := json.Unmarshal(*userByte, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return &models.User{}
	}
	return &user
}

func getUsersFromBytes(w http.ResponseWriter, userByte *[]byte) *[]models.User {
	var user []models.User
	err := json.Unmarshal(*userByte, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return &[]models.User{}
	}
	return &user
}

func filterKeysUser(user *models.User) *models.UserUpdate {
	return &models.UserUpdate{ID: strconv.Itoa(user.ID), Login: user.Login,
		Password:   user.Password,
		First_name: user.First_name,
		Last_name:  user.Last_name,
		Email:      user.Email}
}

//**********************************
// Helper for Authentification
//**********************************
func isAuthorised(w http.ResponseWriter, h *BaseHandler, login *string, password *string) bool {
	firstname := "login"

	userByte, err := h.userRepo.FindOneRecord(&firstname, login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	user := getUserFromBytes(w, userByte)
	pass := user.Password

	return common.CreateHash(password) == pass
}

func IsUserAuth(h *BaseHandler, w http.ResponseWriter, r *http.Request) *string {
	w.Header().Add("Content-Type", "application/json")

	log.Println("Starting Login...")
	login, password, ok := r.BasicAuth()

	if !ok {
		w.Header().Add("WWW-Authenticate", `Basic realm="Give login and password"`)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message": "No basic auth present"}`))
		return nil
	}

	if !isAuthorised(w, h, &login, &password) {
		w.Header().Add("WWW-Authenticate", `Basic realm="Give login and password"`)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message": "Invalid username or password"}`))
		return nil
	}
	log.Println("End Login.")

	return &login
}

//**********************************
// Handlers
//**********************************

// Home Page
func (h *BaseHandler) Status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Welcome to Simple Test User Management!"}`))
}

// List all Users
func (h *BaseHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting ListUsers...")

	allUsersBytes, err := h.userRepo.FindAllRecords()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	allUsers := getUsersFromBytes(w, allUsersBytes)

	log.Println("End ListUsers...")
	common.SerializeAndSendResponse(&w, allUsers)
}

// Create One User
func (h *BaseHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting CreateUser...")
	vars := mux.Vars(r)
	r.ParseForm()

	var login, password, email, first_name, last_name string

	// Seems that POSTMAN is not working
	// or Golang has some missing Content-Type application/json implementation
	if len(vars) != 0 {
		login = vars["login"]
		password = vars["password"]
		first_name = vars["first_name"]
		last_name = vars["last_name"]
		email = vars["email"]

	} else {
		login = r.FormValue("login")
		password = r.FormValue("password")
		first_name = r.FormValue("first_name")
		last_name = r.FormValue("last_name")
		email = r.FormValue("email")
	}

	if password != "" {
		password = common.CreateHash(&password)
	}

	var user = &models.User{Login: login, Password: password,
		First_name: first_name,
		Last_name:  last_name,
		Email:      email}

	// Check if user exists
	var column = "login"
	userByte, err := h.userRepo.FindOneRecord(&column, &user.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userExistInDB := getUserFromBytes(w, userByte)

	// Reject the request if the user already exists
	if userExistInDB.Login != user.Login {

		_, err := h.userRepo.Save(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Mask fields
		userDisplay := filterKeysUser(user)

		log.Println("End CreateUser.")
		common.SerializeAndSendResponse(&w, userDisplay)

	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "User already exist!"}`))
	}
}

func (h *BaseHandler) FindUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting FindUser...")

	userLogged := IsUserAuth(h, w, r)
	if userLogged != nil {

		var column = "id"
		vars := mux.Vars(r)
		id := vars[column]

		userByte, err := h.userRepo.FindOneRecord(&column, &id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user := getUserFromBytes(w, userByte)

		// Check if its the correct User
		if user.Login != *userLogged {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Mask fields
		displayUser := filterKeysUser(user)

		common.SerializeAndSendResponse(&w, displayUser)
	}

	log.Println("End FindUser...")
}

// Delete One User
func (h *BaseHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting DeleteUser...")

	userLogged := IsUserAuth(h, w, r)
	if userLogged != nil {

		vars := mux.Vars(r)
		id := vars["id"]

		// Check if its the correct User
		var column = "id"
		userByte, err := h.userRepo.FindOneRecord(&column, &id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user := getUserFromBytes(w, userByte)

		// Check if its the correct User
		if user.Login != *userLogged {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		deletionStatus := h.userRepo.Delete(id)
		if deletionStatus != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("End DeleteUser.")

		w.Write([]byte(`{"message": "User has been deleted!"}`))

	}

	log.Println("End DeleteUser.")
}

// Update One User
func (h *BaseHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting UpdateUser...")

	userLogged := IsUserAuth(h, w, r)
	if userLogged != nil {

		r.ParseForm()

		vars := mux.Vars(r)
		id := vars["id"]

		var login, password, email, first_name, last_name string

		login = r.FormValue("login")
		password = r.FormValue("password")
		first_name = r.FormValue("first_name")
		last_name = r.FormValue("last_name")
		email = r.FormValue("email")

		if password != "" {
			password = common.CreateHash(&password)
		}

		updatedUser := &models.UserUpdate{ID: id, Login: login, Password: password, First_name: first_name, Last_name: last_name, Email: email}

		var column = "id"
		existingUserBytes, err := h.userRepo.FindOneRecord(&column, &id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		existingUser := getUserFromBytes(w, existingUserBytes)

		// Check if its the correct User
		if existingUser.Login != *userLogged {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Apply Patch
		if updatedUser.Email != "" {
			existingUser.Email = updatedUser.Email
		}
		if updatedUser.Login != "" {
			existingUser.Login = updatedUser.Login
		}
		if updatedUser.Password != "" {
			existingUser.Password = updatedUser.Password
		}
		if updatedUser.First_name != "" {
			existingUser.First_name = updatedUser.First_name
		}
		if updatedUser.Last_name != "" {
			existingUser.Last_name = updatedUser.Last_name
		}

		userUpdatedByte, err := h.userRepo.Update(id, existingUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userUpdated := getUserFromBytes(w, userUpdatedByte)

		// Mask fields
		displayUser := filterKeysUser(userUpdated)

		log.Println("End UpdateUser.")
		common.SerializeAndSendResponse(&w, displayUser)
	}

	log.Println("End UpdateUser.")
}
