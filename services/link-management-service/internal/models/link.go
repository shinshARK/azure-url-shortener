package models

import "time"

type Link struct {
	ShortCode   string     `json:"shortCode"`
	OriginalUrl string     `json:"originalUrl"` // Renamed from LongUrl
	UserID      *int       `json:"userId,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty"`
	ClickCount  int        `json:"clickCount"`
	CustomAlias string     `json:"customAlias,omitempty"`
	IsActive    bool       `json:"isActive"` // Added IsActive
}

type CreateLinkRequest struct {
	OriginalUrl string `json:"originalUrl" binding:"required,url"` // Renamed from LongUrl
	CustomAlias string `json:"customAlias"`
}
