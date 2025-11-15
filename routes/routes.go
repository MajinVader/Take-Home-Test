package routes

import (
	"take-home-test/handlers"
	"take-home-test/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Auth
	app.Post("/register", handlers.Register)
	app.Post("/login", handlers.Login)

	// Fields
	app.Get("/fields", handlers.GetFields)
	app.Get("/fields/:id", handlers.GetFieldByID)
	app.Post("/fields", middlewares.AdminOnly, handlers.CreateField)
	app.Put("/fields/:id", middlewares.AdminOnly, handlers.UpdateField)
	app.Delete("/fields/:id", middlewares.AdminOnly, handlers.DeleteField)

	// Bookings
	app.Post("/bookings", middlewares.AuthRequired, handlers.CreateBooking)
	app.Get("/my-bookings", middlewares.AuthRequired, handlers.GetMyBookings)
	app.Get("/bookings", middlewares.AdminOnly, handlers.GetAllBookings)

	// Payments (mock)
	app.Post("/payments", middlewares.AuthRequired, handlers.PayBooking)
}
