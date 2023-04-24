package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type MarketDeckReview struct {
	gorm.Model
	ID           uuid.UUID `json:"id" gorm:"primaryKey;type:uuid"`
	MarketDeckID uuid.UUID `json:"marketDeckId" gorm:"type:uuid"`
	UserID       uuid.UUID `json:"userId" gorm:"type:uuid"`
	Title        string    `json:"title" gorm:"not null;"`
	Description  string    `json:"description" gorm:"not null;"`
	Stars        int       `json:"stars" gorm:"not null;"`
	Verified     bool      `json:"verified" gorm:"not null;default:false;"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"UpdatedAt"`
	DeletedAt    time.Time `json:"deletedAt" gorm:"default:null"`
}

func (marketDeckReview *MarketDeckReview) BeforeCreate(tx *gorm.DB) (err error) {
	marketDeckReview.ID = uuid.New()
	marketDeckReview.Verified = false

	return
}
