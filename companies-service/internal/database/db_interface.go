package database

import (
	"xm-company/internal/dto"

	"github.com/google/uuid"
)

type AbstractDB interface {
	TestConn() error
	AddCompany(input dto.CreateCompany) (*uuid.UUID, error)
	UpdateCompany(id uuid.UUID, companyDto dto.UpdateCompany) error
	DeleteCompany(id uuid.UUID) error
	GetCompanyByID(id uuid.UUID) (*Company, error)
}
