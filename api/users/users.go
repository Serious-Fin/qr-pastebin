package users

import (
	"context"
	"math/rand"
	"qr-pastebin-api/common"
	"time"

	"github.com/jackc/pgx/v5"
)

type UserCredentials struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SessionData struct {
	SessionId string `json:"sessionId"`
}

type UserDBHandler struct {
	DB *pgx.Conn
}

func NewUserHandler(db *pgx.Conn) *UserDBHandler {
	return &UserDBHandler{DB: db}
}

func (handler *UserDBHandler) CreateUser(request UserCredentials) error {
	_, err := common.GetUserByName(handler.DB, request.Name)
	if err == nil {
		return &UserAlreadyExistsError{}
	}

	hashedPassword, err := common.CreatePasswordHash(request.Password)
	if err != nil {
		return err
	}
	query := "INSERT INTO users (id, name, passwordHash, role) VALUES ($1, $2, $3, $4);"

	_, err = handler.DB.Exec(context.Background(), query, rand.Intn(10000), request.Name, hashedPassword, 0)
	if err != nil {
		return err
	}
	return nil
}

func (handler *UserDBHandler) CreateSession(request UserCredentials) (*SessionData, error) {
	user, err := common.GetUserByName(handler.DB, request.Name)
	if err != nil {
		return nil, &common.PasswordIncorrectError{}
	}

	passwordOk := common.IsPasswordCorrect(user.PasswordHash, request.Password)
	if !passwordOk {
		return nil, &common.PasswordIncorrectError{}
	}

	// Try get active session for this user
	sessionId, err := handler.getActiveSession(user.Id)
	if err == nil {
		return &SessionData{SessionId: sessionId}, nil
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
	return &SessionData{SessionId: sessionId}, nil
}

func (handler *UserDBHandler) GetUserFromSession(sessionId string) (*common.User, error) {
	var user common.User
	err := handler.DB.QueryRow(context.Background(), "SELECT u.id, u.name, u.passwordHash, u.role FROM users AS u RIGHT JOIN sessions AS s ON u.id = s.user_id WHERE expire_at > $1 AND s.session_id = $2;", time.Now(), sessionId).Scan(&user.Id, &user.Name, &user.PasswordHash, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (handler *UserDBHandler) getActiveSession(userId int) (string, error) {
	var session struct {
		UserId    int
		SessionId string
	}

	err := handler.DB.QueryRow(context.Background(), "SELECT session_id FROM sessions WHERE user_id = $1 AND expire_at > $2;", userId, time.Now()).Scan(&session.SessionId)
	if err != nil {
		return "", err
	}
	return session.SessionId, nil
}

func (handler *UserDBHandler) deleteSessions(userId int) error {
	_, err := handler.DB.Exec(context.Background(), "DELETE FROM sessions WHERE user_id = $1;", userId)
	if err != nil {
		return err
	}
	return nil
}

func (handler *UserDBHandler) createNewSession(userId int) (string, error) {
	query := "INSERT INTO sessions (session_id, user_id, expire_at) VALUES ($1, $2, $3);"
	sessionId := common.CreateRandomId(10)

	_, err := handler.DB.Exec(context.Background(), query, sessionId, userId, time.Now().AddDate(0, 0, 7))
	if err != nil {
		return "", err
	}
	return sessionId, nil
}
