package main

import (
	"backend/bootstrap"
	"backend/config"
	"backend/app/models"
	"log"
	"os"
)

func main() {
	config.LoadEnv()
	
	config.ConnectDatabase()
	
	// Auto migrate
	config.DB.AutoMigrate(
		&models.User{},
		&models.Supplier{},
		&models.Item{},
		&models.Purchasing{},
		&models.PurchasingDetail{},
	)
	
	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		log.Println("Warning: WEBHOOK_URL tidak ditemukan di .env, webhook notification tidak akan terkirim")
	} else {
		log.Println("Webhook URL configured:", webhookURL)
	}
	
	bootstrap.StartServer()
}