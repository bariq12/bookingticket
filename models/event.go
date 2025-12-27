package models

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Event struct {
	ID       uint`json:"id" gorm:"primarykey"`
	Name     string	`json:"name"`
	Location string	`json:"location"`
	TotalTicketPurchased int64	`json:"totalticketpurchased" gorm:"-"`
	TotalTicketEntered int64	`json:"totalticketentered" gorm:"-"`
	Date     time.Time `json:"date"`
	CreatedAt   time.Time	`json:"createdat"`
	UpdateAt     time.Time	`json:"updateat"`

}

type EventRepository interface {
	GetMany(ctx context.Context) ([]*Event,error)
	GetOne(ctx context.Context, eventID uint) (*Event,error)
	CreateOne(ctx context.Context, event *Event) (*Event,error)
	UpdateOne(ctx context.Context, event uint, updateData map[string]interface{}) (*Event,error)
	DeleteOne(ctx context.Context, eventID uint) error
}

func (e *Event) AfterFind(db *gorm.DB) (err error) {
	baseQuery := db.Model(&Ticket{}).Where(&Ticket{EventID: e.ID})

	if res := baseQuery.Count(&e.TotalTicketPurchased); res.Error != nil {
		return res.Error
	}
	if res := baseQuery.Where("entered = ?", true).Count(&e.TotalTicketEntered); res.Error != nil {
		return res.Error
	}
	return nil
}