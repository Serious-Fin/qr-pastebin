package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

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

// TODO: if logged in, can see created shares
// TODO: create share editing page
// TODO: host via docker
// TODO: setup HTTPS

// TODO: add user_id to share creation if such exists during creation process
// TODO: add page which displays all created shares
// TODO: add button to each share to edit it
// TODO: add form page where user can edit share and save it
// TODO: add share save after edit functionality
// TODO: clean up error handling
// TODO: clean up objects, seems like I have hundreds of different interfaces
