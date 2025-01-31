package handlers

import "go.mongodb.org/mongo-driver/v2/bson"

type HealthcheckResponse struct {
	Services Services `json:"services"`
}

type Services struct {
	Mongo HealthcheckMessage `json:"Mongo"`
}

type HealthcheckMessage struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type RegisterSuccessResponse struct {
	CreatedUserId bson.ObjectID `json:"created_user_id"`
}

type JWTResponse struct {
	Token string `json:"token"`
}
