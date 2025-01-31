package database

import "github.com/google/uuid"

type Company struct {
	ID                uuid.UUID `json:"id" db:"id"`
	Name              string    `json:"name" db:"name"`
	Description       *string   `json:"description,omitempty" db:"description"`
	AmountOfEmployees int       `json:"amount_of_employees" db:"amount_of_employees"`
	Registered        bool      `json:"registered" db:"registered"`
	Type              string    `json:"type" db:"type"`
}
