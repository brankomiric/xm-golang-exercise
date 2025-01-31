package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"xm-company/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gookit/event"
)

func (h *Handler) Create(ctx *gin.Context) {
	jsonData, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	var company dto.CreateCompany
	err = json.Unmarshal(jsonData, &company)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	err = company.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	result, err := h.DB.AddCompany(company)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to add company"})
		return
	}
	event.MustFire("company_created", event.M{"id": result.String()})
	ctx.JSON(http.StatusOK, UuidResponse{CompanyID: result.String()})
}

func (h *Handler) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	jsonData, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	var company dto.UpdateCompany
	err = json.Unmarshal(jsonData, &company)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	err = company.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	err = h.DB.UpdateCompany(id, company)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to update company"})
		return
	}
	event.MustFire("company_updated", event.M{"id": id})
	ctx.JSON(http.StatusOK, gin.H{"message": "Company updated"})
}

func (h *Handler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	err = h.DB.DeleteCompany(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to delete company"})
		return
	}
	event.MustFire("company_deleted", event.M{"id": id})
	ctx.JSON(http.StatusOK, gin.H{"message": "Company deleted"})
}

func (h *Handler) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	company, err := h.DB.GetCompanyByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: "Company not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to get company"})
		return
	}
	ctx.JSON(http.StatusOK, company)
}
