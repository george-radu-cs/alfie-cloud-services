package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type MarketDeckCard struct {
	gorm.Model
	ID                       uuid.UUID     `json:"id" gorm:"primaryKey;type:uuid"`
	MarketDeckID             uuid.UUID     `json:"marketDeckId" gorm:"type:uuid"`
	LocalCardID              int           `json:"localCardId" gorm:"not null;"`
	LocalDeckID              int           `json:"localDeckId" gorm:"not null;"`
	CardType                 string        `json:"cardType" gorm:"not null;"`
	Question                 string        `json:"question" gorm:"not null;"`
	QuestionTextType         string        `json:"questionTextType" gorm:"not null;"`
	QuestionAttachmentType   string        `json:"questionAttachmentType" gorm:"not null;"`
	QuestionAttachmentURI    string        `json:"questionAttachmentURI" gorm:"not null;"`
	Answer                   string        `json:"answer" gorm:"not null;"`
	AnswerTextType           string        `json:"answerTextType" gorm:"not null;"`
	AnswerNumberOfOptions    int           `json:"answerNumberOfOptions"`
	AnswerCorrectOptionIndex int           `json:"answerCorrectOptionIndex"`
	AnswerAttachmentType     string        `json:"answerAttachmentType" gorm:"not null;"`
	AnswerAttachmentURI      string        `json:"answerAttachmentURI" gorm:"not null;"`
	TimeToAnswer             time.Duration `json:"timeToAnswer" gorm:"not null;"`
	CreatedAt                time.Time     `json:"createdAt"`
	UpdatedAt                time.Time     `json:"UpdatedAt"`
	DeletedAt                time.Time     `json:"deletedAt" gorm:"default:null"`
}

func (marketDeckCard *MarketDeckCard) BeforeCreate(tx *gorm.DB) (err error) {
	marketDeckCard.ID = uuid.New()

	return
}
