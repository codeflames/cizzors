package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// func GetAllCizzors() ([]Cizzor, error) {
// 	var cizzors []Cizzor

// 	tx := db.Find(&cizzors)
// 	if tx.Error != nil {
// 		return []Cizzor{}, tx.Error
// 	}

// 	return cizzors, nil
// }

// func GetCizzorById(id uint64) (Cizzor, error) {
// 	var cizzor Cizzor

// 	tx := db.First(&cizzor, id)
// 	if tx.Error != nil {
// 		return Cizzor{}, tx.Error
// 	}

// 	return cizzor, nil
// }

func GetCizzorByShortUrl(shortUrl string) (Cizzor, error) {
	var cizzor Cizzor

	tx := db.Where("short_url = ?", shortUrl).First(&cizzor)
	if tx.Error != nil {
		return Cizzor{}, tx.Error
	}

	return cizzor, nil
}

// func CreateCizzor(cizzor Cizzor) (Cizzor, error) {
// 	tx := db.Create(&cizzor)
// 	if tx.Error != nil {
// 		return Cizzor{}, tx.Error
// 	}

// 	return cizzor, nil
// }

// func UpdateCizzor(cizzor Cizzor) (Cizzor, error) {
// 	tx := db.Save(&cizzor)
// 	if tx.Error != nil {
// 		return Cizzor{}, tx.Error
// 	}

// 	return cizzor, nil
// }

func UpdateCizzorCount(cizzor Cizzor) (Cizzor, error) {
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

func GetAllCizzors(userID uint64) ([]Cizzor, error) {
	var cizzors []Cizzor

	tx := db.Where("owner_id = ?", userID).Find(&cizzors)
	if tx.Error != nil {
		return []Cizzor{}, tx.Error
	}

	return cizzors, nil
}

func GetCizzorById(id uint64) (Cizzor, error) {
	var cizzor Cizzor

	tx := db.Where("id = ?", id).First(&cizzor)
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

func UpdateCizzor(userID uint64, cizzor Cizzor) (Cizzor, error) {
	// Find the existing cizzor by ID and owner ID
	existingCizzor := Cizzor{}
	if err := db.Where("owner_id = ? AND id = ?", userID, cizzor.ID).First(&existingCizzor).Error; err != nil {
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
	tx := db.Where("id = ?", id).Delete(&Cizzor{})
	// For hard delete use
	//tx := db.Unscoped().Where("owner_id = ? AND id = ?", userID, id).Delete(&Cizzor{})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (c *Cizzor) UpdateClickSourceCount(ipAddress, location string) {
	// Check if the click source already exists in the database
	var existingClickSource ClickSource
	err := db.Where("cizzor_id = ? AND ip_address = ? AND location = ?", c.ID, ipAddress, location).First(&existingClickSource).Error
	if err == nil {
		// Click source already exists, increment the count
		existingClickSource.Count++
		db.Save(&existingClickSource)
		return
	}

	// If click source doesn't exist, create a new one with count = 1
	newClickSource := ClickSource{
		CizzorID:  c.ID,
		IpAddress: ipAddress,
		Location:  location,
		Count:     1,
	}
	c.ClickSources = append(c.ClickSources, newClickSource)

	// Save the cizzor
	db.Save(&c)
}

func GetClickSources(userID, cizzorID uint64) (ClickSource, error) {
	var clickSources ClickSource

	tx := db.Joins("JOIN cizzors ON click_sources.cizzor_id = cizzors.id").
		Where("cizzors.owner_id = ? AND cizzors.id = ?", userID, cizzorID).
		Find(&clickSources)
	if tx.Error != nil {
		return ClickSource{}, tx.Error
	}

	return clickSources, nil
}

// func GetClickSources(userID, cizzorID uint64) ([]ClickSource, error) {
// 	var clickSources []ClickSource

// 	tx := db.Where("cizzor_id = ?", cizzorID).Find(&clickSources)
// 	if tx.Error != nil {
// 		return []ClickSource{}, tx.Error
// 	}

// 	return clickSources, nil
// }

// func DeleteCizzor(id uint64) error {
// 	tx := db.Delete(&Cizzor{}, id)
// 	//for hard delete use
// 	//tx := db.Unscoped().Delete(&cizzor)
// 	if tx.Error != nil {
// 		return tx.Error
// 	}

// 	return nil
// }
