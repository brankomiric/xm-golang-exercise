package database

import (
	"context"
	"xm-auth/internal/structs"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AbstractDB interface {
	TestConn(ctx context.Context) error
	CreateUser(ctx context.Context, user structs.User) (*mongo.InsertOneResult, error)
	ValidateUsername(ctx context.Context, username string) (bool, error)
	GetUser(ctx context.Context, username string) (bson.M, error)
}
