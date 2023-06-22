// url shortener model

package models

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Cizzor struct {
	ID       uint64 `json:"id" gorm:"primary_key"`
	Url      string `json:"url" gorm:"not null"`
	ShortUrl string `json:"short_url" gorm:"not null"`
	Count    int    `json:"count"`
	Random   bool   `json:"random"`
}

var db *gorm.DB

func Setup() {
	dsn := "host=localhost user=usercizzors password=password dbname=cizzors port=5432 sslmode=disable"

	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// clear the DB for testing
	// err = db.Migrator().DropTable(&Cizzor{})
	// if err != nil {
	// 	fmt.Println(err)
	// }

	err = db.AutoMigrate(&Cizzor{})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("DB migrated")
	}

}
