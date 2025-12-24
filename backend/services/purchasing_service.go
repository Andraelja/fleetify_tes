package services

import (
	"backend/app/models"
	"context"
)

type PurchasingService interface {
	Create(
		ctx context.Context,
		model models.PurchasingCreateOrUpdateModel,
	) models.Purchasing

	Update(
		ctx context.Context,
		model models.PurchasingCreateOrUpdateModel,
		id string,
	) models.Purchasing

	Delete(ctx context.Context, id string)

	FindById(ctx context.Context, id string) models.Purchasing

	FindAll(ctx context.Context) []models.Purchasing
}
