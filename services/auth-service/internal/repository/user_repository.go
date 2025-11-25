package repository

import (
	"database/sql"
	"fmt"

	"github.com/shinshark/azure-url-shortener/services/auth-service/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO Users (Username, PasswordHash, Role)
		OUTPUT INSERTED.ID, INSERTED.CreatedAt
		VALUES (@p1, @p2, @p3)
	`
	// Use sql.Named args if driver supports it, or just ?/param placeholders depending on driver
	// mssql driver uses @p1, @p2 or ?
	// Let's use standard ? for simplicity if supported, or named args.
	// The go-mssqldb driver supports named parameters.
	
	err := r.DB.QueryRow(query, user.Username, user.PasswordHash, user.Role).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT ID, Username, PasswordHash, Role, CreatedAt
		FROM Users
		WHERE Username = @p1
	`
	err := r.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
