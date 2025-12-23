package services

import (
	"backend/app/models"
	"backend/config"
	"context"
	"strconv"
)

type supplierService struct {
	// nanti isi repository/db
}

func NewSupplierService() SupplierService {
	return &supplierService{}
}

func (s *supplierService) Create(
	ctx context.Context,
	model models.SupplierCreateOrUpdateModel,
) models.SupplierCreateOrUpdateModel {
	supplier := models.Supplier{
		Name:    model.Name,
		Email:   model.Email,
		Address: model.Address,
	}
	config.DB.Create(&supplier)
	return model
}

func (s *supplierService) Update(
	ctx context.Context,
	model models.SupplierCreateOrUpdateModel,
	id string,
) models.SupplierCreateOrUpdateModel {
	var supplier models.Supplier
	supplierID, _ := strconv.Atoi(id)
	config.DB.First(&supplier, supplierID)
	supplier.Name = model.Name
	supplier.Email = model.Email
	supplier.Address = model.Address
	config.DB.Save(&supplier)
	return model
}

func (s *supplierService) Delete(ctx context.Context, id string) {
	var supplier models.Supplier
	supplierID, _ := strconv.Atoi(id)
	config.DB.Delete(&supplier, supplierID)
}

func (s *supplierService) FindById(
	ctx context.Context,
	id string,
) models.Supplier {
	var supplier models.Supplier
	supplierID, _ := strconv.Atoi(id)
	config.DB.First(&supplier, supplierID)
	return supplier
}

func (s *supplierService) FindAll(
	ctx context.Context,
) []models.Supplier {
	var suppliers []models.Supplier
	config.DB.Find(&suppliers)
	return suppliers
}
