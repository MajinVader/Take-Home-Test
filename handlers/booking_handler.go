package handlers

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"take-home-test/database"
	"take-home-test/models"
)

type BookingInput struct {
	FieldID   uint   `json:"field_id"`
	StartTime string `json:"start_time"` // RFC3339 string
	EndTime   string `json:"end_time"`
}

// POST /bookings (user login)
func CreateBooking(c *fiber.Ctx) error {
	var input BookingInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if input.FieldID == 0 || input.StartTime == "" || input.EndTime == "" {
		return c.Status(400).JSON(fiber.Map{"error": "field_id, start_time, end_time required"})
	}

	// Parse waktu format RFC3339
	start, err := time.Parse(time.RFC3339, input.StartTime)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid start_time format, use RFC3339"})
	}

	end, err := time.Parse(time.RFC3339, input.EndTime)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid end_time format, use RFC3339"})
	}

	if !end.After(start) {
		return c.Status(400).JSON(fiber.Map{"error": "end_time must be after start_time"})
	}

	// Ambil userID dari JWT (dari middleware)
	userIDAny := c.Locals("userID")
	userID, ok := userIDAny.(uint)
	if !ok {
		// kalau tadi disimpan sebagai float64
		if f, ok2 := userIDAny.(float64); ok2 {
			userID = uint(f)
		} else {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid user in context"})
		}
	}

	// Cek dulu field-nya ada atau tidak
	var field models.Field
	if err := database.DB.First(&field, input.FieldID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Field not found"})
	}

	// Cek overlapping booking di lapangan yang sama
	var existing models.Booking
	err = database.DB.
		Where("field_id = ? AND start < ? AND \"end\" > ?", input.FieldID, end, start).
		First(&existing).Error

	if err == nil {
		// ketemu booking yang bentrok
		return c.Status(400).JSON(fiber.Map{
			"error": "Time slot already booked for this field",
		})
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// error DB lain
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check existing bookings"})
	}

	// Kalau sampai sini berarti tidak bentrok → buat booking baru
	booking := models.Booking{
		UserID:  userID,
		FieldID: input.FieldID,
		Start:   start,
		End:     end,
		Status:  "pending",
	}

	if err := database.DB.Create(&booking).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create booking"})
	}

	return c.Status(201).JSON(booking)
}

// GET /my-bookings (user login)
func GetMyBookings(c *fiber.Ctx) error {
	userIDAny := c.Locals("userID")
	var userID uint
	if f, ok := userIDAny.(float64); ok {
		userID = uint(f)
	} else if u, ok := userIDAny.(uint); ok {
		userID = u
	} else {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user in context"})
	}

	var bookings []models.Booking
	if err := database.DB.Where("user_id = ?", userID).Find(&bookings).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get bookings"})
	}

	return c.JSON(bookings)
}

// GET /bookings (admin only) – lihat semua booking
func GetAllBookings(c *fiber.Ctx) error {
	var bookings []models.Booking
	if err := database.DB.Find(&bookings).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get bookings"})
	}
	return c.JSON(bookings)
}
