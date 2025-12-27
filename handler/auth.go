package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/bariq12/bookingticket/models"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)
var validate = validator.New()
type AuthHandler struct {
	service models.AuthService
}

func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	Creds := &models.AuthCrendetial{}
	context, cancel := context.WithTimeout(context.Background(),time.Duration(5*time.Second))
	defer cancel()

	if err := ctx.BodyParser(&Creds); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status" : "fail",
			"message": err.Error(),
			"data" : nil,
		})
	}

	if err := validate.Struct(Creds); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status" : "fail",
			"message": err.Error(),
		})
	}

	token, user, err := h.service.Login(context, Creds)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"status" : "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status" : "success",
		"message": "login successful",	
		"data" : fiber.Map{
			"token": token,
			"user": user,
		},
	})
}

func (h *AuthHandler) Register(ctx *fiber.Ctx) error {
	Creds := &models.AuthCrendetial{}
	context, cancel := context.WithTimeout(context.Background(),time.Duration(5*time.Second))
	defer cancel()

	if err := ctx.BodyParser(&Creds); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status" : "fail",
			"message": err.Error(),
			"data" : nil,
		})
	}

	if err := validate.Struct(Creds); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status" : "fail",
			"message": fmt.Errorf("invalid input: %v", err).Error(),
		})
	}

	token, user, err := h.service.Register(context, Creds)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"status" : "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status" : "success",
		"message": "login successful",	
		"data" : fiber.Map{
			"token": token,
			"user": user,
		},
	})
}

		
func NewAuthHandler(router fiber.Router, service models.AuthService) {
	Handler := &AuthHandler{
		service: service,	
	}
	router.Post("/Login", Handler.Login)
	router.Post("/Register", Handler.Register)
}