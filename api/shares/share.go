package shares

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type Share struct {
	Id       string    `json:"id"`
	Content  string    `json:"content"`
	Title    string    `json:"title,omitempty"`
	ExpireAt time.Time `json:"expireAt,omitempty"`
	Password string    `json:"password,omitempty"`
}

type GetShareResponse struct {
	Id                  string `json:"id"`
	Content             string `json:"content"`
	Title               string `json:"title,omitempty"`
	ExpiresIn           string `json:"expiresIn,omitempty"`
	IsPasswordProtected bool   `json:"isPasswordProtected,omitempty"`
	Author              string `json:"author,omitempty"`
}

type DBShare struct {
	Id       string
	Content  string
	Title    sql.NullString
	ExpireAt sql.NullTime
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
	share, err := createNewShare(request)
	if err != nil {
		return nil, err
	}

	query, arguments := createInsertStatement(*share)
	_, err = handler.DB.Exec(context.Background(), query, arguments...)
	if err != nil {
		return nil, fmt.Errorf("couldn't create new share: %w", err)
	}
	return &CreateShareResponse{ShareId: share.Id}, nil
}

func (handler *ShareDBHandler) GetShare(id string) (*GetShareResponse, error) {
	var dbShare DBShare
	err := handler.DB.QueryRow(context.Background(), "SELECT id, title, content, expire_at, password FROM shares WHERE id = $1;", id).Scan(&dbShare.Id, &dbShare.Title, &dbShare.Content, &dbShare.ExpireAt, &dbShare.Password)
	if err != nil {
		return nil, fmt.Errorf("couldn't find share with id '%s': %w", id, err)
	}

	share := convertShareFromDB(dbShare)
	shareResponse := createGetShareResponse(share)
	return &shareResponse, nil
}

func createNewShare(request CreateShareRequest) (*Share, error) {
	var share Share
	share.Id = createRandomId(7)
	share.Content = request.Content
	share.Title = request.Title
	if request.Password != "" {
		passwordHash, err := createPasswordHash(request.Password)
		if err != nil {
			return nil, err
		}
		share.Password = passwordHash
	}
	if request.ExpireIn != "" {
		expirationDate, err := createExpirationDate(request.ExpireIn)
		if err != nil {
			return nil, err
		}
		share.ExpireAt = expirationDate
	}

	return &share, nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789_")

func createRandomId(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func createInsertStatement(share Share) (string, []any) {
	colNames := []string{}
	args := []any{}
	values := []string{}
	argPos := 1

	colNames, args, values, argPos = tryAddColumnToQuery("id", share.Id, argPos, colNames, args, values)
	colNames, args, values, argPos = tryAddColumnToQuery("content", share.Content, argPos, colNames, args, values)
	colNames, args, values, argPos = tryAddColumnToQuery("title", share.Title, argPos, colNames, args, values)
	colNames, args, values, argPos = tryAddColumnToQuery("expire_at", share.ExpireAt, argPos, colNames, args, values)
	colNames, args, values, _ = tryAddColumnToQuery("password", share.Password, argPos, colNames, args, values)

	return fmt.Sprintf("INSERT INTO shares (%s) VALUES (%s);", strings.Join(colNames, ", "), strings.Join(values, ", ")), args
}

func tryAddColumnToQuery(columnName string, value any, argPos int, columnNames []string, args []any, values []string) ([]string, []any, []string, int) {
	if value == "" {
		return columnNames, args, values, argPos
	}

	columnNames = append(columnNames, columnName)
	values = append(values, fmt.Sprintf("$%d", argPos))
	args = append(args, value)
	return columnNames, args, values, argPos + 1
}

func createPasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func createExpirationDate(expireIn string) (time.Time, error) {
	parts := strings.Split(expireIn, "_")
	if len(parts) != 2 {
		return time.Now(), fmt.Errorf("expiration period is not of correct format '%s', make sure it is formatted as '5_days'", expireIn)
	}

	durationCount, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Now(), fmt.Errorf("could not parse duration '%s' to int: %w", parts[0], err)
	}

	var duration time.Duration
	switch parts[1] {
	case "minutes":
		duration = time.Minute * time.Duration(durationCount)
	case "hours":
		duration = time.Hour * time.Duration(durationCount)
	case "days":
		return time.Now().AddDate(0, 0, durationCount), nil
	case "weeks":
		return time.Now().AddDate(0, 0, durationCount*7), nil
	case "months":
		return time.Now().AddDate(0, durationCount, 0), nil
	case "years":
		return time.Now().AddDate(durationCount, 0, 0), nil
	default:
		return time.Now(), fmt.Errorf("unknown duration type '%s'", parts[1])
	}

	return time.Now().Add(duration), nil
}

func convertShareFromDB(dbShare DBShare) Share {
	var share Share
	share.Id = dbShare.Id
	share.Content = dbShare.Content
	if dbShare.Title.Valid {
		share.Title = dbShare.Title.String
	} else {
		share.Title = ""
	}

	if dbShare.ExpireAt.Valid {
		share.ExpireAt = dbShare.ExpireAt.Time
	} else {
		share.ExpireAt = time.Time{}
	}

	if dbShare.Password.Valid {
		share.Password = dbShare.Password.String
	} else {
		share.Password = ""
	}
	return share
}

func createGetShareResponse(share Share) GetShareResponse {
	var shareResp GetShareResponse
	shareResp.Id = share.Id
	shareResp.Content = share.Content

	if share.Title != "" {
		shareResp.Title = share.Title
	}
	if share.Password != "" {
		shareResp.IsPasswordProtected = true
	}
	if !share.ExpireAt.IsZero() {
		shareResp.ExpiresIn = createExpireInTextFromDate(share.ExpireAt)
	}
	// TODO: add author to db table as optional FK id
	// TODO: create user table (name, email, hashed password)
	// TODO: add one test user

	return shareResp
}

func createExpireInTextFromDate(expireAt time.Time) string {
	now := time.Now()
	if now.After(expireAt) {
		return "Already expired"
	}

	years := expireAt.Year() - now.Year()
	months := expireAt.Month() - now.Month()
	days := expireAt.Day() - now.Day()

	if years == 0 && months == 0 && days == 0 {
		return "Expires today"
	}
	if days < 0 {
		prevMonthDays := time.Date(expireAt.Year(), expireAt.Month(), 0, 0, 0, 0, 0, expireAt.Location())
		days += prevMonthDays.Day()
		months--
	}
	if months < 0 {
		months += 12
		years--
	}

	dateComponents := make([]string, 0)
	dateComponents = append(dateComponents, "Expires in")
	if years == 1 {
		dateComponents = append(dateComponents, fmt.Sprintf("%d year", years))
	} else if years > 1 {
		dateComponents = append(dateComponents, fmt.Sprintf("%d years", years))
	}
	if months == 1 {
		dateComponents = append(dateComponents, fmt.Sprintf("%d month", months))
	} else if months > 1 {
		dateComponents = append(dateComponents, fmt.Sprintf("%d months", months))
	}
	if days == 1 {
		dateComponents = append(dateComponents, fmt.Sprintf("%d day", days))
	} else if days > 1 {
		dateComponents = append(dateComponents, fmt.Sprintf("%d days", days))
	}

	if len(dateComponents) <= 2 {
		return strings.Join(dateComponents, " ")
	}

	allPartsExceptLast := strings.Join(dateComponents[0:len(dateComponents)-1], " ")
	return fmt.Sprintf("%s and %s", allPartsExceptLast, dateComponents[len(dateComponents)-1])
}
