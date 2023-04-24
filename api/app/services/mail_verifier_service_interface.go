package services

type MailVerifierService interface {
	SendMailWithRegistrationCode(email, firstName string) (err error)
	SendMailWith2FALoginCode(email, firstName string) (err error)
	SendMailWithForgotPasswordCode(email string) (err error)
	CheckRegistrationCode(email, code string) (verified bool, err error)
	Check2FALoginCode(email, code string) (verified bool, err error)
	CheckForgotPasswordCode(email, code string) (verified bool, err error)
}
