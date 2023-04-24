package services

import "api/app/models"

type ValidationsService interface {
	IsNameValid(name string) (res bool)
	IsEmailValid(email string) (res bool)
	IsPasswordValid(password string) (res bool)
	UserValidation(user *models.User) (err error)
	UserInfoValidation(user *models.User) (err error)
}
