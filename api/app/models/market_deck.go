package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type MarketDeck struct {
	gorm.Model
	ID          uuid.UUID `json:"id" gorm:"primaryKey;type:uuid"`
	UserID      uuid.UUID `json:"userId" gorm:"type:uuid"`
	LocalDeckID int       `json:"localDeckId" gorm:"not null;"`
	Name        string    `json:"name" gorm:"not null;"`
	Description string    `json:"description" gorm:"not null;"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
	DeletedAt   time.Time `json:"deletedAt" gorm:"default:null"`

	MarketDeckImage  []MarketDeckImage  `json:"marketDeckImage" gorm:"foreignKey:MarketDeckID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MarketDeckCard   []MarketDeckCard   `json:"marketDeckCard" gorm:"foreignKey:MarketDeckID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MarketDeckReview []MarketDeckReview `json:"marketDeckReview" gorm:"foreignKey:MarketDeckID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (marketDeck *MarketDeck) BeforeCreate(tx *gorm.DB) (err error) {
	marketDeck.ID = uuid.New()

	return
}
