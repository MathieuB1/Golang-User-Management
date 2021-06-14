package sqldb

import (
	"fmt"
	"log"
	"user_rest/user_rest/context"
	"user_rest/user_rest/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database Init
func initialMigration(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		panic(err.Error())
	}

}

// Database Conf
type DBInfo struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
}

func setDBConf(dbConf *DBInfo) string {
	return fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbConf.host, dbConf.port, dbConf.user, dbConf.password, dbConf.dbname)
}

// Database Connexion
func ConnectDB() *gorm.DB {

	// Init DB Conf
	pgConf := &DBInfo{
		host:     context.GlobalCtx.DB_HOST,
		port:     context.GlobalCtx.DB_PORT,
		dbname:   context.GlobalCtx.DB_NAME,
		user:     context.GlobalCtx.DB_USER,
		password: context.GlobalCtx.DB_PASSWORD,
	}

	db, err := gorm.Open(postgres.Open(setDBConf(pgConf)), &gorm.Config{})
	if err != nil {

		log.Println("Trying local Database instead...")

		pgConf.host = "localhost"
		pgConf.port = "5432"
		pgConf.dbname = "test"
		pgConf.user = "postgres"
		pgConf.password = "postgres"

		db, err = gorm.Open(postgres.Open(setDBConf(pgConf)), &gorm.Config{})
		if err != nil {
			panic(err.Error())
		}

	}

	initialMigration(db)

	return db
}
