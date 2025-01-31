package dto

import (
	"errors"
	"strings"
)

type CompanyType string

const (
	Corporations       CompanyType = "Corporations"
	NonProfit          CompanyType = "NonProfit"
	Cooperative        CompanyType = "Cooperative"
	SoleProprietorship CompanyType = "Sole Proprietorship"
)

type CreateCompany struct {
	Name              string      `json:"name"`
	Description       *string     `json:"description,omitempty"`
	AmountOfEmployees int         `json:"amount_of_employees"`
	Registered        bool        `json:"registered"`
	Type              CompanyType `json:"type"`
}

func (c *CreateCompany) Validate() error {
	if len(strings.TrimSpace(c.Name)) == 0 || len(c.Name) > 15 {
		return errors.New("name is required and must be 15 characters or fewer")
	}

	if c.AmountOfEmployees < 0 {
		return errors.New("amount_of_employees must be a non-negative integer")
	}

	if !isValidCompanyType(c.Type) {
		return errors.New("type must be one of: Corporations, NonProfit, Cooperative, Sole Proprietorship")
	}

	return nil
}

func isValidCompanyType(t CompanyType) bool {
	switch t {
	case Corporations, NonProfit, Cooperative, SoleProprietorship:
		return true
	default:
		return false
	}
}

type UpdateCompany struct {
	Name              *string      `json:"name,omitempty"`
	Description       *string      `json:"description,omitempty"`
	AmountOfEmployees *int         `json:"amount_of_employees,omitempty"`
	Registered        *bool        `json:"registered,omitempty"`
	Type              *CompanyType `json:"type,omitempty"`
}

func (u *UpdateCompany) Validate() error {
	if u.Name != nil && len(strings.TrimSpace(*u.Name)) > 15 {
		return errors.New("name must be 15 characters or fewer")
	}

	if u.AmountOfEmployees != nil && *u.AmountOfEmployees < 0 {
		return errors.New("amount_of_employees must be a non-negative integer")
	}

	if u.Type != nil && !isValidCompanyType(*u.Type) {
		return errors.New("type must be one of: Corporations, NonProfit, Cooperative, Sole Proprietorship")
	}

	return nil
}
