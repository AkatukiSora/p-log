package security

import (
	"context"
	"errors"
	"fmt"
	"log"

	"backend/api"
	"backend/ent"
	"backend/internal/jwt"
)

var (
	// ErrMissingToken はトークンが存在しない場合のエラーです。
	ErrMissingToken = errors.New("missing token")
)

// SecurityHandler はJWT認証を処理するセキュリティハンドラーです。
type SecurityHandler struct {
	jwtHandler *jwt.JwtHandler
	client     *ent.Client
}

// NewSecurityHandler は新しいSecurityHandlerインスタンスを作成します。
func NewSecurityHandler(s *jwt.JwtHandler, client *ent.Client) *SecurityHandler {
	return &SecurityHandler{
		jwtHandler: s,
		client:     client,
	}
}

// HandleBearerAuth はBearer認証トークンを検証します。
// JWT認証の実装はこのメソッドで行います。
func (s *SecurityHandler) HandleBearerAuth(ctx context.Context, operationName api.OperationName, t api.BearerAuth) (context.Context, error) {
	// トークンの存在確認
	if t.Token == "" {
		return ctx, fmt.Errorf("authentication failed: %w", ErrMissingToken)
	}

	// JWTトークンを検証
	claims, err := s.jwtHandler.ValidateToken(t.Token)
	if err != nil {
		log.Printf("token validation failed: %v", err)
		// JWT関連のエラーをそのまま返して、handler.NewErrorで401として処理されるようにする
		return ctx, fmt.Errorf("authentication failed: %w", err)
	}

	// 検証成功後、ユーザー情報をcontextに保存
	ctx = context.WithValue(ctx, userIDKey, claims.UserID)

	return ctx, nil
}

// ContextKeys for storing user information
type contextKey string

const (
	// userIDKey はcontextからユーザーIDを取得するためのキーです。
	userIDKey contextKey = "user_id"
)

// GetUserIDFromContext はcontextからユーザーIDを取得します。
func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok {
		return "", errors.New("ユーザーIDをコンテキストから取得できません")
	}
	return userID, nil
}
