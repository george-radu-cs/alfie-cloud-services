package services

import (
	"os"
	"time"

	"api/app/utils"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}

type jwtService struct {
	jwtSecretKey    []byte
	tokenExpiration time.Duration
}

func NewJWTService() JWTService {
	tokenDuration, err := time.ParseDuration(os.Getenv("JWT_TOKEN_EXPIRATION"))
	if err != nil {
		utils.ErrorLogger.Fatalf("unable to parse env var JWT_TOKEN_DURATION as time.Duration: %v", err)
	}

	return &jwtService{
		jwtSecretKey:    []byte(os.Getenv("JWT_KEY")),
		tokenExpiration: tokenDuration,
	}
}

func (jwts *jwtService) GenerateToken(email string) (jwtToken string, err error) {
	expirationTime := time.Now().Add(jwts.tokenExpiration)
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: expirationTime,
			},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(jwts.jwtSecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (jwts *jwtService) ValidateToken(tokenString string) (tokenClaims *Claims, err error) {
	token, err := jwt.ParseWithClaims(
		tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwts.jwtSecretKey, nil
		},
	)
	if err != nil {
		return nil, err
	}
	tokenClaims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, err
	}
	return tokenClaims, nil
}
