package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(user User) (User, error) {
	result := db.Create(&user)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

func GetUserByID(id uint64) (User, error) {
	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

// func LoginUser(user LoginUserModel) (User, error) {
// 	var dbUser LoginUserModel
// 	var dbUserResult User
// 	result := db.Where("email = ?", user.Email).First(&dbUser)
// 	if result.Error != nil {
// 		return User{}, result.Error
// 	}
// 	return dbUserResult, nil
// }

func LoginUser(user User) (User, error) {
	var dbUser User
	result := db.Where("email = ?", user.Email).First(&dbUser)
	if result.Error != nil {
		return User{}, result.Error
	}
	// Compare the provided password with the stored password hash
	err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return User{}, errors.New("incorrect password")
	}

	return dbUser, nil
}

func UpdateUser(user User) (User, error) {
	result := db.Save(&user)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

func DeleteUser(id uint64) error {
	result := db.Delete(&User{}, id)
	return result.Error
}
