package users

import (
	"context"
	"fmt"
	"math/rand"
	"qr-pastebin-api/common"
	"time"

	"github.com/jackc/pgx/v5"
)

type UserData struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type CreateSessionResponse struct {
	SessionId string `json:"sessionId"`
}

type Session struct {
	UserId    int
	SessionId string
}

type UserDBHandler struct {
	DB *pgx.Conn
}

func NewUserHandler(db *pgx.Conn) *UserDBHandler {
	return &UserDBHandler{DB: db}
}

func (handler *UserDBHandler) CreateUser(request UserData) error {
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

func (handler *UserDBHandler) CreateSession(request UserData) (*CreateSessionResponse, error) {
	user, err := common.GetUserByName(handler.DB, request.Name)
	if err != nil {
		return nil, &WrongPasswordError{}
	}

	passwordOk := common.IsPasswordCorrect(user.PasswordHash, request.Password)
	if !passwordOk {
		return nil, &WrongPasswordError{}
	}

	// Try get active session for this user
	sessionId, err := handler.getActiveSession(user.Id)
	if err == nil {
		return &CreateSessionResponse{SessionId: sessionId}, nil
	}

	// If no active session, then clean all expired sessions
	err = handler.deleteSessions(user.Id)
	if err != nil {
		return nil, err
	}

	// Create a new session for this user
	sessionId, err = handler.createNewSession(user.Id)
	if err != nil {
		return nil, err
	}
	return &CreateSessionResponse{SessionId: sessionId}, nil
}

func (handler *UserDBHandler) getActiveSession(userId int) (string, error) {
	var session Session
	err := handler.DB.QueryRow(context.Background(), "SELECT sessionId FROM sessions WHERE userId = $1 AND expireAt > $1;", userId, time.Now()).Scan(&session.SessionId)
	if err != nil {
		return "", fmt.Errorf("error getting session with user id '%d': %w", userId, err)
	}
	return session.SessionId, nil
}

func (handler *UserDBHandler) deleteSessions(userId int) error {
	_, err := handler.DB.Exec(context.Background(), "DELETE FROM sessions WHERE userId = $1;", userId)
	if err != nil {
		return fmt.Errorf("error deleting sessions for user with id '%d': %w", userId, err)
	}
	return nil
}

func (handler *UserDBHandler) createNewSession(userId int) (string, error) {
	query := "INSERT INTO sessions (sessionId, userId, expireAt) VALUES ($1, $2, $3);"
	sessionId := common.CreateRandomId(10)

	_, err := handler.DB.Exec(context.Background(), query, sessionId, userId, time.Now().AddDate(0, 0, 7))
	if err != nil {
		return "", fmt.Errorf("couldn't create new session: %w", err)
	}
	return sessionId, nil
}
