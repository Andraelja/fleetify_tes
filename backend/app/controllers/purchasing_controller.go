package controllers

import (
	"backend/app/models"
	"backend/exception"
	"backend/services"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type PurchasingController struct {
	Service services.PurchasingService
}

func NewPurchasingController(service services.PurchasingService) *PurchasingController {
	return &PurchasingController{Service: service}
}

func (controller *PurchasingController) Create(c *fiber.Ctx) error {
	var request models.PurchasingCreateOrUpdateModel
	err := c.BodyParser(&request)
	exception.PanicLogging(err)

	defer func() {
		if r := recover(); r != nil {
			c.Status(fiber.StatusBadRequest).JSON(models.GeneralResponse{
				Code:    400,
				Message: "Failed to create purchasing transaction: " + fmt.Sprintf("%v", r),
				Data:    nil,
			})
		}
	}()

	response := controller.Service.Create(c.Context(), request)
	return c.Status(fiber.StatusCreated).JSON(models.GeneralResponse{
		Code:    201,
		Message: "Success create purchasing transaction!",
		Data:    response,
	})
}

func (controller *PurchasingController) Update(c *fiber.Ctx) error {
	var request models.PurchasingCreateOrUpdateModel
	err := c.BodyParser(&request)
	exception.PanicLogging(err)

	defer func() {
		if r := recover(); r != nil {
			c.Status(fiber.StatusBadRequest).JSON(models.GeneralResponse{
				Code:    400,
				Message: "Failed to update purchasing transaction: " + fmt.Sprintf("%v", r),
				Data:    nil,
			})
		}
	}()

	id := c.Params("id")
	response := controller.Service.Update(c.Context(), request, id)
	return c.Status(fiber.StatusOK).JSON(models.GeneralResponse{
		Code:    200,
		Message: "Success update purchasing transaction!",
		Data:    response,
	})
}

func (controller *PurchasingController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	controller.Service.Delete(c.Context(), id)
	return c.Status(fiber.StatusOK).JSON(models.GeneralResponse{
		Code:    200,
		Message: "Success delete purchasing transaction!",
		Data:    nil,
	})
}

func (controller *PurchasingController) FindById(c *fiber.Ctx) error {
	id := c.Params("id")
	response := controller.Service.FindById(c.Context(), id)
	return c.Status(fiber.StatusOK).JSON(models.GeneralResponse{
		Code:    200,
		Message: "Success get purchasing by id!",
		Data:    response,
	})
}

func (controller *PurchasingController) FindAll(c *fiber.Ctx) error {
	response := controller.Service.FindAll(c.Context())
	return c.Status(fiber.StatusOK).JSON(models.GeneralResponse{
		Code:    200,
		Message: "Success get all purchasing!",
		Data:    response,
	})
}
