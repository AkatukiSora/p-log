package handler

import (
	"backend/api"
	"backend/ent"
	"backend/ent/user"
	"backend/internal/oauth"
	"context"
	"fmt"
	"sync"
)

var (
	oauthConfig     *oauth.OAuthConfig
	oauthConfigOnce sync.Once
)

// getOAuthConfig はOAuthConfigのシングルトンインスタンスを返します
func getOAuthConfig() *oauth.OAuthConfig {
	oauthConfigOnce.Do(func() {
		oauthConfig = oauth.NewOAuthConfig()
	})
	return oauthConfig
}

// AuthCallbackGet implements GET /auth/callback operation.
// OIDCコールバック（コードをトークンに交換）
func (h *Handler) AuthCallbackGet(ctx context.Context, params api.AuthCallbackGetParams) (api.AuthCallbackGetRes, error) {
	oauthCfg := getOAuthConfig()

	// stateパラメータの検証
	if !oauthCfg.ValidateState(params.State) {
		return &api.Error{
			Message: "Invalid state parameter",
		}, nil
	}

	// 認可コードをトークンに交換
	token, err := oauthCfg.ExchangeCode(ctx, params.Code)
	if err != nil {
		return &api.Error{
			Message: fmt.Sprintf("Failed to exchange code: %v", err),
		}, nil
	}

	// Googleからユーザー情報を取得
	userInfo, err := oauthCfg.GetUserInfo(ctx, token)
	if err != nil {
		return &api.Error{
			Message: fmt.Sprintf("Failed to get user info: %v", err),
		}, nil
	}

	// メールアドレスでユーザーを検索
	existingUser, err := h.client.User.Query().
		Where(user.EmailEQ(userInfo.Email)).
		Only(ctx)

	var currentUser *ent.User
	if err != nil {
		if ent.IsNotFound(err) {
			// ユーザーが存在しない場合は新規作成
			currentUser, err = h.client.User.Create().
				SetEmail(userInfo.Email).
				SetName(userInfo.Name).
				Save(ctx)
			if err != nil {
				return &api.Error{
					Message: fmt.Sprintf("Failed to create user: %v", err),
				}, nil
			}
		} else {
			return &api.Error{
				Message: fmt.Sprintf("Database error: %v", err),
			}, nil
		}
	} else {
		currentUser = existingUser
	}

	// JWTトークンを生成
	accessToken, refreshToken, err := h.jwtHandler.GenerateTokens(currentUser.ID, currentUser.Email, ctx)
	if err != nil {
		return &api.Error{
			Message: fmt.Sprintf("Failed to generate tokens: %v", err),
		}, nil
	}

	// Set-Cookieヘッダーの値を構築
	// 開発環境用: Secure属性なし（HTTPでも動作）
	// access_token: HttpOnly, SameSite=Lax
	// refresh_token: HttpOnly, SameSite=Lax, Path=/api/v1/auth/refresh
	setCookieHeader := fmt.Sprintf(
		"access_token=%s; HttpOnly; SameSite=Lax; Path=/; Max-Age=3600, "+
			"refresh_token=%s; HttpOnly; SameSite=Lax; Path=/api/v1/auth/refresh; Max-Age=604800",
		accessToken, refreshToken,
	)

	return &api.AuthCallbackGetNoContent{
		SetCookie: api.NewOptString(setCookieHeader),
	}, nil
}

// AuthRefreshPost implements POST /auth/refresh operation.
// アクセストークンのリフレッシュ
func (h *Handler) AuthRefreshPost(ctx context.Context) (api.AuthRefreshPostRes, error) {
	// TODO: Cookieからリフレッシュトークンを取得する実装が必要
	// 現在はモック実装
	// モックの新しいトークンを生成
	newAccessToken := "new_mock_access_token"
	newRefreshToken := "new_mock_refresh_token"

	// Set-Cookieヘッダーの値を構築（上書き）
	// 開発環境用: Secure属性なし（HTTPでも動作）
	// access_token: HttpOnly, SameSite=Lax
	// refresh_token: HttpOnly, SameSite=Lax, Path=/api/v1/auth/refresh
	setCookieHeader := fmt.Sprintf(
		"access_token=%s; HttpOnly; SameSite=Lax; Path=/; Max-Age=3600, "+
			"refresh_token=%s; HttpOnly; SameSite=Lax; Path=/api/v1/auth/refresh; Max-Age=604800",
		newAccessToken, newRefreshToken,
	)

	return &api.AuthRefreshPostNoContent{
		SetCookie: api.NewOptString(setCookieHeader),
	}, nil
}

// AuthLoginGet implements GET /auth/login operation.
// OIDCログインを開始
func (h *Handler) AuthLoginGet(ctx context.Context) (*api.AuthLoginGetMovedPermanently, error) {
	oauthCfg := getOAuthConfig()

	// Google OAuth2の認証URLを生成してリダイレクト
	authURL := oauthCfg.GetAuthURL()

	// リダイレクトレスポンスを返す
	return &api.AuthLoginGetMovedPermanently{
		Location: api.NewOptString(authURL),
	}, nil
}

// AuthLogoutPost implements POST /auth/logout operation.
// ログアウト
func (h *Handler) AuthLogoutPost(ctx context.Context) error {
	// TODO: Cookieからリフレッシュトークンを取得して無効化
	// TODO: Cookieを削除する処理を実装
	return nil
}

// AuthMeGet implements GET /auth/me operation.
// 現在のユーザー情報を取得
func (h *Handler) AuthMeGet(ctx context.Context) (api.AuthMeGetRes, error) {
	// TODO: Cookieからアクセストークンを取得して検証
	// TODO: トークンからユーザーIDを取得してユーザー情報を返す
	return &api.User{}, nil
}
