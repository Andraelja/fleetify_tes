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

func (s *purchasingService) Create(
	ctx context.Context,
	model models.PurchasingCreateOrUpdateModel,
) models.PurchasingCreateOrUpdateModel {

	date, err := time.Parse("2006-01-02", model.Date)
	exception.PanicLogging(err)

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var (
		grandTotal        int64
		purchasingDetails []models.PurchasingDetail
	)

	for _, detail := range model.Details {
		var item models.Item

		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&item, detail.ItemID).Error; err != nil {

			tx.Rollback()
			if err == gorm.ErrRecordNotFound {
				exception.PanicLogging(
					fmt.Errorf("Item dengan ID %d tidak ditemukan", detail.ItemID),
				)
			}
			exception.PanicLogging(err)
		}

		subTotal := item.Price * int64(detail.Qty)
		grandTotal += subTotal

		purchasingDetails = append(purchasingDetails, models.PurchasingDetail{
			ItemID:   item.ID,
			Qty:      detail.Qty,
			Price:    item.Price,
			SubTotal: subTotal,
		})
	}

	purchasing := models.Purchasing{
		Date:       date,
		SupplierID: model.SupplierID,
		UserID:     model.UserID,
		GrandTotal: grandTotal,
	}

	if err := tx.Create(&purchasing).Error; err != nil {
		tx.Rollback()
		exception.PanicLogging(
			fmt.Errorf("Gagal membuat purchasing: %w", err),
		)
	}

	for i := range purchasingDetails {
		purchasingDetails[i].PurchasingID = purchasing.ID
	}

	if err := tx.Create(&purchasingDetails).Error; err != nil {
		tx.Rollback()
		exception.PanicLogging(
			fmt.Errorf("Gagal membuat purchasing details: %w", err),
		)
	}

	for _, detail := range purchasingDetails {
		if err := tx.Model(&models.Item{}).
			Where("id = ?", detail.ItemID).
			UpdateColumn("stock", gorm.Expr("stock + ?", detail.Qty)).
			Error; err != nil {

			tx.Rollback()
			exception.PanicLogging(
				fmt.Errorf("Gagal update stock item %d: %w", detail.ItemID, err),
			)
		}
	}

	if err := tx.Commit().Error; err != nil {
		exception.PanicLogging(
			fmt.Errorf("Gagal commit transaction: %w", err),
		)
	}

	config.DB.
		Preload("Supplier").
		Preload("User").
		Preload("Details.Item").
		First(&purchasing, purchasing.ID)

	if s.webhookService != nil {
		_ = s.webhookService.SendPurchasingNotification(purchasing)
	}

	model.GrandTotal = grandTotal
	return model
}

func (s *purchasingService) Update(ctx context.Context, model models.PurchasingCreateOrUpdateModel, id string) models.PurchasingCreateOrUpdateModel {
	date, err := time.Parse("2006-01-02", model.Date)
	exception.PanicLogging(err)

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var existingPurchasing models.Purchasing
	purchasingID, err := strconv.Atoi(id)
	if err != nil {
		tx.Rollback()
		exception.PanicLogging(fmt.Errorf("Invalid purchasing ID: %w", err))
	}
	if err := tx.Preload("Details").First(&existingPurchasing, purchasingID).Error; err != nil {
		tx.Rollback()
		exception.PanicLogging(fmt.Errorf("Purchasing not found: %w", err))
	}

	// Subtract old stock
	for _, detail := range existingPurchasing.Details {
		if err := tx.Model(&models.Item{}).
			Where("id = ?", detail.ItemID).
			UpdateColumn("stock", gorm.Expr("stock - ?", detail.Qty)).
			Error; err != nil {
			tx.Rollback()
			exception.PanicLogging(fmt.Errorf("Failed to update stock for item %d: %w", detail.ItemID, err))
		}
	}

	// Delete old details
	if err := tx.Where("purchasing_id = ?", id).Delete(&models.PurchasingDetail{}).Error; err != nil {
		tx.Rollback()
		exception.PanicLogging(fmt.Errorf("Failed to delete old details: %w", err))
	}

	// Now, similar to Create: calculate new grandTotal, create new details, add stock
	var (
		grandTotal        int64
		purchasingDetails []models.PurchasingDetail
	)

	for _, detail := range model.Details {
		var item models.Item
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&item, detail.ItemID).Error; err != nil {
			tx.Rollback()
			if err == gorm.ErrRecordNotFound {
				exception.PanicLogging(fmt.Errorf("Item with ID %d not found", detail.ItemID))
			}
			exception.PanicLogging(err)
		}

		subTotal := item.Price * int64(detail.Qty)
		grandTotal += subTotal

		purchasingDetails = append(purchasingDetails, models.PurchasingDetail{
			PurchasingID: existingPurchasing.ID,
			ItemID:       item.ID,
			Qty:          detail.Qty,
			Price:        item.Price,
			SubTotal:     subTotal,
		})
	}

	// Update purchasing
	existingPurchasing.Date = date
	existingPurchasing.SupplierID = model.SupplierID
	existingPurchasing.UserID = model.UserID
	existingPurchasing.GrandTotal = grandTotal

	if err := tx.Save(&existingPurchasing).Error; err != nil {
		tx.Rollback()
		exception.PanicLogging(fmt.Errorf("Failed to update purchasing: %w", err))
	}

	// Create new details
	if err := tx.Create(&purchasingDetails).Error; err != nil {
		tx.Rollback()
		exception.PanicLogging(fmt.Errorf("Failed to create purchasing details: %w", err))
	}

	// Add new stock
	for _, detail := range purchasingDetails {
		if err := tx.Model(&models.Item{}).
			Where("id = ?", detail.ItemID).
			UpdateColumn("stock", gorm.Expr("stock + ?", detail.Qty)).
			Error; err != nil {
			tx.Rollback()
			exception.PanicLogging(fmt.Errorf("Failed to update stock for item %d: %w", detail.ItemID, err))
		}
	}

	if err := tx.Commit().Error; err != nil {
		exception.PanicLogging(fmt.Errorf("Failed to commit transaction: %w", err))
	}

	// Reload
	config.DB.Preload("Supplier").Preload("User").Preload("Details.Item").First(&existingPurchasing, purchasingID)

	if s.webhookService != nil {
		_ = s.webhookService.SendPurchasingNotification(existingPurchasing)
	}

	model.GrandTotal = grandTotal
	return model
}

func (s *purchasingService) Delete(ctx context.Context, id string) {
	var purchasing models.Purchasing
	if err := config.DB.Preload("Details").First(&purchasing, id).Error; err != nil {
		exception.PanicLogging(fmt.Errorf("Purchasing not found: %w", err))
	}

	// Subtract stock
	for _, detail := range purchasing.Details {
		if err := config.DB.Model(&models.Item{}).
			Where("id = ?", detail.ItemID).
			UpdateColumn("stock", gorm.Expr("stock - ?", detail.Qty)).
			Error; err != nil {
			exception.PanicLogging(fmt.Errorf("Failed to update stock for item %d: %w", detail.ItemID, err))
		}
	}

	// Soft delete
	config.DB.Delete(&purchasing)
}

func (s *purchasingService) FindById(ctx context.Context, id string) models.Purchasing {
	var purchasing models.Purchasing
	config.DB.Preload("Supplier").Preload("User").Preload("Details.Item").First(&purchasing, id)
	return purchasing
}

func (s *purchasingService) FindAll(ctx context.Context) []models.Purchasing {
	var purchasings []models.Purchasing
	config.DB.Preload("Supplier").Preload("User").Preload("Details.Item").Find(&purchasings)
	return purchasings
}
