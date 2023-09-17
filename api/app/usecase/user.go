package usecase

import (
	"context"
	"errors"

	"api/app/models"
	"api/app/utils"
)

func (uc *useCase) UpdateUserInfo(ctx context.Context, firstName, lastName, userEmail string) (err error) {
	err = uc.ValidationsService.UserInfoValidation(
		&models.User{
			FirstName: firstName,
			LastName:  lastName,
		},
	)
	if err != nil {
		utils.ErrorLogger.Printf("Error validating user info update: %s", err.Error())
		return err
	}

	user, err := uc.Repository.GetUserByEmail(userEmail)
	if uc.didNotReceivedUser(err, user) {
		utils.ErrorLogger.Printf("User not found for email %s", userEmail)
		return errors.New("invalid request")
	}

	user.FirstName = firstName
	user.LastName = lastName

	err = uc.Repository.UpdateUser(user)
	if err != nil {
		utils.ErrorLogger.Printf("Error updating user info: %s", err.Error())
		return errors.New("error processing request")
	}

	return nil
}

func (uc *useCase) UpdateUserPassword(ctx context.Context, oldPassword, newPassword, userEmail string) (err error) {
	if !uc.ValidationsService.IsPasswordValid(newPassword) {
		utils.ErrorLogger.Printf("Error validating new user password")
		return errors.New("invalid new password")
	}

	user, err := uc.Repository.GetUserByEmail(userEmail)
	if uc.didNotReceivedUser(err, user) {
		utils.ErrorLogger.Printf("User not found for email %s", userEmail)
		return errors.New("invalid request")
	}

	err = uc.verifyUserPassword(oldPassword, user.Password, user.Salt)
	if err != nil {
		utils.ErrorLogger.Printf(
			"Verification for user's old password failed. User with email %s, error %s",
			userEmail,
			err.Error(),
		)
		return errors.New("invalid old password")
	}

	user.Password, user.Salt, err = uc.generateHashedPasswordWithSalt(newPassword)
	if err != nil {
		utils.ErrorLogger.Printf("Error generating new password for user: %s", err.Error())
		return errors.New("error processing request")
	}

	err = uc.Repository.UpdateUser(user)
	if err != nil {
		utils.ErrorLogger.Printf("Error updating password for user with email %s, error %s", userEmail, err.Error())
		return errors.New("error processing request")
	}

	return nil
}

func (uc *useCase) didNotReceivedUser(err error, user *models.User) (res bool) {
	return err != nil || user == nil
}
