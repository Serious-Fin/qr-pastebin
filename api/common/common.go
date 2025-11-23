package common

import (
	"context"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type Role int

const (
	USER  = 0
	ADMIN = 1
)

var roles = map[Role]string{
	USER:  "user",
	ADMIN: "admin",
}

func (r Role) String() string {
	return roles[r]
}

type User struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	PasswordHash string `json:"password"`
	Role         Role   `json:"role"`
	IsOauth      bool   `json:"isoauth"`
}

type HealthResponse struct {
	Status   string `json:"status"`
	Hostname string `json:"hostname"`
}

func GetHostname() string {
	// Execute the 'hostname' command
	cmd := exec.Command("hostname")

	// Capture the output
	output, err := cmd.Output()
	if err != nil {
		// If the command fails, return a safe default or error message
		return "error_retrieving_hostname"
	}

	// Convert output to string and trim any newline characters
	return strings.TrimSpace(string(output))
}

func CreatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func IsPasswordCorrect(passwordHash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err == nil
}

func GetUserByName(db *pgx.Conn, name string) (*User, error) {
	var user User
	err := db.QueryRow(context.Background(), "SELECT id, name, passwordHash, role, isoauth FROM users WHERE name = $1;", name).Scan(&user.Id, &user.Name, &user.PasswordHash, &user.Role, &user.IsOauth)
	if err != nil {
		return nil, fmt.Errorf("error getting user with name '%s': %w", name, err)
	}
	return &user, nil
}

func GetUserById(db *pgx.Conn, id int) (*User, error) {
	var user User
	err := db.QueryRow(context.Background(), "SELECT id, name, passwordHash, role, isoauth FROM users WHERE id = $1;", id).Scan(&user.Id, &user.Name, &user.PasswordHash, &user.Role, &user.IsOauth)
	if err != nil {
		return nil, fmt.Errorf("error getting user with id '%d': %w", id, err)
	}
	return &user, nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789_")

func CreateRandomId(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
