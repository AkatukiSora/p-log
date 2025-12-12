package security

import (
	"context"
	"errors"

	"backend/api"
	"backend/util"
)

var (
	// ErrMissingToken はトークンが存在しない場合のエラーです。
	ErrMissingToken = errors.New("missing token")
)

// SecurityHandler はJWT認証を処理するセキュリティハンドラーです。
type SecurityHandler struct {
	jwtConfig *util.JWTConfig
}

// NewSecurityHandler は新しいSecurityHandlerインスタンスを作成します。
func NewSecurityHandler() *SecurityHandler {
	return &SecurityHandler{
		jwtConfig: util.NewJWTConfig(),
	}
}

// HandleBearerAuth はBearer認証トークンを検証します。
// JWT認証の実装はこのメソッドで行います。
func (s *SecurityHandler) HandleBearerAuth(ctx context.Context, operationName api.OperationName, t api.BearerAuth) (context.Context, error) {
	// トークンの存在確認
	if t.Token == "" {
		return ctx, ErrMissingToken
	}

	// JWTトークンを検証
	claims, err := s.jwtConfig.ValidateToken(t.Token)
	if err != nil {
		return ctx, err
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
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}
