package services

import (
	"backend/app/models"
	"context"
)

type SupplierService interface {
	Create(ctx context.Context, models models.SupplierCreateOrUpdateModel) models.SupplierCreateOrUpdateModel
	Update(ctx context.Context, Supplier models.SupplierCreateOrUpdateModel, id string) models.SupplierCreateOrUpdateModel
	Delete(ctx context.Context, id string)
	FindById(ctx context.Context, id string) models.Supplier
	FindAll(ctx context.Context) []models.Supplier
}
