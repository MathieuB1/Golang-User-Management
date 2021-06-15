package controllers

import (
	"log"
	"net/http"
	"user_rest/user_rest/common"
	"user_rest/user_rest/controllers_tools"
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
// Helper for Authentification
//**********************************
func (h *BaseHandler) IsAuthorised(login *string, password *string) (bool, error) {

	log.Println("Check Authorization...")

	firstname := "login"

	userByte, err := h.userRepo.FindOneRecord(&firstname, login)
	if err != nil {
		return false, err
	}

	user, err := controllers_tools.GetUserFromBytes(userByte)
	if err != nil {
		return false, err
	}

	pass := user.Password

	log.Println("End Check Authorization...")

	return common.CreateHash(password) == pass, nil
}

func (h *BaseHandler) IsUserAuth(w http.ResponseWriter, r *http.Request) *string {
	w.Header().Add("Content-Type", "application/json")

	log.Println("Starting Login...")
	login, password, ok := r.BasicAuth()

	isLogged := true

	if !ok {
		isLogged = false
	}

	if authStatus, _ := h.IsAuthorised(&login, &password); !authStatus {
		isLogged = false
	}

	if !isLogged {
		w.Header().Add("WWW-Authenticate", `Basic realm="Give login and password"`)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message": "Error bad or missing basic auth"}`))
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

	allUsers, err := controllers_tools.GetUsersFromBytes(allUsersBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("End ListUsers...")

	usersObj, err := common.SerializeSender(allUsers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(usersObj)
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
	}

	userExistInDB, err := controllers_tools.GetUserFromBytes(userByte)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Reject the request if the user already exists
	if userExistInDB.Login != user.Login {

		_, err := h.userRepo.Save(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Mask fields
		userDisplay := controllers_tools.FilterKeysUser(user)

		userObj, err := common.SerializeSender(userDisplay)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		log.Println("End CreateUser.")

		w.Header().Set("Content-Type", "application/json")
		w.Write(userObj)

	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "User "` + user.Login + `" already exist!"}`))
	}
}

func (h *BaseHandler) FindUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting FindUser...")

	userLogged := h.IsUserAuth(w, r)
	if userLogged != nil {

		var column = "id"
		vars := mux.Vars(r)
		id := vars[column]

		userByte, err := h.userRepo.FindOneRecord(&column, &id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := controllers_tools.GetUserFromBytes(userByte)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// Check if its the correct User
		if user.Login != *userLogged {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Mask fields
		displayUser := controllers_tools.FilterKeysUser(user)

		userObj, err := common.SerializeSender(displayUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		log.Println("End FindUser...")

		w.Header().Set("Content-Type", "application/json")
		w.Write(userObj)
	}

	log.Println("End FindUser...")
}

// Delete One User
func (h *BaseHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting DeleteUser...")

	userLogged := h.IsUserAuth(w, r)
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
		user, err := controllers_tools.GetUserFromBytes(userByte)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

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

		w.Write([]byte(`{"message": "User has been deleted!"}`))

	}

	log.Println("End DeleteUser.")
}

// Update One User
func (h *BaseHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting UpdateUser...")

	userLogged := h.IsUserAuth(w, r)
	if userLogged != nil {

		r.ParseForm()

		vars := mux.Vars(r)
		id := vars["id"]

		var login, password, email, first_name, last_name string

		// Seems that POSTMAN is not working
		// or Golang has some missing Content-Type application/json implementation
		if len(vars) > 1 {
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

		updatedUser := &models.UserUpdate{ID: id, Login: login, Password: password, First_name: first_name, Last_name: last_name, Email: email}

		var column = "id"
		existingUserBytes, err := h.userRepo.FindOneRecord(&column, &id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		existingUser, err := controllers_tools.GetUserFromBytes(existingUserBytes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

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
		userUpdated, err := controllers_tools.GetUserFromBytes(userUpdatedByte)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// Mask fields
		displayUser := controllers_tools.FilterKeysUser(userUpdated)

		userObj, err := common.SerializeSender(displayUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		log.Println("End UpdateUser.")

		w.Header().Set("Content-Type", "application/json")
		w.Write(userObj)

	}

	log.Println("End UpdateUser.")
}
