package database

import (
	"context"
	"errors"
	"xm-auth/internal/structs"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	DATABASE   = "xm-auth"
	COLLECTION = "users"
)

type Db struct {
	MdbClient *mongo.Client
}

func New(connStr string, ctx context.Context) (*Db, error) {
	db := &Db{}
	client, err := mongo.Connect(options.Client().ApplyURI(connStr))
	if err != nil {
		return nil, err
	}
	db.MdbClient = client
	return db, nil
}

func (db *Db) TestConn(ctx context.Context) error {
	return db.MdbClient.Ping(ctx, nil)
}

func (db *Db) CreateUser(ctx context.Context, user structs.User) (*mongo.InsertOneResult, error) {
	return db.MdbClient.Database(DATABASE).Collection(COLLECTION).InsertOne(ctx, user)
}

func (db *Db) ValidateUsername(ctx context.Context, username string) (bool, error) {
	filter := bson.M{"username": username}
	var result bson.M
	err := db.MdbClient.Database(DATABASE).Collection(COLLECTION).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (db *Db) GetUser(ctx context.Context, username string) (bson.M, error) {
	filter := bson.M{"username": username}
	var result bson.M
	err := db.MdbClient.Database(DATABASE).Collection(COLLECTION).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
