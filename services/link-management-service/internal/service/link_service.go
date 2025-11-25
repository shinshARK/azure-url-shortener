package service

import (
	"errors"
	"time"

	"github.com/shinshark/azure-url-shortener/services/link-management-service/internal/models"
	"github.com/shinshark/azure-url-shortener/services/link-management-service/internal/repository"
	"github.com/teris-io/shortid"
)

type LinkService struct {
	Repo *repository.LinkRepository
}

func NewLinkService(repo *repository.LinkRepository) *LinkService {
	return &LinkService{Repo: repo}
}

func (s *LinkService) CreateLink(req *models.CreateLinkRequest, userID *int, role string) (*models.Link, error) {
	// 1. Quota Check for Users
	if role == "User" && userID != nil {
		count, err := s.Repo.CountLinksByUserID(*userID)
		if err != nil {
			return nil, err
		}
		if count >= 10 { // Hardcoded quota for now
			return nil, errors.New("quota exceeded: max 10 links allowed")
		}
	}

	// 2. Generate Short Code
	var shortCode string
	if req.CustomAlias != "" {
		if role != "User" && role != "Admin" {
			return nil, errors.New("custom alias is only for registered users")
		}
		// Check if alias exists
		existing, err := s.Repo.GetLinkByShortCode(req.CustomAlias)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("alias already taken")
		}
		shortCode = req.CustomAlias
	} else {
		var err error
		shortCode, err = shortid.Generate()
		if err != nil {
			return nil, err
		}
	}

	// 3. Set Expiry
	var expiresAt *time.Time
	if role == "Guest" {
		t := time.Now().Add(24 * time.Hour) // 24 hours for guests
		expiresAt = &t
	}
	// Users have no expiry by default (nil)

	link := &models.Link{
		ShortCode:   shortCode,
		OriginalUrl: req.OriginalUrl,
		UserID:      userID,
		CreatedAt:   time.Now(),
		ExpiresAt:   expiresAt,
		CustomAlias: req.CustomAlias,
		IsActive:    true, // Default to true
	}

	if err := s.Repo.CreateLink(link); err != nil {
		return nil, err
	}

	return link, nil
}

func (s *LinkService) GetUserLinks(userID int) ([]models.Link, error) {
	return s.Repo.GetLinksByUserID(userID)
}

func (s *LinkService) DeleteLink(shortCode string, userID int, role string) error {
	link, err := s.Repo.GetLinkByShortCode(shortCode)
	if err != nil {
		return err
	}
	if link == nil {
		return errors.New("link not found")
	}

	// Authorization Check
	if role != "Admin" {
		if link.UserID == nil || *link.UserID != userID {
			return errors.New("unauthorized")
		}
	}

	return s.Repo.DeleteLink(shortCode)
}
