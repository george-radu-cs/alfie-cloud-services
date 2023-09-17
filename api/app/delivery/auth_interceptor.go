package delivery

import (
	"context"
	"errors"
	"strings"

	"api/app/services"
	"api/app/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	JWTService                 services.JWTService
	routesWithoutAuthorization map[string]bool
}

func NewAuthInterceptor(jwtService services.JWTService) *AuthInterceptor {
	return &AuthInterceptor{
		JWTService:                 jwtService,
		routesWithoutAuthorization: computeRoutesWhoDoNotRequireAuthorization(),
	}
}

func (interceptor *AuthInterceptor) Authorize() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		utils.InfoLogger.Printf("AuthInterceptor.Authorize() -> info.FullMethod: %s", info.FullMethod)

		claimUserEmail, err := interceptor.authorizeAndGetUserEmailFromClaims(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		md, _ := metadata.FromIncomingContext(ctx)
		md.Append("user_email", claimUserEmail)
		newCtx := metadata.NewIncomingContext(ctx, md)
		return handler(newCtx, req)
	}
}

func (interceptor *AuthInterceptor) authorizeAndGetUserEmailFromClaims(
	ctx context.Context,
	methodName string,
) (string, error) {
	if interceptor.routesWithoutAuthorization[methodName] {
		return "", nil
	}

	authToken, err := extractTokenFromMetadata(ctx)
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, "access token is required: %v", err)
	}

	claims, err := interceptor.JWTService.ValidateToken(authToken)
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	return claims.Email, nil
}

func computeRoutesWhoDoNotRequireAuthorization() map[string]bool {
	const alfieServicePath = "/alfie.protobuf.Alfie/"
	return map[string]bool{
		alfieServicePath + "Register":                   true,
		alfieServicePath + "VerifyUserAccount":          true,
		alfieServicePath + "ResendUserVerificationCode": true,
		alfieServicePath + "Login":                      true,
		alfieServicePath + "VerifyLoginCode":            true,
		alfieServicePath + "ForgotPassword":             true,
		alfieServicePath + "ResetPassword":              true,
	}
}

func extractTokenFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("metadata not provided")
	}
	tokenValues := md.Get("authorization")
	if len(tokenValues) == 0 {
		return "", errors.New("authorization token not provided")
	}
	// get the token from the Bearer scheme
	splitToken := strings.Split(tokenValues[0], "Bearer ")
	if len(splitToken) != 2 {
		return "", errors.New("malformed token")
	}
	return splitToken[1], nil
}

func getUserEmailFromValidatedContext(ctx context.Context) string {
	md, _ := metadata.FromIncomingContext(ctx)
	return md.Get("user_email")[0]
}
