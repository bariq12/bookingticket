package db

import (
	"github.com/bariq12/bookingticket/models"
	"gorm.io/gorm"
)

func DBMigrator(db *gorm.DB) error{
	return db.AutoMigrate(&models.Event{},&models.Ticket{},&models.User{})
}