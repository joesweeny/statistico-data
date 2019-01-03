package model

import (
	"github.com/satori/go.uuid"
	"time"
)

type Country struct {
	ID         uuid.UUID `json:"id"`
	ExternalID int       `json:"external_id"`
	Name       string    `json:"name"`
	Continent  string    `json:"continent"`
	ISO        string    `json:"iso"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
