package handlers

import (
	"xm-company/internal/database"
)

type Handler struct {
	DB database.AbstractDB
}

type HealthcheckResponse struct {
	Services Services `json:"services"`
}

type Services struct {
	Postgres HealthcheckMessage `json:"Postgres"`
}

type HealthcheckMessage struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type UuidResponse struct {
	CompanyID string `json:"company_id"`
}
