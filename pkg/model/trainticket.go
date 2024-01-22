package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/priyanshu/trainservice-database/pkg/config"
)

type User struct {
	gorm.Model
	UserID    string `gorm:"column:user_id" json:"user_id"`
	FirstName string `gorm:"column:first_name" json:"first_name"`
	LastName  string `gorm:"column:last_name" json:"last_name"`
	Email     string `gorm:"column:email" json:"email"`
}

type Ticket struct {
	gorm.Model
	From      string  `gorm:"column:from" json:"from"`
	To        string  `gorm:"column:to" json:"to"`
	UserID    string  `gorm:"column:user_id" json:"user_id"`
	PricePaid float32 `gorm:"column:price_paid" json:"price_paid"`
	Seat      string  `gorm:"column:seat" json:"seat"`
}

var db *gorm.DB

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&User{}, &Ticket{})
	fmt.Println("AutoMigrate completed")
}

func (user *User) CreateUser() *User {
	db.NewRecord(user)
	if err := db.Create(&user).Error; err != nil {
		// Handle the error, e.g., log it or return an error
		fmt.Println(err)
		return nil
	}
	return user
}

func (ticket *Ticket) CreateTicket() *Ticket {
	db.NewRecord(ticket)
	if err := db.Create(&ticket).Error; err != nil {
		// Handle the error, e.g., log it or return an error
		fmt.Println(err)
		return nil
	}
	return ticket
}



func RemoveTicket(userID string) error {
	var ticket Ticket
	err := db.Where("user_id = ?", userID).Delete(&ticket).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateUserSection(userID, newSection string) error {
	var user User
	if newSection == "A" {
		err := db.Model(&user).Where("user_id = ?", userID).Update("seat", "B").Error
		if err != nil {
			return err
		}
	} else {
		// Handle the other condition
		// Example: update to a different seat or perform other logic
		err := db.Model(&user).Where("user_id = ?", userID).Update("seat", "A").Error
		if err != nil {
			return err
		}
	}
	return nil
}
