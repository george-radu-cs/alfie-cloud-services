package usecase

import (
	"api/app/models"
	"api/app/utils"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"time"
)

const (
	passwordSaltBytes = 16
	argonTime         = 1
	argonMemoryUsed   = 65536 // 64 * 1024
	argonThreadsUsed  = 1
	argonKeyLen       = 32
)

func (uc *useCase) Register(ctx context.Context, user *models.User) (err error) {
	err = uc.ValidationsService.UserValidation(user)
	if err != nil {
		utils.ErrorLogger.Printf("Error validating user registration: %s", err.Error())
		return err
	}

	user.Password, user.Salt, err = uc.generateHashedPasswordWithSalt(user.Password)
	if err != nil {
		utils.ErrorLogger.Printf("Error generating password for user: %s", err.Error())
		return errors.New("error processing request")
	}

	err = uc.Repository.CreateUser(user)
	if err != nil {
		utils.ErrorLogger.Printf("Error creating user: %s", err.Error())
		return errors.New("error processing request")
	}
	utils.InfoLogger.Printf("Registered user with email %s", user.Email)

	err = uc.MailVerifierService.SendMailWithRegistrationCode(ctx, user.Email, user.FirstName)
	if err != nil {
		utils.ErrorLogger.Printf("Error sending verification email to verify registration for user: %s", err.Error())
		return errors.New("error processing request")
	}
	utils.InfoLogger.Printf("Successfully sent verification email for created user with email %s", user.Email)

	return nil
}

func (uc *useCase) VerifyUserAccount(ctx context.Context, email, code string) (err error) {
	user, err := uc.Repository.GetUserByEmail(email)
	if uc.didNotReceivedUser(err, user) {
		utils.ErrorLogger.Printf("User not found for email %s", email)
		return errors.New("invalid request")
	}

	if user.Verified {
		utils.InfoLogger.Printf("User with email %s already verified", email)
		return nil
	}

	verified, err := uc.MailVerifierService.CheckRegistrationCode(ctx, email, code)
	if err != nil || !verified {
		utils.ErrorLogger.Printf("Error verifying user with email %s, error %s", email, err.Error())
		return errors.New("invalid request")
	}

	err = uc.createMediaFolderForUser(ctx, user)
	if err != nil {
		utils.ErrorLogger.Printf("Error creating media folder for user with email %s, error %s", email, err.Error())
		return errors.New("error processing request")
	}

	err = uc.markUserAccountAsVerified(user)
	if err != nil {
		utils.ErrorLogger.Printf("Error marking user's account with email %s as verified, error %s", email, err.Error())
		return errors.New("error processing request")
	}

	utils.InfoLogger.Printf("Successfully verified account for user with email %s", email)

	return nil
}

func (uc *useCase) ResendUserVerificationCode(ctx context.Context, email, password string) (err error) {
	user, err := uc.Repository.GetUserByEmail(email)
	if uc.didNotReceivedUser(err, user) {
		utils.ErrorLogger.Printf("User not found for email %s", email)
		return errors.New("invalid credentials")
	}

	if user.Verified {
		utils.InfoLogger.Printf("User with email %s already verified", email)
		return nil
	}

	err = uc.verifyUserPassword(password, user.Password, user.Salt)
	if err != nil {
		utils.ErrorLogger.Printf("Login failed for user with email %s, error %s", email, err.Error())
		return errors.New("invalid credentials")
	}

	err = uc.MailVerifierService.SendMailWithRegistrationCode(ctx, user.Email, user.FirstName)
	if err != nil {
		utils.ErrorLogger.Printf("Error sending verification email to verify registration for user: %s", err.Error())
		return errors.New("error processing request")
	}
	utils.InfoLogger.Printf("Successfully sent verification email for created user with email %s", user.Email)

	return nil
}

func (uc *useCase) Login(ctx context.Context, email, password string) (err error) {
	user, err := uc.Repository.GetUserByEmail(email)
	if uc.didNotReceivedUser(err, user) {
		utils.ErrorLogger.Printf("User not found for email %s", email)
		return errors.New("invalid credentials")
	}

	if !user.Verified {
		utils.ErrorLogger.Printf("User with email %s hasn't verified his email", email)
		return errors.New("email not verified")
	}

	err = uc.verifyUserPassword(password, user.Password, user.Salt)
	if err != nil {
		utils.ErrorLogger.Printf("Login failed for user with email %s, error %s", email, err.Error())
		return errors.New("invalid credentials")
	}

	err = uc.markUserCanGoTo2FAStep(user)
	if err != nil {
		utils.ErrorLogger.Printf("Error marking user with email %s to go to 2fa step, error %s", email, err.Error())
		return errors.New("invalid credentials")
	}

	err = uc.MailVerifierService.SendMailWith2FALoginCode(ctx, user.Email, user.FirstName)
	if err != nil {
		utils.ErrorLogger.Printf("Error sending verification email to log in 2fa for user: %v", err)
		return errors.New("error processing request")
	}
	utils.InfoLogger.Printf("Successfully sent verification email for login request of user with email %s", user.Email)

	utils.InfoLogger.Printf("User with email %s passed the login step, waiting for 2fa", email)

	return nil
}

