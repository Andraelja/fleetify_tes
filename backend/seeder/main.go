package main

import (
	"backend/app/models"
	"backend/config"
	"backend/utils"
	"log"
	"time"

	"gorm.io/gorm"
)

func main() {
	runSeeder()
}

func runSeeder() {
	config.LoadEnv()
	config.ConnectDatabase()

	config.DB.AutoMigrate(
		&models.User{},
		&models.Supplier{},
		&models.Item{},
		&models.Purchasing{},
		&models.PurchasingDetail{},
	)

	log.Println("Starting database seeding...")

	seedUsers()
	seedSuppliers()
	seedItems()
	seedPurchasingTransactions()

	log.Println("Database seeding completed successfully!")
}

func seedUsers() {
	users := []models.User{
		{
			Username: "admin",
			Password: hashPassword("admin123"),
			Role:     "admin",
		},
	}

	for _, user := range users {
		var existing models.User
		if err := config.DB.Where("username = ?", user.Username).First(&existing).Error; err != nil {
			config.DB.Create(&user)
		}
	}
}

func seedSuppliers() {
	suppliers := []models.Supplier{
		{Name: "PT Teknologi Maju", Email: "contact@tm.com", Address: "Jakarta"},
		{Name: "CV Sumber Barokah", Email: "info@sb.com", Address: "Bandung"},
	}

	for _, s := range suppliers {
		var existing models.Supplier
		if err := config.DB.Where("email = ?", s.Email).First(&existing).Error; err != nil {
			config.DB.Create(&s)
		}
	}
}

func seedItems() {
	items := []models.Item{
		{Name: "Laptop", Stock: 10, Price: 15000000},
		{Name: "Mouse", Stock: 50, Price: 300000},
	}

	for _, i := range items {
		var existing models.Item
		if err := config.DB.Where("name = ?", i.Name).First(&existing).Error; err != nil {
			config.DB.Create(&i)
		}
	}
}

func seedPurchasingTransactions() {
	var user models.User
	if err := config.DB.First(&user).Error; err != nil {
		return
	}

	var supplier models.Supplier
	if err := config.DB.First(&supplier).Error; err != nil {
		return
	}

	var item models.Item
	if err := config.DB.First(&item).Error; err != nil {
		return
	}

	tx := config.DB.Begin()

	p := models.Purchasing{
		Date:       time.Now(),
		SupplierID: supplier.ID,
		UserID:     user.ID,
		GrandTotal: item.Price * 2,
	}

	if err := tx.Create(&p).Error; err != nil {
		tx.Rollback()
		return
	}

	d := models.PurchasingDetail{
		PurchasingID: p.ID,
		ItemID:       item.ID,
		Qty:          2,
		SubTotal:     item.Price * 2,
	}

	tx.Create(&d)

	tx.Model(&models.Item{}).
		Where("id = ?", item.ID).
		UpdateColumn("stock", gorm.Expr("stock + ?", 2))

	tx.Commit()
}

func hashPassword(password string) string {
	hashed, _ := utils.HashPassword(password)
	return hashed
}
