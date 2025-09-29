package users

import (
	"context"
	"fmt"
	"math/rand"
	"qr-pastebin-api/common"

	"github.com/jackc/pgx/v5"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type GetUserRequest struct {
	Name           string `json:"name"`
	HashedPassword string `json:"hashedPassword"`
}

type GetUserResponse struct {
	Name      string `json:"name"`
	SessionId string `json:"sessionId"`
}

type UserDBHandler struct {
	DB *pgx.Conn
}

func NewUserHandler(db *pgx.Conn) *UserDBHandler {
	return &UserDBHandler{DB: db}
}

func (handler *UserDBHandler) CreateUser(request CreateUserRequest) error {
	_, err := common.GetUserByName(handler.DB, request.Name)
	if err == nil {
		return &UserAlreadyExistsError{}
	}

	hashedPassword, err := common.CreatePasswordHash(request.Password)
	if err != nil {
		return err
	}
	query := "INSERT INTO users (id, name, password) VALUES ($1, $2, $3);"

	_, err = handler.DB.Exec(context.Background(), query, rand.Intn(10000), request.Name, hashedPassword)
	if err != nil {
		return fmt.Errorf("couldn't create new user: %w", err)
	}
	return nil
}

func (handler *UserDBHandler) GetUser(request GetUserRequest) (*GetUserResponse, error) {
	// get user by name
	// check if password is correct
	// delete all sessions related to this user
	// create a new session and return it
	return nil, nil
}
