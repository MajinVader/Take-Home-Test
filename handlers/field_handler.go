package handlers

import (
	"strconv"

	"take-home-test/database"
	"take-home-test/models"

	"github.com/gofiber/fiber/v2"
)

type FieldInput struct {
	Name         string `json:"name"`
	PricePerHour int    `json:"price_per_hour"`
	Location     string `json:"location"`
}

func CreateField(c *fiber.Ctx) error {
	var input FieldInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	field := models.Field{
		Name:         input.Name,
		PricePerHour: input.PricePerHour,
		Location:     input.Location,
	}

	if err := database.DB.Create(&field).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create field"})
	}

	return c.Status(201).JSON(field)
}

func GetFields(c *fiber.Ctx) error {
	var fields []models.Field
	if err := database.DB.Find(&fields).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get fields"})
	}
	return c.JSON(fields)
}

func GetFieldByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var field models.Field

	if err := database.DB.First(&field, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Field not found"})
	}

	return c.JSON(field)
}

func UpdateField(c *fiber.Ctx) error {
	id := c.Params("id")
	var field models.Field

	if err := database.DB.First(&field, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Field not found"})
	}

	var input FieldInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	field.Name = input.Name
	field.PricePerHour = input.PricePerHour
	field.Location = input.Location

	database.DB.Save(&field)

	return c.JSON(field)
}

func DeleteField(c *fiber.Ctx) error {
	id := c.Params("id")

	if _, err := strconv.Atoi(id); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := database.DB.Delete(&models.Field{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete field"})
	}

	return c.JSON(fiber.Map{"message": "Field deleted"})
}
