package repository

import (
	"database/sql"
	"fmt"

	"github.com/shinshark/azure-url-shortener/services/link-management-service/internal/models"
)

type LinkRepository struct {
	DB *sql.DB
}

func NewLinkRepository(db *sql.DB) *LinkRepository {
	return &LinkRepository{DB: db}
}

func (r *LinkRepository) CreateLink(link *models.Link) error {
	query := `
		INSERT INTO Links (ShortCode, OriginalUrl, UserID, CreatedAt, ExpiresAt, ClickCount, CustomAlias, IsActive)
		VALUES (@p1, @p2, @p3, @p4, @p5, 0, @p6, @p7)
	`
	_, err := r.DB.Exec(query, link.ShortCode, link.OriginalUrl, link.UserID, link.CreatedAt, link.ExpiresAt, link.CustomAlias, link.IsActive)
	if err != nil {
		return fmt.Errorf("failed to create link: %w", err)
	}
	return nil
}

func (r *LinkRepository) GetLinksByUserID(userID int) ([]models.Link, error) {
	query := `
		SELECT ShortCode, OriginalUrl, UserID, CreatedAt, ExpiresAt, ClickCount, CustomAlias, IsActive
		FROM Links
		WHERE UserID = @p1
	`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []models.Link
	for rows.Next() {
		var l models.Link
		if err := rows.Scan(&l.ShortCode, &l.OriginalUrl, &l.UserID, &l.CreatedAt, &l.ExpiresAt, &l.ClickCount, &l.CustomAlias, &l.IsActive); err != nil {
			return nil, err
		}
		links = append(links, l)
	}
	return links, nil
}

func (r *LinkRepository) CountLinksByUserID(userID int) (int, error) {
	query := "SELECT COUNT(*) FROM Links WHERE UserID = @p1"
	var count int
	err := r.DB.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *LinkRepository) GetLinkByShortCode(code string) (*models.Link, error) {
	query := `
		SELECT ShortCode, OriginalUrl, UserID, CreatedAt, ExpiresAt, ClickCount, CustomAlias, IsActive
		FROM Links
		WHERE ShortCode = @p1
	`
	var l models.Link
	err := r.DB.QueryRow(query, code).Scan(&l.ShortCode, &l.OriginalUrl, &l.UserID, &l.CreatedAt, &l.ExpiresAt, &l.ClickCount, &l.CustomAlias, &l.IsActive)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *LinkRepository) DeleteLink(code string) error {
	query := "DELETE FROM Links WHERE ShortCode = @p1"
	_, err := r.DB.Exec(query, code)
	return err
}

func (r *LinkRepository) CountCustomLinksByUserID(userID int) (int, error) {
	query := "SELECT COUNT(*) FROM Links WHERE UserID = @p1 AND CustomAlias IS NOT NULL AND CustomAlias <> ''"
	var count int
	err := r.DB.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *LinkRepository) CountStandardLinksByUserID(userID int) (int, error) {
	query := "SELECT COUNT(*) FROM Links WHERE UserID = @p1 AND (CustomAlias IS NULL OR CustomAlias = '')"
	var count int
	err := r.DB.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *LinkRepository) UpdateLink(link *models.Link) error {
	query := `
		UPDATE Links
		SET OriginalUrl = @p1
		WHERE ShortCode = @p2
	`
	_, err := r.DB.Exec(query, link.OriginalUrl, link.ShortCode)
	if err != nil {
		return fmt.Errorf("failed to update link: %w", err)
	}
	return nil
}
