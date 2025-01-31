package database

import (
	"fmt"
	"strconv"
	"strings"
	"xm-company/internal/dto"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	DB *sqlx.DB
}

func Initialize(connectionStr string) (*Database, error) {
	db, err := connect(connectionStr)
	if err != nil {
		return nil, err
	}
	return &Database{db}, nil
}

func connect(connectionStr string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", connectionStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateConnectionString(host string, port string, user string, password string, dbname string) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}

func (dbObj *Database) TestConn() error {
	err := dbObj.DB.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (dbObj *Database) AddCompany(input dto.CreateCompany) (*uuid.UUID, error) {
	var ID *uuid.UUID
	err := dbObj.DB.QueryRow("INSERT INTO company (name, description, amount_of_employees, registered, type) VALUES ($1, $2, $3, $4, $5) RETURNING id", input.Name, input.Description, input.AmountOfEmployees, input.Registered, input.Type).Scan(&ID)
	if err != nil {
		return nil, err
	}
	return ID, nil
}

func (dbObj *Database) UpdateCompany(id uuid.UUID, companyDto dto.UpdateCompany) error {
	query := "UPDATE company SET "
	args := []interface{}{}
	argIndex := 1

	if companyDto.Name != nil {
		query += "name = $" + strconv.Itoa(argIndex) + ", "
		args = append(args, *companyDto.Name)
		argIndex++
	}
	if companyDto.Description != nil {
		query += "description = $" + strconv.Itoa(argIndex) + ", "
		args = append(args, *companyDto.Description)
		argIndex++
	}
	if companyDto.AmountOfEmployees != nil {
		query += "amount_of_employees = $" + strconv.Itoa(argIndex) + ", "
		args = append(args, *companyDto.AmountOfEmployees)
		argIndex++
	}
	if companyDto.Registered != nil {
		query += "registered = $" + strconv.Itoa(argIndex) + ", "
		args = append(args, *companyDto.Registered)
		argIndex++
	}
	if companyDto.Type != nil {
		query += "type = $" + strconv.Itoa(argIndex) + ", "
		args = append(args, *companyDto.Type)
		argIndex++
	}

	// Remove trailing comma and space
	query = strings.TrimSuffix(query, ", ")
	query += " WHERE id = $" + strconv.Itoa(argIndex)
	args = append(args, id)

	_, err := dbObj.DB.Exec(query, args...)
	return err
}

func (dbObj *Database) DeleteCompany(id uuid.UUID) error {
	query := "DELETE FROM company WHERE id = $1"
	_, err := dbObj.DB.Exec(query, id)
	return err
}

func (dbObj *Database) GetCompanyByID(id uuid.UUID) (*Company, error) {
	query := "SELECT id, name, description, amount_of_employees, registered, type FROM company WHERE id = $1"
	row := dbObj.DB.QueryRow(query, id)

	var company Company
	err := row.Scan(&company.ID, &company.Name, &company.Description, &company.AmountOfEmployees, &company.Registered, &company.Type)
	if err != nil {
		return nil, err
	}
	return &company, nil
}
