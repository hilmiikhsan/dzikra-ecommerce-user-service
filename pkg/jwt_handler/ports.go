package jwt_handler

import "context"

type JWT interface {
	GenerateTokenString(ctx context.Context, payload CostumClaimsPayload) (*GenerateTokenResponse, error)
	ParseTokenString(ctx context.Context, tokenString, email, sessionID, tokenType string) (*CustomClaims, error)
	ParseMiddlewareTokenString(ctx context.Context, tokenString string) (*CustomClaims, error)
	ParseMiddlewareRefreshTokenString(ctx context.Context, refreshTokenString string) (*CustomClaims, error)
}
