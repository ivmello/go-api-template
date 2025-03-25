package auth

import (
	"time"

	"github.com/ivmello/go-api-template/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewUser creates a new user with a hashed password
func NewUser(email, password, name string) (*User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		Name:         name,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// ComparePassword checks if the provided password matches the stored hash
func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// GenerateToken creates a new JWT token for the user
func (u *User) GenerateToken(expirationHours int) (string, error) {
	return auth.GenerateToken(u.ID, u.Email, expirationHours)
}