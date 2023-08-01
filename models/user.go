package models

import (
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uint
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) Create(email, password string) (*User, error) {
	email = strings.ToLower(email)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Create user: %w", err)
	}
	passwordHash := string(hashedBytes)
	row := us.DB.QueryRow(`
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2) RETURNING id;
	`, email, passwordHash)

	user := User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	err = row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("error when scan user id: %w", err)
	}

	return &user, nil
}
