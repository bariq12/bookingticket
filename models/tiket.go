package models

import (
	"context"
	"time"
)

type Ticket struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	EventID   uint      `json:"eventId"`
	UserID    uint    `json:"userId" gorm:"foreignkey:UserID;contraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Event     Event     `json:"event" gorm:"foreignkey:EventID;contraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Entered   bool      `json:"entered" default:"false"`
	CreatedAt time.Time`json:"createdat"`
	UpdateAt  time.Time `json:"updateat"`
}

type TicketRepository interface{
	GetMany(ctx context.Context,UserID uint) ([]*Ticket,error)
	GetOne(ctx context.Context,UserID uint ,ticketID uint) (*Ticket,error)
	CreateOne(ctx context.Context,UserID uint, ticket *Ticket) (*Ticket,error)
	UpdateOne(ctx context.Context,UserID uint, ticketId uint, updateData map[string]interface{}) (*Ticket,error)
}

type ValidateTicket struct{
	TicketId uint `json:"ticketId"`
	OwnerId uint `json:"ownerId"`
}