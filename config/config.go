package config

import (
	"fmt"
	"os"
	"stockz-app/model"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DatabaseInit() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connection := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbHost, dbUser, dbPassword, dbName, dbPort)
	db, e := gorm.Open(postgres.Open(connection), &gorm.Config{})
	if e != nil {
		panic("Error Connecting Databse")
	}
	fmt.Print("Connect to Database")
	db.AutoMigrate(&model.Post{}, model.User{}, model.Comment{})

	DB = db
}
