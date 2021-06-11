package routing

import (
	"user_rest/user_rest/controllers"

	"github.com/gorilla/mux"
)

func InitRooting(h controllers.BaseHandler) *mux.Router {

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", h.Status).Methods("GET")
	myRouter.HandleFunc("/users/", h.CreateUser).Methods("POST")
	myRouter.HandleFunc("/users/{id}", h.DeleteUser).Methods("DELETE")
	myRouter.HandleFunc("/users/{id}", h.UpdateUser).Methods("PUT")
	myRouter.HandleFunc("/users/", h.ListUsers).Methods("GET")

	return myRouter
}
