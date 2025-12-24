package services

import (
	"backend/app/models"
	"backend/config"
	"backend/utils"
)

func RegisterUser(username, password string, role string) error {
	hashed, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	user := models.User{Username: username, Password: hashed, Role: role}
	return config.DB.Create(&user).Error
}

func LoginUser(username, password string) (string, error) {
	var user models.User
	err := config.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return "", err
	}
	err = utils.CheckPassword(user.Password, password)
	if err != nil {
		return "", err
	}
	return utils.GenerateToken(username)
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
