package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Healthcheck(ctx *gin.Context) {
	httpStatus := http.StatusOK
	resp := HealthcheckResponse{
		Services{
			Postgres: HealthcheckMessage{Status: "OK"},
		},
	}

	err := h.DB.TestConn()
	if err != nil {
		log.Printf("health check failed: Postgres: %s\n", err.Error())

		resp.Services.Postgres.Status = "Not OK"
		httpStatus = http.StatusInternalServerError
	}

	ctx.JSON(httpStatus, resp)
}
