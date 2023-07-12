package model

import (
	"encoding/json"
	"github.com/google/uuid"
	"io"
)

type Order struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	ItemId         string    `gorm:"not null"`
	UserId         string    `gorm:"not null"`
	Quantity       int       `gorm:"not null"`
	ItemName       string
	PricePerItem   float64 `gorm:"not null"`
	OrderStatus    string  `gorm:"not null"`
	AccountBalance float64 `gorm:"not null"`
}

type Orders []*Orders

func (order *Order) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(order)
}

func (order *Order) FromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(order)
}
