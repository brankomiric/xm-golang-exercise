package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"xm-auth/internal/database"
	"xm-auth/internal/structs"
	"xm-auth/internal/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Handler struct {
	DB database.AbstractDB
}

func (h *Handler) Healthcheck(ctx *gin.Context) {
	httpStatus := http.StatusOK
	resp := HealthcheckResponse{
		Services{
			Mongo: HealthcheckMessage{Status: "OK"},
		},
	}

	err := h.DB.TestConn(ctx)
	if err != nil {
		log.Printf("health check failed: Mongo: %s\n", err.Error())

		resp.Services.Mongo.Status = "Not OK"
		httpStatus = http.StatusInternalServerError
	}

	ctx.JSON(httpStatus, resp)
}

func (h *Handler) Register(ctx *gin.Context) {
	jsonData, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	var user structs.User
	err = json.Unmarshal(jsonData, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	// Check If Username is taken
	exists, err := h.DB.ValidateUsername(ctx, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	if exists {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: "Username is already taken"})
		return
	}

	// Hash Password
	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	// Insert User
	created, err := h.DB.CreateUser(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	createdUserID, ok := created.InsertedID.(bson.ObjectID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to assert type of InsertedID"})
		return
	}
	ctx.JSON(http.StatusOK, RegisterSuccessResponse{CreatedUserId: createdUserID})
}

func (h *Handler) Login(ctx *gin.Context) {
	jsonData, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	var user structs.User
	err = json.Unmarshal(jsonData, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	// Get User
	dbUser, err := h.DB.GetUser(ctx, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	if dbUser == nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid Username or Password"})
		return
	}

	// Compare Password
	isValidPassword := utils.CheckPassword(dbUser["password"].(string), user.Password)
	if !isValidPassword {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Invalid Username or Password"})
		return
	}

	// Generate JWT
	token, err := utils.CreateJWT(dbUser["username"].(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, JWTResponse{Token: token})
}
