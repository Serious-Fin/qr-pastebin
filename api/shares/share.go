package shares

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"qr-pastebin-api/common"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type Share struct {
	Id         string    `json:"id"`
	Content    string    `json:"content"`
	Title      string    `json:"title,omitempty"`
	ExpireAt   time.Time `json:"expireAt,omitempty"`
	Password   string    `json:"password,omitempty"`
	AuthorId   int       `json:"authorId,omitempty"`
	HideAuthor bool      `json:"hideAuthor"`
}

type GetShareResponse struct {
	Id                  string `json:"id"`
	Content             string `json:"content"`
	Title               string `json:"title,omitempty"`
	ExpiresIn           string `json:"expiresIn,omitempty"`
	IsPasswordProtected bool   `json:"isPasswordProtected,omitempty"`
	AuthorName          string `json:"authorName,omitempty"`
	HideAuthor          bool   `json:"hideAuthor"`
}

type IsPasswordProtectedResponse struct {
	IsPasswordProtected bool `json:"isPasswordProtected"`
}

type DBShare struct {
	Id         string
	Content    string
	Title      sql.NullString
	ExpireAt   sql.NullTime
	Password   sql.NullString
	AuthorId   sql.NullInt32
	HideAuthor bool
}

type CreateShareRequest struct {
	Content    string `json:"content"`
	Title      string `json:"title,omitempty"`
	ExpireIn   string `json:"expireIn,omitempty"`
	Password   string `json:"password,omitempty"`
	AuthorId   int32  `json:"authorId"`
	HideAuthor bool   `json:"hideAuthor"`
}

type CreateShareResponse struct {
	ShareId string `json:"id"`
}

type EditShareRequest struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	SetPassword bool   `json:"setPassword"`
	Password    string `json:"password"`
	ExpireIn    string `json:"expireIn"`
	HideAuthor  bool   `json:"hideAuthor"`
}

type GetProtectedShareRequest struct {
	Password string `json:"password"`
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
	share, err := handler.readShareFromDB(id)
	if err != nil {
		return nil, err
	}

	if !share.ExpireAt.IsZero() && time.Now().After(share.ExpireAt) {
		return nil, errors.New("trying to access expired share")
	}

	shareResponse, err := handler.createGetShareResponse(share)
	if err != nil {
		return nil, err
	}

	if shareResponse.HideAuthor {
		shareResponse.AuthorName = ""
	}

	return shareResponse, nil
}

func (handler *ShareDBHandler) GetShareForEdit(shareId string, userId int) (*GetShareResponse, error) {
	share, err := handler.readShareFromDBFromUser(shareId, userId)
	if err != nil {
		return nil, err
	}

	shareResponse, err := handler.createGetShareResponse(share)
	if err != nil {
		return nil, err
	}

	return shareResponse, nil
}

func (handler *ShareDBHandler) UpdateShare(shareId string, userId int, shareBody EditShareRequest) error {
	colNames := make([]string, 0)
	values := make([]string, 0)
	args := make([]interface{}, 0)
	argCount := 1

	colNames = append(colNames, "title")
	values = append(values, fmt.Sprintf("$%d", argCount))
	args = append(args, shareBody.Title)
	argCount++

	colNames = append(colNames, "content")
	values = append(values, fmt.Sprintf("$%d", argCount))
	args = append(args, shareBody.Content)
	argCount++

	if shareBody.SetPassword {
		passwordHash, err := common.CreatePasswordHash(shareBody.Password)
		if err != nil {
			return err
		}
		colNames = append(colNames, "password")
		values = append(values, fmt.Sprintf("$%d", argCount))
		args = append(args, passwordHash)
		argCount++
	}

	if shareBody.ExpireIn != "no-change" {
		expirationDate, err := createExpirationDate(shareBody.ExpireIn)
		if err != nil {
			return err
		}
		colNames = append(colNames, "expire_at")
		values = append(values, fmt.Sprintf("$%d", argCount))
		args = append(args, expirationDate)
		argCount++
	}

	colNames = append(colNames, "hide_author")
	values = append(values, fmt.Sprintf("$%d", argCount))
	args = append(args, shareBody.HideAuthor)
	argCount++

	colNamesString := strings.Join(colNames, ", ")
	valueIndexString := strings.Join(values, ", ")

	shareIdIndex := fmt.Sprintf("$%d", argCount)
	args = append(args, shareId)
	argCount++

	authorIdIndex := fmt.Sprintf("$%d", argCount)
	args = append(args, userId)
	argCount++

	query := fmt.Sprintf("UPDATE shares SET (%s) VALUES (%s) WHERE id = %s AND author_id = %s;", colNamesString, valueIndexString, shareIdIndex, authorIdIndex)
	_, err := handler.DB.Exec(context.Background(), query, args...)
	return err
}

func (handler *ShareDBHandler) GetShares(userId int) ([]GetShareResponse, error) {
	shares, err := handler.readSharesFromDB(userId)
	if err != nil {
		return nil, err
	}

	shareResponses := make([]GetShareResponse, 0)
	for _, share := range shares {
		newShareResponse, err := handler.createGetShareResponse(share)
		if err != nil {
			return nil, err
		}
		shareResponses = append(shareResponses, *newShareResponse)
	}

	return shareResponses, nil
}

func (handler *ShareDBHandler) GetProtectedShare(id string, password string) (*GetShareResponse, error) {
	share, err := handler.readShareFromDB(id)
	if err != nil {
		return nil, err
	}

	if !share.ExpireAt.IsZero() && time.Now().After(share.ExpireAt) {
		return nil, errors.New("trying to access expired share")
	}

	passwordOk := common.IsPasswordCorrect(share.Password, password)
	if !passwordOk {
		return nil, &PasswordIncorrectError{}
	}

	shareResponse, err := handler.createGetShareResponse(share)
	if err != nil {
		return nil, err
	}
	return shareResponse, nil
}