func (uc *useCase) VerifyLoginCode(ctx context.Context, email, code string) (user *models.User, err error) {
	user, err = uc.Repository.GetUserByEmail(email)
	if uc.didNotReceivedUser(err, user) {
		utils.ErrorLogger.Printf("User not found for email %s", email)
		return nil, errors.New("invalid credentials")
	}

	if !user.LoginCanCheck2FA {
		utils.ErrorLogger.Printf("User with email %s can't check 2fa code", email)
		return nil, errors.New("invalid credentials")
	}

	if userExceeded2FACodeTime(user) {
		utils.ErrorLogger.Printf("User with email %s has exceeded the 2fa code validity time", email)
		return nil, errors.New("invalid credentials")
	}

	verified, err := uc.MailVerifierService.Check2FALoginCode(ctx, email, code)
	if err != nil || !verified {
		utils.ErrorLogger.Printf("Error verifying 2fa for user with email %s, error %s", email, err.Error())
		return nil, errors.New("invalid credentials")
	}

	err = uc.markUserPassed2FAStep(user)
	if err != nil {
		utils.ErrorLogger.Printf(
			"Error updating user with email %s as passed the 2fa step, error %s", email,
			err.Error(),
		)
		return nil, errors.New("invalid credentials")
	}

	utils.InfoLogger.Printf("Successfully login with 2fa for user with email %s", email)

	return user, nil
}

func (uc *useCase) ForgotPassword(ctx context.Context, email string) (err error) {
	user, err := uc.Repository.GetUserByEmail(email)
	if uc.didNotReceivedUser(err, user) {
		utils.ErrorLogger.Printf("User not found for email %s", email)
		return nil // if user not found, we don't want to tell the user that
	}

	if !user.Verified {
		utils.ErrorLogger.Printf("User with email %s hasn't verified his email", email)
		return errors.New("email not verified")
	}

	err = uc.MailVerifierService.SendMailWithForgotPasswordCode(ctx, user.Email)
	if err != nil {
		utils.ErrorLogger.Printf("Error sending verification email to reset password for user: %v", err)
		return errors.New("error processing request")
	}
	utils.InfoLogger.Printf("Successfully sent verification email for login request of user with email %s", user.Email)

	return nil
}

func (uc *useCase) ResetPassword(ctx context.Context, email, code, newPassword string) (err error) {
	validRequest := uc.ValidationsService.IsPasswordValid(newPassword)
	if !validRequest {
		utils.ErrorLogger.Printf("Invalid new password for user with email %s", email)
		return errors.New("invalid password")
	}

	user, err := uc.Repository.GetUserByEmail(email)
	if uc.didNotReceivedUser(err, user) {
		utils.ErrorLogger.Printf("User not found for email %s", email)
		return nil // if user not found, we don't want to tell the user that
	}

	if !user.Verified {
		utils.ErrorLogger.Printf("User with email %s hasn't verified his email", email)
		return errors.New("email not verified")
	}

	verified, err := uc.MailVerifierService.CheckForgotPasswordCode(ctx, email, code)
	if err != nil || !verified {
		utils.ErrorLogger.Printf("Error verifying 2fa for user with email %s, error %s", email, err.Error())
		return errors.New("invalid credentials")
	}

	user.Password, user.Salt, err = uc.generateHashedPasswordWithSalt(newPassword)
	if err != nil {
		utils.ErrorLogger.Printf("Error generating password for user: %s", err.Error())
		return errors.New("error processing request")
	}

	err = uc.Repository.UpdateUser(user)
	if err != nil {
		utils.ErrorLogger.Printf("Error updating user: %s", err.Error())
		return errors.New("error processing request")
	}
	utils.InfoLogger.Printf("Successfully reseted forgotten password for user with email %s", user.Email)

	return nil
}

func (uc *useCase) generateSalt() (salt []byte, err error) {
	salt = make([]byte, passwordSaltBytes)
	_, err = rand.Read(salt)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error generating salt: %s", err.Error()))
	}
	return salt, nil
}

func (uc *useCase) hashPassword(password string, salt []byte) (hashedPassword []byte) {
	key := argon2.IDKey([]byte(password), salt, argonTime, argonMemoryUsed, argonThreadsUsed, argonKeyLen)
	return key
}

// returns the hashed password and the salt as hex strings and possible error
func (uc *useCase) generateHashedPasswordWithSalt(password string) (
	hashedPasswordAsHex string, saltAsHex string, err error,
) {
	salt, err := uc.generateSalt()
	if err != nil {
		return "", "", err
	}
	hashedPassword := uc.hashPassword(password, salt)

	// encode the salt and the password as hex strings
	return hex.EncodeToString(hashedPassword), hex.EncodeToString(salt), nil
}

func (uc *useCase) verifyUserPassword(password, storedHashedPasswordAsHex, saltAsHex string) (err error) {
	salt, err := hex.DecodeString(saltAsHex)
	if err != nil {
		return err
	}

	givenHashedPassword := uc.hashPassword(password, salt)
	givenHashedPasswordAsHex := hex.EncodeToString(givenHashedPassword)

	if givenHashedPasswordAsHex != storedHashedPasswordAsHex {
		return errors.New("password mismatch")
	}

	return nil
}

func (uc *useCase) markUserAccountAsVerified(user *models.User) (err error) {
	user.Verified = true
	err = uc.Repository.UpdateUser(user)
	return err
}

func (uc *useCase) markUserCanGoTo2FAStep(user *models.User) (err error) {
	user.LoginCanCheck2FA = true
	requestTime := time.Now()
	user.O2FARequestedAt = &requestTime
	err = uc.Repository.UpdateUser(user)
	return err
}

func (uc *useCase) markUserPassed2FAStep(user *models.User) (err error) {
	user.LoginCanCheck2FA = false
	user.O2FARequestedAt = nil
	err = uc.Repository.UpdateUser(user)
	return err
}

func userExceeded2FACodeTime(user *models.User) (res bool) {
	return time.Now().Sub(*user.O2FARequestedAt) > time.Duration(10)*time.Minute
}
