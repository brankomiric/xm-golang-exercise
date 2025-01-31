package database

import (
	"fmt"
	"os"
)

type DBConnParams struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func ReadConnectionStringParams() (*DBConnParams, error) {
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")
	dbname := os.Getenv("PG_DB")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		return nil, fmt.Errorf("missing DB connection parameters")

	}
	return &DBConnParams{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbname,
	}, nil
}
