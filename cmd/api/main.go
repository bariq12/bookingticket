package main

import (
	"github.com/gofiber/fiber/v2"

	"github.com/bariq12/bookingticket/config"
	"github.com/bariq12/bookingticket/db"
	"github.com/bariq12/bookingticket/handler"
	"github.com/bariq12/bookingticket/middleware"
	"github.com/bariq12/bookingticket/repositories"
	"github.com/bariq12/bookingticket/services"
)

func main() {
	EnvConfig := config.NewENVConfig()
	db := db.Init(EnvConfig, db.DBMigrator)
	app := fiber.New(fiber.Config{
		AppName:      "TicketBooking",
		ServerHeader: "Fiber",
	})

	// Repository initialization
	eventRepo := repositories.NewEventRepository(db)
	ticketRepo := repositories.NewTicketRepository(db)
	authRepo := repositories.NewAuthRepository(db)

	// services
	AuthServices := services.NewAuthService(authRepo)

	// Group base routing under /api
	api := app.Group("/api")
	handler.NewAuthHandler(api.Group("/auth"), AuthServices)

	PrivateRoute := api.Use(middleware.Authprotected(db))

	// Register event routes under /api/event
	handler.NewEventHandler(PrivateRoute.Group("/event"), eventRepo)
	handler.NewTicketHandler(PrivateRoute.Group("/ticket"), ticketRepo)

	// Start server
	if err := app.Listen(":" + EnvConfig.ServerPort); err != nil {
	panic(err)
	}
}
