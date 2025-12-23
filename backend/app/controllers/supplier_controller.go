package controllers

import (
	"backend/app/models"
	"backend/exception"
	"backend/services"
	"github.com/gofiber/fiber/v2"
)

type SupplierController struct {
	Service services.SupplierService
}

func NewSupplierController(service services.SupplierService) *SupplierController {
	return &SupplierController{Service: service}
}

func (controller SupplierController) Create(c *fiber.Ctx) error {
	var request models.SupplierCreateOrUpdateModel
	err := c.BodyParser(&request)
	exception.PanicLogging(err)

	response := controller.Service.Create(c.Context(), request)
	return c.Status(fiber.StatusCreated).JSON(models.GeneralResponse{
		Code:    201,
		Message: "Success create data!",
		Data:    response,
	})
}

func (controller SupplierController) GetAll(c *fiber.Ctx) error {
	response := controller.Service.FindAll(c.Context())
	return c.Status(fiber.StatusOK).JSON(models.GeneralResponse{
		Code:    200,
		Message: "Success get data!",
		Data:    response,
	})
}

func (controller SupplierController) GetById(c *fiber.Ctx) error {
	id := c.Params("id")
	response := controller.Service.FindById(c.Context(), id)
	return c.Status(fiber.StatusOK).JSON(models.GeneralResponse{
		Code:    200,
		Message: "Success get data by id!",
		Data:    response,
	})
}

func (controller SupplierController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var request models.SupplierCreateOrUpdateModel
	err := c.BodyParser(&request)
	exception.PanicLogging(err)

	response := controller.Service.Update(c.Context(), request, id)
	return c.Status(fiber.StatusOK).JSON(models.GeneralResponse{
		Code:    200,
		Message: "Success update data!",
		Data:    response,
	})
}

func (controller SupplierController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	controller.Service.Delete(c.Context(), id)
	return c.Status(fiber.StatusOK).JSON(models.GeneralResponse{
		Code:    200,
		Message: "Success delete data!",
		Data:    nil,
	})
}
