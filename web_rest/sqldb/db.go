package sqldb

import (
	"api-test/web_rest/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB opens a connection to the database
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "test"
)

// Database Init
func initialMigration(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
}

// Database Connexion
func ConnectDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	initialMigration(db)

	return db
}
