package main

import (
	"encoding/json"
	"fmt"
	"time"

	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB
var DBErr error
var err error

type User struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	FirstName string
}

func main() {

	// Create DB connection
	db := initDB()
	defer db.Close()

	// createUsers()
	users := getUsers(getDB(DB))

	jsonres, _ := json.MarshalIndent(users, "", " ")
	fmt.Println(string(jsonres))

}

func getUsers(db *gorm.DB) []User {
	users := []User{}

	db.Find(&users)

	return users
}

func getUser(db *gorm.DB, id uint) (*User, error) {
	user := new(User)

	err := db.Where("id = ?", id).Find(user).Error

	return user, err
}

func saveUser(db *gorm.DB, user User) (User, error) {
	err := db.Save(&user).Error
	return user, err
}

func initDB() *gorm.DB {

	database := "postgres"
	username := "postgres"
	password := "postgres"
	host := "localhost"
	port := "5432"

	// POSTGRES
	db, err := gorm.Open("postgres", "host="+host+" port="+port+" user="+username+" dbname="+database+" sslmode=disable password="+password)
	if err != nil {
		DBErr = err
		log.Fatal(err)
	}

	db.LogMode(true)
	// db.DropTableIfExists(&User{})
	db.AutoMigrate(&User{})
	DB = db

	return DB
}

// GetDB helps you to get a connection
func getDB(db *gorm.DB) *gorm.DB {
	return DB
}

// GetDBErr helps you to get a connection
func getDBErr() error {
	return DBErr
}
