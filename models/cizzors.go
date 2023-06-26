package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func GetAllCizzors() ([]Cizzor, error) {
	var cizzors []Cizzor

	tx := db.Find(&cizzors)
	if tx.Error != nil {
		return []Cizzor{}, tx.Error
	}

	return cizzors, nil
}

func GetCizzorById(id uint64) (Cizzor, error) {
	var cizzor Cizzor

	tx := db.First(&cizzor, id)
	if tx.Error != nil {
		return Cizzor{}, tx.Error
	}

	return cizzor, nil
}

func GetCizzorByShortUrl(shortUrl string) (Cizzor, error) {
	var cizzor Cizzor

	tx := db.Where("short_url = ?", shortUrl).First(&cizzor)
	if tx.Error != nil {
		return Cizzor{}, tx.Error
	}

	return cizzor, nil
}

func CreateCizzor(cizzor Cizzor) (Cizzor, error) {
	tx := db.Create(&cizzor)
	if tx.Error != nil {
		return Cizzor{}, tx.Error
	}

	return cizzor, nil
}

// func UpdateCizzor(cizzor Cizzor) (Cizzor, error) {
// 	tx := db.Save(&cizzor)
// 	if tx.Error != nil {
// 		return Cizzor{}, tx.Error
// 	}

// 	return cizzor, nil
// }

func UpdateCizzor(cizzor Cizzor) (Cizzor, error) {
	// Find the existing cizzor by ID
	existingCizzor := Cizzor{}
	if err := db.First(&existingCizzor, cizzor.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Cizzor{}, fmt.Errorf("cizzor not found")
		}
		return Cizzor{}, err
	}

	// Exclude the ID field from update
	cizzor.ID = existingCizzor.ID

	// Perform the update
	if err := db.Model(&existingCizzor).Updates(cizzor).Error; err != nil {
		return Cizzor{}, err
	}

	return existingCizzor, nil
}

func DeleteCizzor(id uint64) error {
	tx := db.Delete(&Cizzor{}, id)
	//for hard delete use
	//tx := db.Unscoped().Delete(&cizzor)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
