package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type MarketDeckImage struct {
	gorm.Model
	ID           uuid.UUID `json:"id" gorm:"primaryKey;type:uuid"`
	MarketDeckID uuid.UUID `json:"marketDeckId" gorm:"type:uuid"`
	URI          string    `json:"uri" gorm:"not null;"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"UpdatedAt"`
	DeletedAt    time.Time `json:"deletedAt" gorm:"default:null"`
}

func (marketDeckImage *MarketDeckImage) BeforeCreate(tx *gorm.DB) (err error) {
	marketDeckImage.ID = uuid.New()

	return
}
