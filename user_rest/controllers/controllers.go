package controllers

import (
	"log"
	"net/http"
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
// Helper for Authentification
//**********************************
func isAuthorised(h *BaseHandler, login *string, password *string) bool {
	firstname := "login"
	user, _ := h.userRepo.FindOneRecord(&firstname, login)
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

	if !isAuthorised(h, &login, &password) {
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

	allUsers, err := h.userRepo.FindAllRecords()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	password = common.CreateHash(&password)

	var user = &models.User{Login: login, Password: password, First_name: first_name, Last_name: last_name, Email: email}

	_, err := h.userRepo.Save(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("End CreateUser.")
	common.SerializeAndSendResponse(&w, user)
}

func (h *BaseHandler) FindUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting FindUser...")

	userLogged := IsUserAuth(h, w, r)
	if userLogged != nil {

		var column = "id"
		vars := mux.Vars(r)
		id := vars[column]

		user, err := h.userRepo.FindOneRecord(&column, &id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if its the correct User
		if user.Login != *userLogged {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		common.SerializeAndSendResponse(&w, user)
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
		user, err := h.userRepo.FindOneRecord(&column, &id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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

		password = common.CreateHash(&password)

		updatedUser := &models.UserUpdate{ID: id, Login: login, Password: password, First_name: first_name, Last_name: last_name, Email: email}

		var column = "id"
		existingUser, err := h.userRepo.FindOneRecord(&column, &id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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

		userUpdated, err := h.userRepo.Update(id, existingUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("End UpdateUser.")
		common.SerializeAndSendResponse(&w, userUpdated)
	}

	log.Println("End UpdateUser.")
}
