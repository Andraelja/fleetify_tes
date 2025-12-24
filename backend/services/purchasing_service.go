package services

import (
	"backend/app/models"
	"context"
)

type PurchasingService interface {
	Create(ctx context.Context, models models.PurchasingCreateOrUpdateModel) models.PurchasingCreateOrUpdateModel
	Update(ctx context.Context, Purchasing models.PurchasingCreateOrUpdateModel, id string) models.PurchasingCreateOrUpdateModel
	Delete(ctx context.Context, id string)
	FindById(ctx context.Context, id string) models.Purchasing
	FindAll(ctx context.Context) []models.Purchasing
}
