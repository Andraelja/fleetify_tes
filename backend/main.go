package main

import (
	"backend/bootstrap"
	"backend/config"
	"backend/app/models"
)

func main() {
	config.ConnectDatabase()
	config.DB.AutoMigrate(
		&models.User{},
		&models.Supplier{},
		&models.Item{},
	)
	bootstrap.StartServer()
}
