package shares

import (
	"context"
	"fmt"
	"qr-pastebin-api/common"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type ShareRequest struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	SetPassword bool   `json:"setPassword"`
	Password    string `json:"password"`
	ExpireIn    string `json:"expireIn"`
	HideAuthor  bool   `json:"hideAuthor"`
	AuthorId    int    `json:"authorId"`
}

type ShareResponse struct {
	Id                  string `json:"id"`
	Title               string `json:"title"`
	Content             string `json:"content"`
	IsPasswordProtected bool   `json:"isPasswordProtected"`
	ExpiresIn           string `json:"expiresIn"`
	AuthorName          string `json:"authorName"`
	HideAuthor          bool   `json:"hideAuthor"`
}

type Share struct {
	Id           string
	Title        string
	Content      string
	PasswordHash string
	ExpireAt     time.Time
	AuthorId     int
	HideAuthor   bool
}

type IsPasswordProtectedResponse struct {
	IsPasswordProtected bool `json:"isPasswordProtected"`
}

type CreateShareResponse struct {
	ShareId string `json:"id"`
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

func (handler *ShareDBHandler) CreateShare(shareBody ShareRequest) (*CreateShareResponse, error) {
	colNames := []string{}
	args := []any{}
	values := []string{}
	argPos := 1

	shareId := common.CreateRandomId(7)
	colNames = append(colNames, "id")
	values = append(values, fmt.Sprintf("$%d", argPos))
	args = append(args, shareId)
	argPos++

	colNames = append(colNames, "title")
	values = append(values, fmt.Sprintf("$%d", argPos))
	args = append(args, shareBody.Title)
	argPos++

	colNames = append(colNames, "content")
	values = append(values, fmt.Sprintf("$%d", argPos))
	args = append(args, shareBody.Content)
	argPos++

	if shareBody.SetPassword {
		passwordHash, err := common.CreatePasswordHash(shareBody.Password)
		if err != nil {
			return nil, err
		}
		colNames = append(colNames, "password")
		values = append(values, fmt.Sprintf("$%d", argPos))
		args = append(args, passwordHash)
		argPos++
	} else {
		colNames = append(colNames, "password")
		values = append(values, fmt.Sprintf("$%d", argPos))
		args = append(args, "")
		argPos++
	}

	expirationDate, err := createExpirationDate(shareBody.ExpireIn)
	if err != nil {
		return nil, err
	}
	colNames = append(colNames, "expire_at")
	values = append(values, fmt.Sprintf("$%d", argPos))
	args = append(args, expirationDate)
	argPos++

	colNames = append(colNames, "hide_author")
	values = append(values, fmt.Sprintf("$%d", argPos))
	args = append(args, shareBody.HideAuthor)
	argPos++

	colNames = append(colNames, "author_id")
	values = append(values, fmt.Sprintf("$%d", argPos))
	args = append(args, shareBody.AuthorId)
	argPos++

	query := fmt.Sprintf("INSERT INTO shares (%s) VALUES (%s);", strings.Join(colNames, ", "), strings.Join(values, ", "))

	_, err = handler.DB.Exec(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("couldn't create new share: %w", err)
	}
	return &CreateShareResponse{ShareId: shareId}, nil
}

func (handler *ShareDBHandler) UpdateShare(shareId string, userId int, shareBody ShareRequest) error {
	setParts := []string{}
	args := []any{}
	argCount := 1

	setParts = append(setParts, fmt.Sprintf("%s = $%d", "title", argCount))
	args = append(args, shareBody.Title)
	argCount++

	setParts = append(setParts, fmt.Sprintf("%s = $%d", "content", argCount))
	args = append(args, shareBody.Content)
	argCount++

	if shareBody.SetPassword {
		if shareBody.Password == "" {
			setParts = append(setParts, fmt.Sprintf("%s = $%d", "password", argCount))
			args = append(args, "")
			argCount++
		} else {
			passwordHash, err := common.CreatePasswordHash(shareBody.Password)
			if err != nil {
				return err
			}
			setParts = append(setParts, fmt.Sprintf("%s = $%d", "password", argCount))
			args = append(args, passwordHash)
			argCount++
		}
	}

	if shareBody.ExpireIn != "no-change" {
		expirationDate, err := createExpirationDate(shareBody.ExpireIn)
		if err != nil {
			return err
		}
		setParts = append(setParts, fmt.Sprintf("%s = $%d", "expire_at", argCount))
		args = append(args, expirationDate)
		argCount++
	}
	setParts = append(setParts, fmt.Sprintf("%s = $%d", "hide_author", argCount))
	args = append(args, shareBody.HideAuthor)
	argCount++

	setQueryPart := strings.Join(setParts, ", ")

	shareIdIndex := fmt.Sprintf("$%d", argCount)
	args = append(args, shareId)
	argCount++

	authorIdIndex := fmt.Sprintf("$%d", argCount)
	args = append(args, userId)
	argCount++

	query := fmt.Sprintf("UPDATE shares SET %s WHERE id = %s AND author_id = %s;", setQueryPart, shareIdIndex, authorIdIndex)
	_, err := handler.DB.Exec(context.Background(), query, args...)
	return err
}

func (handler *ShareDBHandler) GetShareForPublic(id string) (*ShareResponse, error) {
	share, err := handler.readShare(id)
	if err != nil {
		return nil, err
	}

	if !share.ExpireAt.IsZero() && time.Now().After(share.ExpireAt) {
		return nil, &ExpiredShareError{}
	}

	shareResponse, err := handler.transformToShareResponse(share)
	if err != nil {
		return nil, err
	}

	if shareResponse.HideAuthor {
		shareResponse.AuthorName = ""
	}

	return shareResponse, nil
}

func (handler *ShareDBHandler) GetShareForOwner(shareId string, userId int) (*ShareResponse, error) {
	permit, err := handler.HasAccessToShare(userId, shareId, 0)
	if err != nil {
		return nil, err
	}
	if !permit {
		return nil, &common.NotFoundError{}
	}

	share, err := handler.readShare(shareId)
	if err != nil {
		return nil, err
	}

	shareResponse, err := handler.transformToShareResponse(share)
	if err != nil {
		return nil, err
	}

	return shareResponse, nil
}

func (handler *ShareDBHandler) GetShares(userId int) ([]ShareResponse, error) {
	shares, err := handler.readShares(userId)
	if err != nil {
		return nil, err
	}

	shareResponses := make([]ShareResponse, 0)
	for _, share := range shares {
		newShareResponse, err := handler.transformToShareResponse(&share)
		if err != nil {
			return nil, err
		}
		shareResponses = append(shareResponses, *newShareResponse)
	}

	return shareResponses, nil
}

func (handler *ShareDBHandler) GetProtectedShare(id string, password string) (*ShareResponse, error) {
	share, err := handler.readShare(id)
	if err != nil {
		return nil, err
	}

	if !share.ExpireAt.IsZero() && time.Now().After(share.ExpireAt) {
		return nil, &ExpiredShareError{}
	}

	passwordOk := common.IsPasswordCorrect(share.PasswordHash, password)
	if !passwordOk {
		return nil, &common.PasswordIncorrectError{}
	}

	shareResponse, err := handler.transformToShareResponse(share)
	if err != nil {
		return nil, err
	}

	if shareResponse.HideAuthor {
		shareResponse.AuthorName = ""
	}
	return shareResponse, nil
}

func (handler *ShareDBHandler) DeleteShare(shareId string) error {
	query := "DELETE FROM shares WHERE id = $1;"
	_, err := handler.DB.Exec(context.Background(), query, shareId)
	if err != nil {
		return err
	}
	return nil
}

func (handler *ShareDBHandler) HasAccessToShare(userId int, shareId string, role common.Role) (bool, error) {
	if role.String() == "admin" {
		return true, nil
	}

	var count int
	query := "SELECT COUNT(*) FROM shares WHERE author_id = $1 AND id = $2;"
	err := handler.DB.QueryRow(context.Background(), query, userId, shareId).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func (handler *ShareDBHandler) IsPasswordProtected(id string) (*IsPasswordProtectedResponse, error) {
	share, err := handler.readShare(id)
	if err != nil {
		return nil, err
	}
	if share.PasswordHash == "" {
		return &IsPasswordProtectedResponse{IsPasswordProtected: false}, nil
	} else {
		return &IsPasswordProtectedResponse{IsPasswordProtected: true}, nil
	}
}

func (handler *ShareDBHandler) readShare(shareId string) (*Share, error) {
	var share Share
	err := handler.DB.QueryRow(context.Background(), "SELECT id, title, content, expire_at, password, author_id, hide_author FROM shares WHERE id = $1;", shareId).Scan(&share.Id, &share.Title, &share.Content, &share.ExpireAt, &share.PasswordHash, &share.AuthorId, &share.HideAuthor)
	if err != nil {
		return nil, err
	}
	return &share, nil
}

func (handler *ShareDBHandler) readShares(userId int) ([]Share, error) {
	rows, err := handler.DB.Query(context.Background(), "SELECT s.id, s.title, s.content, s.expire_at, s.password, s.author_id, s.hide_author FROM users AS u RIGHT JOIN shares AS s ON u.id = s.author_id WHERE u.id = $1;", userId)
	if err != nil {
		return nil, fmt.Errorf("error querying shares: %w", err)
	}
	defer rows.Close()

	shares := make([]Share, 0)
	for rows.Next() {
		var share Share
		err := rows.Scan(&share.Id, &share.Title, &share.Content, &share.ExpireAt, &share.PasswordHash, &share.AuthorId, &share.HideAuthor)
		if err != nil {
			return nil, err
		}
		shares = append(shares, share)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return shares, nil
}

func createExpirationDate(expireIn string) (time.Time, error) {
	if expireIn == "never" {
		return time.Time{}, nil
	}

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

func (handler *ShareDBHandler) transformToShareResponse(share *Share) (*ShareResponse, error) {
	var shareResp ShareResponse
	shareResp.Id = share.Id
	shareResp.Content = share.Content
	shareResp.HideAuthor = share.HideAuthor
	shareResp.Title = share.Title

	if share.PasswordHash != "" {
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
