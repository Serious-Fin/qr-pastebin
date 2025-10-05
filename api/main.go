package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

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
	apiError := APIError{Message: message}
	if gin.IsDebugging() {
		apiError.Details = err.Error()
	}
	c.IndentedJSON(statusCode, apiError)
}

var wrongPasswordErr *shares.PasswordIncorrectError
var userAlreadyExistsErr *users.UserAlreadyExistsError
var wrongNameOrPassErr *users.WrongPasswordError

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			statusCode := http.StatusInternalServerError
			message := "An unexpected server error encountered"

			if errors.As(err, &wrongPasswordErr) {
				statusCode = http.StatusUnauthorized
				message = "Wrong password"
			}

			if errors.As(err, &userAlreadyExistsErr) {
				statusCode = http.StatusConflict
				message = "User already exists"
			}

			if errors.As(err, &wrongNameOrPassErr) {
				statusCode = http.StatusUnauthorized
				message = wrongNameOrPassErr.Error()
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

		user, err := userHandler.ValidateSession(sessionId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid session token. %s", err)})
			return
		}

		c.Set("userId", user.Id)

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
	var body shares.CreateShareRequest
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
	response, err := shareHandler.GetShare(shareId)
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

	response, err := shareHandler.GetShareForEdit(shareId, userId)
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

	var body shares.CreateShareRequest
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

	err = shareHandler.DeleteShare(shareId, userId)
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

func IsPasswordProtected(c *gin.Context) {
	shareId := c.Param("id")
	response, err := shareHandler.GetShare(shareId)
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
	var body users.UserData
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
	var body users.UserData
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

// TODO: add API functionality to update share info
// TODO: hook up "save edit" button to submit changes to API
// TODO: add toast message on success or failure

// TODO: add user roles (user/admin)
// TODO: if user is admin, he has a "delete share" button

// TODO: clean up objects in API (seems like I have 5 different User objects, 5 different share objects)
// TODO: clean up error handling in API
// TODO: add logging to discord of errors
// TODO: clean up API code (methods that do similar things, naming)

// TODO: clean up objects in WEB
// TODO: clean up error handling in WEB
// TODO: add logging to discord of errors
// TODO: clean up WEB code

// TODO: host via docker
// TODO: setup HTTPS
