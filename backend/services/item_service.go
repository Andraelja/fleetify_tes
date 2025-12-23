package services

import (
	"backend/app/models"
	"context"
)

type ItemService interface {
	Create(ctx context.Context, models models.ItemCreateOrUpdateModel) models.ItemCreateOrUpdateModel
	Update(ctx context.Context, Item models.ItemCreateOrUpdateModel, id string) models.ItemCreateOrUpdateModel
	Delete(ctx context.Context, id string)
	FindById(ctx context.Context, id string) models.Item
	FindAll(ctx context.Context) []models.Item
}
