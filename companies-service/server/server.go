package server

import (
	"xm-company/internal/database"
	"xm-company/server/handlers"
	"xm-company/server/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db database.AbstractDB, isDevMode bool) *gin.Engine {
	mode := gin.ReleaseMode
	if isDevMode {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)
	r := gin.New()
	// This way healthcheck logs will not be printed
	r.Use(gin.LoggerWithWriter(gin.DefaultWriter, "/health"), gin.Recovery())

	h := handlers.Handler{DB: db}

	r.GET("/health", h.Healthcheck)
	r.POST("/company", middleware.JWTAuthMiddleware(), h.Create)
	r.PATCH("/company/:id", middleware.JWTAuthMiddleware(), h.Update)
	r.DELETE("/company/:id", middleware.JWTAuthMiddleware(), h.Delete)
	r.GET("/company/:id", h.GetByID)

	return r
}
