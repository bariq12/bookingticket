package handler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/bariq12/bookingticket/models"
	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
)

type TicketHandler struct { 
	repository models.TicketRepository
}

func (h *TicketHandler) GetMany(ctx *fiber.Ctx) error{
	context, cancel := context.WithTimeout(context.Background(),time.Duration(5*time.Second))
	defer cancel()
	userID := uint(ctx.Locals("userId").(float64))
	tickets, err := h.repository.GetMany(context, userID)

	if err != nil {
		return ctx.Status(fiber.StatusBadGateway).JSON(&fiber.Map{
			"status" : "fail",
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status" : "success",
		"message": "",
		"data" : tickets,
	})
}

func (h *TicketHandler) GetOne(ctx *fiber.Ctx) error{
	context, cancel := context.WithTimeout(context.Background(),time.Duration(5*time.Second))
	defer cancel()

	ticketId, _ := strconv.Atoi(ctx.Params("ticketId"))
	userID := uint(ctx.Locals("userId").(float64))

	ticket, err := h.repository.GetOne(context,userID ,uint(ticketId))

	
	if err != nil {
		return ctx.Status(fiber.StatusBadGateway).JSON(&fiber.Map{
			"status" : "fail",
			"message": err.Error(),
		})
	}

	var QRCode []byte
	QRCode, err = qrcode.Encode(
		fmt.Sprintf("ticketId: %v","ownerId: %v",ticket.ID, userID),
		qrcode.Medium,
		256,
	)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status" : "fail",
			"message": "Failed to generate QR code",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status" : "success",
		"message": "",
		"data" : &fiber.Map{
			"ticket": ticket,
			"QRCode": QRCode,
		},
	})
}

func (h *TicketHandler) CreateOne(ctx *fiber.Ctx) error{
	context, cancel := context.WithTimeout(context.Background(),time.Duration(5*time.Second))
	defer cancel()

	ticket := &models.Ticket{}

	if err := ctx.BodyParser(ticket); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status" : "fail",
			"message": err.Error(),
			"data" : nil,
		})
	}
	userId := uint(ctx.Locals("userId").(float64))

	ticket, err := h.repository.CreateOne(context,userId, ticket)

	
	if err != nil {
		return ctx.Status(fiber.StatusBadGateway).JSON(&fiber.Map{
			"status" : "fail",
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"status" : "success",
		"message": "",
		"data" : ticket,
	})
}
func (h *TicketHandler) ValidateOne(ctx *fiber.Ctx) error{
	context, cancel := context.WithTimeout(context.Background(),time.Duration(5*time.Second))
	defer cancel()

	validateBody := &models.ValidateTicket{}

	if err := ctx.BodyParser(validateBody); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status" : "fail",
			"message": err.Error(),
			"data" : nil,
		})
	}
	ValidateData := make(map[string]interface{})
	ValidateData["entered"] = true

	ticket, err := h.repository.UpdateOne(context,validateBody.OwnerId, validateBody.TicketId, ValidateData)

	
	if err != nil {
		return ctx.Status(fiber.StatusBadGateway).JSON(&fiber.Map{
			"status" : "fail",
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status" : "success",
		"message": "Welcome the Show!",
		"data" : ticket,
	})
}

func NewTicketHandler(router fiber.Router, repository models.TicketRepository){
	handler := TicketHandler{
		repository: repository,
	}

	router.Get("/", handler.GetMany)
	router.Post("/", handler.CreateOne)
	router.Get("/:ticketId", handler.GetOne)
	router.Post("/validate", handler.ValidateOne)
}

