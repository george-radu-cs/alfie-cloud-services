package services

import (
	"errors"
	"net/mail"

	"api/app/models"

	"github.com/dlclark/regexp2"
)

const (
	// include at least one uppercase letter, one lowercase letter, one number, one special character
	// and be at least 12 characters long
	passwordRegex = `^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[!@#\$&\(\)\[\]<>%^-_=+*~]).{12,}$`
	minNameLength = 3
)

type validationsService struct {
	passwordRegex *regexp2.Regexp
}

func NewValidationsService() ValidationsService {
	return &validationsService{
		passwordRegex: regexp2.MustCompile(passwordRegex, 0),
	}
}

func (vs *validationsService) IsNameValid(name string) (res bool) {
	return len(name) >= minNameLength
}

func (vs *validationsService) IsEmailValid(email string) (res bool) {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (vs *validationsService) IsPasswordValid(password string) (res bool) {
	passwordValid, _ := vs.passwordRegex.MatchString(password)
	return passwordValid
}

func (vs *validationsService) UserValidation(user *models.User) (err error) {
	if vs.IsEmailValid(user.Email) && vs.IsPasswordValid(user.Password) &&
		vs.IsNameValid(user.FirstName) && vs.IsNameValid(user.LastName) {
		return nil
	}

	return errors.New("invalid user registration data")
}

func (vs *validationsService) UserInfoValidation(user *models.User) (err error) {
	if vs.IsNameValid(user.FirstName) && vs.IsNameValid(user.LastName) {
		return nil
	}

	return errors.New("invalid user info data")
}
