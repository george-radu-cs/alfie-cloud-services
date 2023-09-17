package services

import "context"

type MailVerifierService interface {
	SendMailWithRegistrationCode(ctx context.Context, email, firstName string) (err error)
	SendMailWith2FALoginCode(ctx context.Context, email, firstName string) (err error)
	SendMailWithForgotPasswordCode(ctx context.Context, email string) (err error)
	CheckRegistrationCode(ctx context.Context, email, code string) (verified bool, err error)
	Check2FALoginCode(ctx context.Context, email, code string) (verified bool, err error)
	CheckForgotPasswordCode(ctx context.Context, email, code string) (verified bool, err error)
}
