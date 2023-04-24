package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	ID                      uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid"`
	Email                   string     `json:"email" gorm:"unique;not null"`
	FirstName               string     `json:"firstName" gorm:"not null;"`
	LastName                string     `json:"lastName" gorm:"not null;"`
	Password                string     `json:"password" gorm:"not null;"`
	Salt                    string     `json:"salt" gorm:"not null;"`
	S3ID                    string     `json:"s3Id" gorm:"unique"`
	S3MaxNumberOfMediaFiles int32      `json:"s3MaxNumberOfMediaFiles" gorm:"not null;default:1000"`
	Verified                bool       `json:"verified" gorm:"not null;default:false;"`
	LoginCanCheck2FA        bool       `json:"passLogin" gorm:"not null;default:false;"`
	O2FARequestedAt         *time.Time `json:"o2FARequestedAt" gorm:"default:null;"`
	CreatedAt               *time.Time `json:"createdAt"`
	UpdatedAt               *time.Time `json:"updatedAt"`
	DeletedAt               *time.Time `json:"deletedAt" gorm:"default:null"`

	MarketDeck       []MarketDeck       `json:"marketDeck" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MarketDeckReview []MarketDeckReview `json:"marketDeckReview" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	user.S3ID = uuid.New().String()
	user.S3MaxNumberOfMediaFiles = 1000
	user.Verified = false
	user.LoginCanCheck2FA = false
	user.O2FARequestedAt = nil

	return
}
