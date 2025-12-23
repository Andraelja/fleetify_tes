package services

import (
	"backend/app/models"
	"backend/config"
	"context"
	"strconv"
)

type itemService struct {
	// nanti isi repository/db
}

func NewItemService() ItemService {
	return &itemService{}
}

func (s *itemService) Create(
	ctx context.Context,
	model models.ItemCreateOrUpdateModel,
) models.ItemCreateOrUpdateModel {
	item := models.Item{
		Name:  model.Name,
		Stock: model.Stock,
		Price: model.Price,
	}
	config.DB.Create(&item)
	return model
}

func (s *itemService) Update(
	ctx context.Context,
	model models.ItemCreateOrUpdateModel,
	id string,
) models.ItemCreateOrUpdateModel {
	var item models.Item
	itemID, _ := strconv.Atoi(id)
	config.DB.First(&item, itemID)
	item.Name = model.Name
	item.Stock = model.Stock
	item.Price = model.Price
	config.DB.Save(&item)
	return model
}

func (s *itemService) Delete(ctx context.Context, id string) {
	var item models.Item
	itemID, _ := strconv.Atoi(id)
	config.DB.Delete(&item, itemID)
}

func (s *itemService) FindById(
	ctx context.Context,
	id string,
) models.Item {
	var item models.Item
	itemID, _ := strconv.Atoi(id)
	config.DB.First(&item, itemID)
	return item
}

func (s *itemService) FindAll(
	ctx context.Context,
) []models.Item {
	var item []models.Item
	config.DB.Find(&item)
	return item
}