func (handler *ShareDBHandler) DeleteShare(shareId string, userId int) error {
	query := "DELETE FROM shares WHERE id = $1 AND author_id = $2;"
	_, err := handler.DB.Exec(context.Background(), query, shareId, userId)
	if err != nil {
		return fmt.Errorf("couldn't delete share with id '%s' for user with id '%d': %w", shareId, userId, err)
	}
	return nil
}

func (handler *ShareDBHandler) IsPasswordProtected(id string) (*IsPasswordProtectedResponse, error) {
	share, err := handler.readShareFromDB(id)
	if err != nil {
		return nil, err
	}
	if share.Password == "" {
		return &IsPasswordProtectedResponse{IsPasswordProtected: false}, nil
	} else {
		return &IsPasswordProtectedResponse{IsPasswordProtected: true}, nil
	}
}

func (handler *ShareDBHandler) readShareFromDB(shareId string) (Share, error) {
	var dbShare DBShare
	err := handler.DB.QueryRow(context.Background(), "SELECT id, title, content, expire_at, password, author_id, hide_author FROM shares WHERE id = $1;", shareId).Scan(&dbShare.Id, &dbShare.Title, &dbShare.Content, &dbShare.ExpireAt, &dbShare.Password, &dbShare.AuthorId, &dbShare.HideAuthor)
	if err != nil {
		return Share{}, fmt.Errorf("error getting share with id '%s': %w", shareId, err)
	}
	return convertShareFromDB(dbShare), nil
}

func (handler *ShareDBHandler) readShareFromDBFromUser(shareId string, userId int) (Share, error) {
	var dbShare DBShare
	err := handler.DB.QueryRow(context.Background(), "SELECT id, title, content, expire_at, password, author_id, hide_author FROM shares WHERE id = $1 AND author_id = $2;", shareId, userId).Scan(&dbShare.Id, &dbShare.Title, &dbShare.Content, &dbShare.ExpireAt, &dbShare.Password, &dbShare.AuthorId, &dbShare.HideAuthor)
	if err != nil {
		return Share{}, fmt.Errorf("error getting share with id '%s': %w", shareId, err)
	}
	return convertShareFromDB(dbShare), nil
}

func (handler *ShareDBHandler) readSharesFromDB(userId int) ([]Share, error) {
	rows, err := handler.DB.Query(context.Background(), "SELECT s.id, s.title, s.content, s.expire_at, s.password, s.author_id, s.hide_author FROM users AS u RIGHT JOIN shares AS s ON u.id = s.author_id WHERE u.id = $1;", userId)
	if err != nil {
		return nil, fmt.Errorf("error querying shares: %w", err)
	}
	defer rows.Close()

	shares := make([]Share, 0)
	for rows.Next() {
		var dbShare DBShare
		err := rows.Scan(&dbShare.Id, &dbShare.Title, &dbShare.Content, &dbShare.ExpireAt, &dbShare.Password, &dbShare.AuthorId, &dbShare.HideAuthor)
		if err != nil {
			return nil, fmt.Errorf("error scanning share row: %w", err)
		}
		shares = append(shares, convertShareFromDB(dbShare))
	}
	// check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return shares, nil
}

func createNewShare(request CreateShareRequest) (*Share, error) {
	var share Share
	share.Id = common.CreateRandomId(7)
	share.Content = request.Content
	share.Title = request.Title
	share.AuthorId = int(request.AuthorId)
	share.HideAuthor = request.HideAuthor
	if request.Password != "" {
		passwordHash, err := common.CreatePasswordHash(request.Password)
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

func createInsertStatement(share Share) (string, []any) {
	colNames := []string{}
	args := []any{}
	values := []string{}
	argPos := 1

	colNames, args, values, argPos = tryAddColumnToQuery("id", share.Id, argPos, colNames, args, values)
	colNames, args, values, argPos = tryAddColumnToQuery("content", share.Content, argPos, colNames, args, values)
	colNames, args, values, argPos = tryAddColumnToQuery("title", share.Title, argPos, colNames, args, values)
	colNames, args, values, argPos = tryAddColumnToQuery("expire_at", share.ExpireAt, argPos, colNames, args, values)
	colNames, args, values, argPos = tryAddColumnToQuery("author_id", share.AuthorId, argPos, colNames, args, values)
	colNames, args, values, argPos = tryAddColumnToQuery("hide_author", share.HideAuthor, argPos, colNames, args, values)
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
	share.HideAuthor = dbShare.HideAuthor
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

	if dbShare.AuthorId.Valid {
		share.AuthorId = int(dbShare.AuthorId.Int32)
	} else {
		share.AuthorId = -1
	}
	return share
}

func (handler *ShareDBHandler) createGetShareResponse(share Share) (*GetShareResponse, error) {
	var shareResp GetShareResponse
	shareResp.Id = share.Id
	shareResp.Content = share.Content
	shareResp.HideAuthor = share.HideAuthor

	if share.Title != "" {
		shareResp.Title = share.Title
	}
	if share.Password != "" {
		shareResp.IsPasswordProtected = true
	}
	if !share.ExpireAt.IsZero() {
		shareResp.ExpiresIn = createExpireInTextFromDate(share.ExpireAt)
	}
	if share.AuthorId != -1 {
		author, err := common.GetUserById(handler.DB, share.AuthorId)
		if err != nil {
			return nil, err
		}
		shareResp.AuthorName = author.Name
	}
	return &shareResp, nil
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
