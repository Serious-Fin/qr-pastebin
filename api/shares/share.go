package shares

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Share struct {
	Id 		 string `json:"id,omitempty"`
	Content  string `json:"content"`
	Title    string `json:"title,omitempty"`
	ExpireAt string `json:"expireAt,omitempty"`
	Password string `json:"password,omitempty"`
}

type DBShare struct {
	Content  string
	Title    sql.NullString
	ExpireAt sql.NullString
	Password sql.NullString
}

type CreateShareRequest struct {
	Content  string `json:"content"`
	Title    string `json:"title,omitempty"`
	ExpireIn string `json:"expireIn,omitempty"`
	Password string `json:"password,omitempty"`
}

type CreateShareResponse struct {
	ShareId string `json:"id"`
}

type ShareDBHandler struct {
	DB *pgx.Conn
}

func NewShareHandler(db *pgx.Conn) *ShareDBHandler {
	return &ShareDBHandler{DB: db}
}

func (handler *ShareDBHandler) CreateShare(request CreateShareRequest) (*CreateShareResponse, error) {
	id := uuid.NewString()
	_, err := handler.DB.Exec(context.Background(), "INSERT INTO shares (id, title, content, password, expireat) VALUES ($1, $2, $3, $4, $5);", id, request.Title, request.Content, request.Password, request.ExpireIn)
	if err != nil {
		return nil, fmt.Errorf("couldn't create new share: %w", err)
	}
	return &CreateShareResponse{ ShareId: id }, nil
}

func (handler *ShareDBHandler) GetShare(id string) (*Share, error) {
	var dbShare DBShare
	err := handler.DB.QueryRow(context.Background(), "SELECT title, content, expire_at, password FROM shares WHERE id = $1;", id).Scan(&dbShare.Title, &dbShare.Content, &dbShare.ExpireAt, &dbShare.Password)
	if err != nil {
		return nil, fmt.Errorf("couldn't find share with id '%s': %w", id, err)
	}

	share := convertShareFromDB(dbShare)
	return &share, nil
}

func createNewShare(request CreateShareRequest) Share {
	var share Share
	share.Id = uuid.NewString()
	share.Content = request.Content
	share.Title = request.Title
	// TODO: hash the password
	share.Password = share.Password
	// TODO: calculate NOW + expireIn
	share.ExpireAt = request.ExpireIn

	return share
}

func convertShareFromDB(dbShare DBShare) Share {
	var share Share
	share.Content = dbShare.Content
	if dbShare.Title.Valid {
		share.Title = dbShare.Title.String
	} else {
		share.Title = ""
	}

	if dbShare.ExpireAt.Valid {
		share.ExpireAt = dbShare.ExpireAt.String
	} else {
		share.ExpireAt = ""
	}

	if dbShare.Password.Valid {
		share.Password = dbShare.Password.String
	} else {
		share.Password = ""
	}
	return share
}