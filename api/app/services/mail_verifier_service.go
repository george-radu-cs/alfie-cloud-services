package services

import (
	"api/app/utils"
	"fmt"
	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
	"os"
)

const (
	substitutionsKey                 = "substitutions"
	firstNameTemplateSubstitutionKey = "first_name"
)

type mailVerifierService struct {
	verifyRegistrationSID   string
	verifyLoginSID          string
	verifyForgotPasswordSID string
	client                  *twilio.RestClient
}

func NewMailVerifierService() MailVerifierService {
	return &mailVerifierService{
		verifyRegistrationSID:   os.Getenv("TWILIO_VERIFY_REGISTRATION_SID"),
		verifyLoginSID:          os.Getenv("TWILIO_VERIFY_LOGIN_SID"),
		verifyForgotPasswordSID: os.Getenv("TWILIO_VERIFY_FORGOT_PASSWORD_SID"),
		client: twilio.NewRestClientWithParams(
			twilio.ClientParams{
				Username: os.Getenv("TWILIO_ACCOUNT_SID"),
				Password: os.Getenv("TWILIO_AUTH_TOKEN"),
			},
		),
	}
}

func (mvs *mailVerifierService) SendMailWithRegistrationCode(email, firstName string) (err error) {
	params := &verify.CreateVerificationParams{}
	params.SetChannelConfiguration(
		map[string]interface{}{
			substitutionsKey: map[string]interface{}{
				firstNameTemplateSubstitutionKey: firstName,
			},
		},
	)
	params.SetTo(email)
	params.SetChannel("email")

	resp, err := mvs.client.VerifyV2.CreateVerification(
		mvs.verifyRegistrationSID, params,
	)
	if err != nil {
		return fmt.Errorf("failed to send verify registration email for %s: %w", email, err)
	} else {
		if resp.Sid != nil {
			utils.InfoLogger.Printf("Sent verify registration email to %s; Response SID %s", email, *resp.Sid)
		} else {
			utils.InfoLogger.Printf("Sent verify registration email to %s; Response SID %s", email, resp.Sid)
		}

		return nil
	}
}

func (mvs *mailVerifierService) SendMailWith2FALoginCode(email, firstName string) (err error) {
	params := &verify.CreateVerificationParams{}
	params.SetChannelConfiguration(
		map[string]interface{}{
			substitutionsKey: map[string]interface{}{
				firstNameTemplateSubstitutionKey: firstName,
			},
		},
	)
	params.SetTo(email)
	params.SetChannel("email")

	resp, err := mvs.client.VerifyV2.CreateVerification(
		mvs.verifyLoginSID, params,
	)
	if err != nil {
		return fmt.Errorf("failed to send verify 2fa login email for %s: %w", email, err)
	} else {
		if resp.Sid != nil {
			utils.InfoLogger.Printf("Sent verify 2fa login email to %s; Response SID %s", email, *resp.Sid)
		} else {
			utils.InfoLogger.Printf("Sent verify 2fa login email to %s; Response SID %s", email, resp.Sid)
		}

		return nil
	}
}

func (mvs *mailVerifierService) SendMailWithForgotPasswordCode(email string) (err error) {
	params := &verify.CreateVerificationParams{}
	params.SetTo(email)
	params.SetChannel("email")

	resp, err := mvs.client.VerifyV2.CreateVerification(
		mvs.verifyForgotPasswordSID, params,
	)
	if err != nil {
		return fmt.Errorf("failed to send verify forgot password email for %s: %w", email, err)
	} else {
		if resp.Sid != nil {
			utils.InfoLogger.Printf("Sent verify forgot password email to %s; Response SID %s", email, *resp.Sid)
		} else {
			utils.InfoLogger.Printf("Sent verify forgot password email to %s; Response SID %s", email, resp.Sid)
		}

		return nil
	}
}

func (mvs *mailVerifierService) CheckRegistrationCode(email, code string) (verified bool, err error) {
	verified, err = mvs.checkVerificationCodeForService(
		email, code,
		mvs.verifyRegistrationSID,
	)
	if err != nil {
		return false, fmt.Errorf("failed to verify registration code for email %s: %w", email, err)
	}

	return verified, nil
}

func (mvs *mailVerifierService) Check2FALoginCode(email, code string) (verified bool, err error) {
	verified, err = mvs.checkVerificationCodeForService(
		email, code,
		mvs.verifyLoginSID,
	)
	if err != nil {
		return false, fmt.Errorf("failed to verify 2fa login code for email %s: %w", email, err)
	}

	return verified, nil
}

func (mvs *mailVerifierService) CheckForgotPasswordCode(email, code string) (verified bool, err error) {
	verified, err = mvs.checkVerificationCodeForService(
		email, code,
		mvs.verifyForgotPasswordSID,
	)
	if err != nil {
		return false, fmt.Errorf("failed to verify forgot password code for email %s: %w", email, err)
	}

	return verified, nil
}

func (mvs *mailVerifierService) checkVerificationCodeForService(email, code, serviceId string) (
	verified bool, err error,
) {
	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(email)
	params.SetCode(code)

	resp, err := mvs.client.VerifyV2.CreateVerificationCheck(serviceId, params)
	if err != nil {
		return false, err
	} else {
		if resp.Sid != nil {
			utils.InfoLogger.Printf(
				"Verified code for %s; Response SID %s received status %s", email, *resp.Sid, *resp.Status,
			)

			if *resp.Status != "approved" {
				return false, nil
			}

			return true, nil

		} else {
			utils.InfoLogger.Printf(
				"Verified code for %s; Response SID %s received status %s", email, resp.Sid, resp.Status,
			)

			if *resp.Status != "approved" {
				return false, nil
			}

			return true, nil
		}
	}
}
