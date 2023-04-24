package repository

import (
	"api/app/models"
)

type Repository interface {
	CreateUser(user *models.User) (err error)
	UpdateUser(user *models.User) (err error)
	GetUserByEmail(email string) (user *models.User, err error)
}
