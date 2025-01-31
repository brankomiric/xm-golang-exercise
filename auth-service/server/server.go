package server

import (
	"xm-auth/internal/database"
	"xm-auth/server/handlers"

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
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)

	return r
}
