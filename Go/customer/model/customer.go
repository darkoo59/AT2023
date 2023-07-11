package model

import (
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"net/mail"
)

type Customer struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Email    string    `gorm:"uniqueIndex;not null"`
	Password string    `gorm:"not null"`
	Balance  float64
}

type Customers []*Customers

func (customer *Customer) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(customer)
}

func (customer *Customer) FromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(customer)
}

func (customer *Customer) IsEmailValid() bool {
	_, err := mail.ParseAddress(customer.Email)
	return err == nil
}
