package services

import (
	"backend/app/models"
	"backend/config"
	"backend/exception"
	"context"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type purchasingService struct {
	webhookService WebhookService
}

func NewPurchasingService(webhookService WebhookService) PurchasingService {
	return &purchasingService{
		webhookService: webhookService,
	}
}

/* =========================
	CREATE
========================= */
func (s *purchasingService) Create(
	ctx context.Context,
	model models.PurchasingCreateOrUpdateModel,
) models.Purchasing {

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var (
		grandTotal int64
		details    []models.PurchasingDetail
	)

	for _, d := range model.Details {
		var item models.Item

		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&item, d.ItemID).Error; err != nil {

			tx.Rollback()
			if err == gorm.ErrRecordNotFound {
				exception.PanicLogging(fmt.Errorf("item %d tidak ditemukan", d.ItemID))
			}
			exception.PanicLogging(err)
		}

		subTotal := item.Price * int64(d.Qty)
		grandTotal += subTotal

		details = append(details, models.PurchasingDetail{
			ItemID:   item.ID,
			Qty:      d.Qty,
			Price:    item.Price,
			SubTotal: subTotal,
		})
	}

	purchasing := models.Purchasing{
		Date:       time.Now(), // 🔥 SET BACKEND
		SupplierID: model.SupplierID,
		UserID:     model.UserID,
		GrandTotal: grandTotal,
	}

	if err := tx.Create(&purchasing).Error; err != nil {
		tx.Rollback()
		exception.PanicLogging(err)
	}

	for i := range details {
		details[i].PurchasingID = purchasing.ID
	}

	if err := tx.Create(&details).Error; err != nil {
		tx.Rollback()
		exception.PanicLogging(err)
	}

	// Update stock
	for _, d := range details {
		if err := tx.Model(&models.Item{}).
			Where("id = ?", d.ItemID).
			UpdateColumn("stock", gorm.Expr("stock + ?", d.Qty)).
			Error; err != nil {

			tx.Rollback()
			exception.PanicLogging(err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		exception.PanicLogging(err)
	}

	config.DB.
		Preload("Supplier").
		Preload("User").
		Preload("Details.Item").
		First(&purchasing, purchasing.ID)

	if s.webhookService != nil {
		_ = s.webhookService.SendPurchasingNotification(purchasing)
	}

	return purchasing
}

/* =========================
	UPDATE
========================= */
func (s *purchasingService) Update(
	ctx context.Context,
	model models.PurchasingCreateOrUpdateModel,
	id string,
) models.Purchasing {

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	purchasingID, err := strconv.Atoi(id)
	if err != nil {
		exception.PanicLogging(err)
	}

	var purchasing models.Purchasing
	if err := tx.Preload("Details").First(&purchasing, purchasingID).Error; err != nil {
		tx.Rollback()
		exception.PanicLogging(err)
	}

	// Rollback old stock
	for _, d := range purchasing.Details {
		tx.Model(&models.Item{}).
			Where("id = ?", d.ItemID).
			UpdateColumn("stock", gorm.Expr("stock - ?", d.Qty))
	}

	tx.Where("purchasing_id = ?", purchasingID).
		Delete(&models.PurchasingDetail{})

	var (
		grandTotal int64
		newDetails []models.PurchasingDetail
	)

	for _, d := range model.Details {
		var item models.Item
		tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&item, d.ItemID)

		subTotal := item.Price * int64(d.Qty)
		grandTotal += subTotal

		newDetails = append(newDetails, models.PurchasingDetail{
			PurchasingID: purchasing.ID,
			ItemID:       item.ID,
			Qty:          d.Qty,
			Price:        item.Price,
			SubTotal:     subTotal,
		})
	}

	purchasing.SupplierID = model.SupplierID
	purchasing.UserID = model.UserID
	purchasing.GrandTotal = grandTotal
	purchasing.Date = time.Now()

	tx.Save(&purchasing)
	tx.Create(&newDetails)

	for _, d := range newDetails {
		tx.Model(&models.Item{}).
			Where("id = ?", d.ItemID).
			UpdateColumn("stock", gorm.Expr("stock + ?", d.Qty))
	}

	tx.Commit()

	config.DB.Preload("Supplier").Preload("User").Preload("Details.Item").
		First(&purchasing, purchasingID)

	return purchasing
}

/* =========================
	DELETE
========================= */
func (s *purchasingService) Delete(ctx context.Context, id string) {
	var purchasing models.Purchasing
	config.DB.Preload("Details").First(&purchasing, id)

	for _, d := range purchasing.Details {
		config.DB.Model(&models.Item{}).
			Where("id = ?", d.ItemID).
			UpdateColumn("stock", gorm.Expr("stock - ?", d.Qty))
	}

	config.DB.Delete(&purchasing)
}

func (s *purchasingService) FindById(ctx context.Context, id string) models.Purchasing {
	var purchasing models.Purchasing
	config.DB.Preload("Supplier").Preload("User").Preload("Details.Item").
		First(&purchasing, id)
	return purchasing
}

func (s *purchasingService) FindAll(ctx context.Context) []models.Purchasing {
	var data []models.Purchasing
	config.DB.Preload("Supplier").Preload("User").Preload("Details.Item").
		Find(&data)
	return data
}
