package services

type JWTService interface {
	GenerateToken(email string) (jwtToken string, err error)
	ValidateToken(tokenString string) (tokenClaims *Claims, err error)
}
