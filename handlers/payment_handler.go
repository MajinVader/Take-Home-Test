package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"take-home-test/database"
	"take-home-test/models"
)

type PaymentInput struct {
	BookingID uint `json:"booking_id"`
}

// POST /payments
func PayBooking(c *fiber.Ctx) error {
	var input PaymentInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	if input.BookingID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "booking_id is required",
		})
	}

	// ambil user dari JWT
	userIDAny := c.Locals("userID")
	var userID uint
	if f, ok := userIDAny.(float64); ok {
		userID = uint(f)
	} else if u, ok := userIDAny.(uint); ok {
		userID = u
	} else {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user in context"})
	}

	roleAny := c.Locals("role")
	role, _ := roleAny.(string)

	// cari booking
	var booking models.Booking
	if err := database.DB.First(&booking, input.BookingID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "Booking not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get booking"})
	}

	// hanya pemilik booking atau admin yang boleh bayar
	if booking.UserID != userID && role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "You are not allowed to pay this booking",
		})
	}

	if booking.Status == "paid" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Booking already paid",
		})
	}

	// payment gateway sukses (mock)
	booking.Status = "paid"

	if err := database.DB.Save(&booking).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update booking status",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Payment successful",
		"booking": booking,
	})
}
