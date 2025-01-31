package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"xm-auth/internal/database"
	"xm-auth/internal/utils"
	"xm-auth/server/handlers"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestHealthcheck(t *testing.T) {
	db := database.NewMockAbstractDB(t)

	router := SetupRouter(db, true)

	db.On("TestConn", mock.Anything).Return(nil).Once()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)

	router.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)
	var resp handlers.HealthcheckResponse
	json.Unmarshal(responseData, &resp)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "OK", resp.Services.Mongo.Status)
}

func TestRegisterSuccess(t *testing.T) {
	db := database.NewMockAbstractDB(t)

	router := SetupRouter(db, true)

	db.On("ValidateUsername", mock.Anything, mock.Anything).Return(false, nil).Once()

	id := bson.NewObjectID()
	insertOneResult := &mongo.InsertOneResult{
		InsertedID:   id,
		Acknowledged: true,
	}
	db.On("CreateUser", mock.Anything, mock.Anything).Return(insertOneResult, nil).Once()

	user := `{"username": "test", "password": "test"}`
	reqBody := []byte(user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", nil)
	req.Body = io.NopCloser(bytes.NewReader(reqBody))

	router.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)
	var resp handlers.RegisterSuccessResponse
	json.Unmarshal(responseData, &resp)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, id, resp.CreatedUserId)
}

func TestRegisterUsernameExists(t *testing.T) {
	db := database.NewMockAbstractDB(t)

	router := SetupRouter(db, true)

	db.On("ValidateUsername", mock.Anything, mock.Anything).Return(true, nil).Once()

	user := `{"username": "test", "password": "test"}`
	reqBody := []byte(user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", nil)
	req.Body = io.NopCloser(bytes.NewReader(reqBody))

	router.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)
	var resp handlers.ErrorResponse
	json.Unmarshal(responseData, &resp)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Equal(t, "Username is already taken", resp.Message)
}

func TestLoginSuccess(t *testing.T) {
	os.Setenv("JWT_SECRET", "test")

	db := database.NewMockAbstractDB(t)

	router := SetupRouter(db, true)

	password := "test"

	passwordHash, _ := utils.HashPassword(password)

	dbUser := bson.M{"username": "test", "password": passwordHash}

	db.On("GetUser", mock.Anything, mock.Anything).Return(dbUser, nil).Once()

	user := `{"username": "test", "password": "test"}`
	reqBody := []byte(user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", nil)
	req.Body = io.NopCloser(bytes.NewReader(reqBody))

	router.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)
	var resp handlers.JWTResponse
	json.Unmarshal(responseData, &resp)

	require.Equal(t, http.StatusOK, w.Code)
	isValid, err := utils.IsValidJWT(resp.Token)
	require.True(t, isValid)
	require.Nil(t, err)
}

func TestLoginInvalidUsername(t *testing.T) {
	os.Setenv("JWT_SECRET", "test")

	db := database.NewMockAbstractDB(t)

	router := SetupRouter(db, true)

	db.On("GetUser", mock.Anything, mock.Anything).Return(nil, nil).Once()

	user := `{"username": "test", "password": "test"}`
	reqBody := []byte(user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", nil)
	req.Body = io.NopCloser(bytes.NewReader(reqBody))

	router.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)
	var resp handlers.ErrorResponse
	json.Unmarshal(responseData, &resp)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Equal(t, "Invalid Username or Password", resp.Message)
}

func TestLoginInvalidPassword(t *testing.T) {
	os.Setenv("JWT_SECRET", "test")

	db := database.NewMockAbstractDB(t)

	router := SetupRouter(db, true)

	password := "test"

	passwordHash, _ := utils.HashPassword(password)

	dbUser := bson.M{"username": "test", "password": passwordHash}

	db.On("GetUser", mock.Anything, mock.Anything).Return(dbUser, nil).Once()

	user := `{"username": "test", "password": "wrong_password"}` // incorrect password
	reqBody := []byte(user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", nil)
	req.Body = io.NopCloser(bytes.NewReader(reqBody))

	router.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)
	var resp handlers.ErrorResponse
	json.Unmarshal(responseData, &resp)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	require.Equal(t, "Invalid Username or Password", resp.Message)
}
