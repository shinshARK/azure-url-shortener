package service

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/shinshark/azure-url-shortener/services/link-management-service/internal/models"
	"github.com/shinshark/azure-url-shortener/services/link-management-service/internal/repository"
	"github.com/teris-io/shortid"
)

type LinkService struct {
	Repo             *repository.LinkRepository
	CacheEvictionUrl string
}

func NewLinkService(repo *repository.LinkRepository, cacheEvictionUrl string) *LinkService {
	return &LinkService{
		Repo:             repo,
		CacheEvictionUrl: cacheEvictionUrl,
	}
}

func (s *LinkService) evictCache(shortCode string) {
	if s.CacheEvictionUrl == "" {
		return
	}
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", s.CacheEvictionUrl, shortCode), nil)
	if err != nil {
		fmt.Printf("Failed to create cache eviction request: %v\n", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to evict cache: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		fmt.Printf("Cache eviction failed with status: %d\n", resp.StatusCode)
	}
}

func (s *LinkService) CreateLink(req *models.CreateLinkRequest, userID *int, role string) (*models.Link, error) {
	// ... (existing code) ...
	// 1. Quota Check for Users
	if role == "User" && userID != nil {
		if req.CustomAlias != "" {
			// Custom Link Quota
			count, err := s.Repo.CountCustomLinksByUserID(*userID)
			if err != nil {
				return nil, err
			}
			if count >= 2 {
				return nil, errors.New("quota exceeded: max 2 custom links allowed")
			}
		} else {
			// Standard Link Quota
			count, err := s.Repo.CountStandardLinksByUserID(*userID)
			if err != nil {
				return nil, err
			}
			if count >= 20 {
				return nil, errors.New("quota exceeded: max 20 standard links allowed")
			}
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

func (s *LinkService) UpdateLink(shortCode string, req *models.UpdateLinkRequest, userID int, role string) (*models.Link, error) {
	link, err := s.Repo.GetLinkByShortCode(shortCode)
	if err != nil {
		return nil, err
	}
	if link == nil {
		return nil, errors.New("link not found")
	}

	// Authorization Check
	if role != "Admin" {
		if link.UserID == nil || *link.UserID != userID {
			return nil, errors.New("unauthorized")
		}
	}

	// Validation: Only custom alias links can be edited
	if link.CustomAlias == "" {
		return nil, errors.New("only custom alias links can be edited")
	}

	link.OriginalUrl = req.OriginalUrl
	if err := s.Repo.UpdateLink(link); err != nil {
		return nil, err
	}

	// Evict Cache
	go s.evictCache(shortCode)

	return link, nil
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

	if err := s.Repo.DeleteLink(shortCode); err != nil {
		return err
	}

	// Evict Cache
	go s.evictCache(shortCode)

	return nil
}
