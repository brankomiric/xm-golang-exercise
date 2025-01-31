package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"xm-company/internal/database"
	"xm-company/server/handlers"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHealthcheck(t *testing.T) {
	db := database.NewMockAbstractDB(t)

	router := SetupRouter(db, true)

	db.On("TestConn").Return(nil).Once()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)

	router.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)
	var resp handlers.HealthcheckResponse
	json.Unmarshal(responseData, &resp)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "OK", resp.Services.Postgres.Status)
}

func TestCreateCompanySuccess(t *testing.T) {
	os.Setenv("JWT_SECRET", "test")

	db := database.NewMockAbstractDB(t)

	router := SetupRouter(db, true)

	id := uuid.New()

	db.On("AddCompany", mock.Anything).Return(&id, nil).Once()

	company := `{"name": "test", "description": "test", "amount_of_employees": 1, "registered": true, "type": "Corporations"}`
	reqBody := []byte(company)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/company", nil)
	token, _ := createJWT("test", "test")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Body = io.NopCloser(bytes.NewReader(reqBody))

	router.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)
	var resp handlers.UuidResponse
	json.Unmarshal(responseData, &resp)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, id.String(), resp.CompanyID)
}

func createJWT(username string, jwtSecret string) (string, error) {
	claims := jwt.MapClaims{
		"name": username,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // Expiration time (24 hours)
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
