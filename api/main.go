package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"qr-pastebin-api/common"
	"qr-pastebin-api/shares"
	"qr-pastebin-api/users"

	"github.com/jackc/pgx/v5"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

type APIError struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func sendError(c *gin.Context, statusCode int, message string, err error) {
	if statusCode >= 500 {
		sendToDiscord(err.Error())
	}

	apiError := APIError{Message: message}
	if gin.IsDebugging() {
		apiError.Details = err.Error()
	}
	c.IndentedJSON(statusCode, apiError)
}

func sendToDiscord(message string) {
	runes := []rune(message)
	sentChars := 0
	for sentChars < len(runes) {
		end := min(sentChars+2000, len(runes))
		sendDiscordMessage(string(runes[sentChars:end]))
		sentChars = end
	}
}

func sendDiscordMessage(message string) {
	discordToken := os.Getenv("DISCORD_TOKEN")
	channelId := os.Getenv("DISCORD_CHANNEL_ID")
	discordApiUrl := fmt.Sprintf("https://discord.com/api/channels/%s/messages", channelId)
	body := map[string]string{"content": fmt.Sprintf("backend-msg: %s", message)}
	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		// printing this to standard output because this function is supposed to send errors to discord normally
		fmt.Printf("error marshaling discord message body: %v", err)
	}
	req, err := http.NewRequest("POST", discordApiUrl, bytes.NewBuffer(bodyAsBytes))
	if err != nil {
		fmt.Printf("could not create new discord request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", discordToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error while sending request to discord: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("got response status code %d, when expected %d", resp.StatusCode, http.StatusOK)
	}
}

var wrongPasswordErr *common.PasswordIncorrectError
var userAlreadyExistsErr *users.UserAlreadyExistsError
var expiredShareError *shares.ExpiredShareError
var notFoundError *common.NotFoundError

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			statusCode := http.StatusInternalServerError
			message := "An unexpected server error encountered"

			if errors.As(err, &wrongPasswordErr) {
				statusCode = http.StatusUnauthorized
				message = wrongPasswordErr.Error()
			}

			if errors.As(err, &userAlreadyExistsErr) {
				statusCode = http.StatusConflict
				message = "User already exists"
			}

			if errors.As(err, &expiredShareError) {
				statusCode = http.StatusNotFound
				message = expiredShareError.Error()
			}

			if errors.As(err, &notFoundError) {
				statusCode = http.StatusNotFound
				message = notFoundError.Error()
			}

			if errors.Is(err, sql.ErrNoRows) {
				statusCode = http.StatusNotFound
				message = "Resource not found"
			}

			sendError(c, statusCode, message, err)
		}
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}

		sessionId := parts[1]

		user, err := userHandler.GetUserFromSession(sessionId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid session token. %s", err)})
			return
		}

		c.Set("userId", user.Id)
		c.Set("userRole", user.Role)

		c.Next()
	}
}

var shareHandler shares.ShareDBHandler
var userHandler users.UserDBHandler

func main() {
	fmt.Println(os.Getenv("DATABASE_URL"))
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	shareHandler = *shares.NewShareHandler(conn)
	userHandler = *users.NewUserHandler(conn)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	router.Use(ErrorHandlerMiddleware())

	api := router.Group("/")
	api.Use(AuthMiddleware())
	{
		api.GET("/shares", GetShares)
		api.DELETE("/share/:id", DeleteShare)
		api.GET("/share/:id/edit", GetShareForEdit)
		api.PATCH("/share/:id/edit", UpdateShare)
	}

	router.POST("/share", CreateShare)
	router.GET("/share/:id", GetShare)
	router.POST("/share/:id/protected", GetProtectedShare)
	router.GET("/share/:id/protected", IsPasswordProtected)
	router.POST("/user", CreateUser)
	router.GET("/user/session/:sessionId", GetUser)
	router.POST("/user/session", CreateSession)
	router.Run("0.0.0.0:8080")
}

func CreateShare(c *gin.Context) {
	var body shares.ShareRequest
	if err := c.ShouldBind(&body); err != nil {
		c.Error(err)
		return
	}

	response, err := shareHandler.CreateShare(body)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, response)
}

func GetShare(c *gin.Context) {
	shareId := c.Param("id")
	response, err := shareHandler.GetShareForPublic(shareId)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, response)
}

func GetShareForEdit(c *gin.Context) {
	shareId := c.Param("id")
	userId, err := getUserIdFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	response, err := shareHandler.GetShareForOwner(shareId, userId)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, response)
}

func UpdateShare(c *gin.Context) {
	shareId := c.Param("id")
	userId, err := getUserIdFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	var body shares.ShareRequest
	if err := c.ShouldBind(&body); err != nil {
		c.Error(err)
		return
	}

	err = shareHandler.UpdateShare(shareId, userId, body)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, nil)
}

func DeleteShare(c *gin.Context) {
	shareId := c.Param("id")
	userId, err := getUserIdFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}
	userRole, err := getUserRoleFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	permit, err := shareHandler.HasAccessToShare(userId, shareId, userRole)
	if err != nil {
		c.Error(err)
		return
	}

	if !permit {
		c.IndentedJSON(http.StatusUnauthorized, nil)
	}

	err = shareHandler.DeleteShare(shareId)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, nil)
}

func GetShares(c *gin.Context) {
	userId, err := getUserIdFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	response, err := shareHandler.GetShares(userId)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, response)
}

func getUserIdFromContext(c *gin.Context) (int, error) {
	userIdInterface, exists := c.Get("userId")
	if !exists {
		return -1, errors.New("userId not set in context")
	}

	userId, ok := userIdInterface.(int)
	if !ok {
		return -1, errors.New("userId in context is not an int")
	}
	return userId, nil
}

func getUserRoleFromContext(c *gin.Context) (common.Role, error) {
	userIdInterface, exists := c.Get("userRole")
	if !exists {
		return -1, errors.New("userRole not set in context")
	}

	userRole, ok := userIdInterface.(common.Role)
	if !ok {
		return -1, errors.New("userRole in context is not an int")
	}
	return userRole, nil
}

func IsPasswordProtected(c *gin.Context) {
	shareId := c.Param("id")
	response, err := shareHandler.GetShareForPublic(shareId)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, response)
}

func GetProtectedShare(c *gin.Context) {
	shareId := c.Param("id")
	var body shares.GetProtectedShareRequest
	if err := c.ShouldBind(&body); err != nil {
		c.Error(err)
		return
	}

	response, err := shareHandler.GetProtectedShare(shareId, body.Password)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, response)
}

func CreateUser(c *gin.Context) {
	var body users.UserCredentials
	if err := c.ShouldBind(&body); err != nil {
		c.Error(err)
		return
	}

	err := userHandler.CreateUser(body)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, nil)
}

func GetUser(c *gin.Context) {
	sessionId := c.Param("sessionId")
	user, err := userHandler.GetUserFromSession(sessionId)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}

func CreateSession(c *gin.Context) {
	var body users.UserCredentials
	if err := c.ShouldBind(&body); err != nil {
		c.Error(err)
		return
	}

	response, err := userHandler.CreateSession(body)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, response)
}
