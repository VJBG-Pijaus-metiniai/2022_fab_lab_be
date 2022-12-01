package database

import (
	"fablab-project/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Db struct {
	DB *gorm.DB
}

var Database Db

func ConnectDB() {
	dbErr := godotenv.Load(".env")

	if dbErr != nil {
		log.Fatal("Failed to connect to db! \n", dbErr.Error())
		os.Exit(2)
	}

	dsn := os.Getenv("CONN_STRING")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to db! \n", err.Error())
		os.Exit(2)
	}

	log.Println("Connected to db!")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running migrations")
	db.AutoMigrate(&models.User{}, &models.Project{})

	Database = Db{DB: db}
}