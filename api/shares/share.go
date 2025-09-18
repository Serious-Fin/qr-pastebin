package shares

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Share struct {
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

func (handler *ShareDBHandler) CreateShare(request Share) (*CreateShareResponse, error) {
	id := uuid.NewString()
	_, err := handler.DB.Exec(context.Background(), "INSERT INTO shares (id, title, content, password, expireat) VALUES ($1, $2, $3, $4, $5);", id, request.Title, request.Content, request.Password, request.ExpireIn)
	if err != nil {
		return nil, fmt.Errorf("couldn't create new share: %w", err)
	}
	return &CreateShareResponse{ ShareId: id }, nil
}

func (handler *ShareDBHandler) GetShare(id string) (*Share, error) {
	var share Share
	err := handler.DB.QueryRow(context.Background(), "SELECT content, title, expire_at, password FROM shares WHERE id = $1;", id).Scan(&share.Content, &share.Title, &share.ExpireIn, &share.Password)
	if err != nil {
		return nil, fmt.Errorf("couldn't find share with id '%s': %w", id, err)
	}
	return &share, nil
}