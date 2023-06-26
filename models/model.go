// url shortener model

package models

import (
	"fmt"
	"log"

	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Cizzor struct {
	ID       uint64 `json:"id" gorm:"primary_key"`
	Url      string `json:"url" gorm:"not null"`
	ShortUrl string `json:"short_url" gorm:"not null"`
	Count    int    `json:"count"`
	Random   bool   `json:"random"`
	OwnerId  uint64 `json:"owner_id" gorm:"foreignKey:UserID"`
}

type User struct {
	ID       uint64 `json:"id" gorm:"primary_key"`
	Username string `json:"username" gorm:"not null;unique"`
	Email    string `json:"email" gorm:"not null;unique"`
	Password string `json:"password" gorm:"not null"`
}

// type LoginUserModel struct {
// 	Email    string `json:"email" gorm:"not null"`
// 	Password string `json:"password" gorm:"not null"`
// }

var db *gorm.DB

func Setup() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// clear the DB for testing
	// err = db.Migrator().DropTable(&Cizzor{})
	// if err != nil {
	// 	fmt.Println(err)
	// }123456

	err = db.AutoMigrate(&Cizzor{}, &User{})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("DB migrated")
	}

}
