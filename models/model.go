// url shortener model

package models

import (
	"fmt"
	"log"
	"time"

	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ClickSource struct {
	ID        uint64    `json:"id" gorm:"primary_key"`
	CizzorID  uint64    `json:"cizzor_id" gorm:"foreignKey:CizzorID"`
	IpAddress string    `json:"ip_address" gorm:"not null"`
	Location  string    `json:"location" gorm:"not null"`
	Count     uint64    `json:"count"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp with time zone;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp with time zone;autoUpdateTime"`
}
type Cizzor struct {
	ID           uint64        `json:"id" gorm:"primary_key"`
	Url          string        `json:"url" gorm:"not null"`
	ShortUrl     string        `json:"short_url" gorm:"not null"`
	Count        int           `json:"count"`
	Random       bool          `json:"random"`
	OwnerId      uint64        `json:"owner_id" gorm:"foreignKey:UserID"`
	ClickSources []ClickSource `json:"-" gorm:"foreignKey:CizzorID"`
	CreatedAt    time.Time     `json:"created_at" gorm:"type:timestamp with time zone;autoCreateTime"`
	UpdatedAt    time.Time     `json:"updated_at" gorm:"type:timestamp with time zone;autoUpdateTime"`
}

type User struct {
	ID        uint64    `json:"id" gorm:"primary_key"`
	Username  string    `json:"username" gorm:"not null;unique"`
	Email     string    `json:"email" gorm:"not null;unique"`
	Password  string    `json:"password" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp with time zone;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp with time zone;autoUpdateTime"`
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

	// host := os.Getenv("DB_HOST")
	// host, err := os.Hostname()
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// fmt.Println(host)
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

	err = db.AutoMigrate(&Cizzor{}, &User{}, &ClickSource{})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("DB migrated")
	}

}
