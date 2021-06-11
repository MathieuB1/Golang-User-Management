package main

import (
	"api-test/web_rest/controllers"
	"api-test/web_rest/repositories"
	"api-test/web_rest/routing"
	"api-test/web_rest/sqldb"
	"log"
	"net/http"
)

func main() {

	// Init the Database Socket
	db := sqldb.ConnectDB()
	log.Println("Database OK!")

	userRepo := repositories.NewUserRepo(db)

	// Init Handlers
	h := controllers.NewBaseHandler(userRepo)
	log.Println("Handlers OK!")

	// Init Routes
	myRouter := routing.InitRooting(*h)
	log.Println("Routes OK!")

	// For admin
	log.Println("Starting Server...")
	log.Fatal(http.ListenAndServe(":8000", myRouter))
}
