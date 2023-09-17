package repository

import (
	"gorm.io/gorm"
)

type repository struct {
	Db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{Db: db}
}
