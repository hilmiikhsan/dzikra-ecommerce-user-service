package jwt_handler

import "context"

type JWT interface {
	GenerateTokenString(ctx context.Context, payload CostumClaimsPayload) (string, string, error)
	ParseTokenString(ctx context.Context, tokenString, username, tokenType string) (*CustomClaims, error)
	ParseMiddlewareTokenString(ctx context.Context, tokenString string) (*CustomClaims, error)
}
