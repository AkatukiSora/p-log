package handler

import (
"backend/api"
"context"
"fmt"
)

// AuthCallbackGet implements GET /auth/callback operation.
// OIDCコールバック（コードをトークンに交換）
func (h *Handler) AuthCallbackGet(ctx context.Context, params api.AuthCallbackGetParams) (api.AuthCallbackGetRes, error) {
// TODO: APIの処理を実装
// モックのトークンを生成
accessToken := "mock_access_token"
refreshToken := "mock_refresh_token"

// Set-Cookieヘッダーの値を構築
// access_token: HttpOnly, Secure, SameSite=Lax
// refresh_token: HttpOnly, Secure, SameSite=Lax, Path=/api/v1/auth/refresh
setCookieHeader := fmt.Sprintf(
"access_token=%s; HttpOnly; Secure; SameSite=Lax; Path=/; Max-Age=3600, "+
"refresh_token=%s; HttpOnly; Secure; SameSite=Lax; Path=/api/v1/auth/refresh; Max-Age=604800",
accessToken, refreshToken,
)

return &api.AuthCallbackGetNoContent{
SetCookie: api.NewOptString(setCookieHeader),
}, nil
}

// AuthRefreshPost implements POST /auth/refresh operation.
// アクセストークンのリフレッシュ
func (h *Handler) AuthRefreshPost(ctx context.Context) (api.AuthRefreshPostRes, error) {
// TODO: APIの処理を実装
// モックの新しいトークンを生成
newAccessToken := "new_mock_access_token"
newRefreshToken := "new_mock_refresh_token"

// Set-Cookieヘッダーの値を構築（上書き）
// access_token: HttpOnly, Secure, SameSite=Lax
// refresh_token: HttpOnly, Secure, SameSite=Lax, Path=/api/v1/auth/refresh
setCookieHeader := fmt.Sprintf(
"access_token=%s; HttpOnly; Secure; SameSite=Lax; Path=/; Max-Age=3600, "+
"refresh_token=%s; HttpOnly; Secure; SameSite=Lax; Path=/api/v1/auth/refresh; Max-Age=604800",
newAccessToken, newRefreshToken,
)

return &api.AuthRefreshPostNoContent{
SetCookie: api.NewOptString(setCookieHeader),
}, nil
}

// AuthLoginGet implements GET /auth/login operation.
// OIDCログインを開始
func (h *Handler) AuthLoginGet(ctx context.Context) error {
// TODO: APIの処理を実装
return nil
}

// AuthLogoutPost implements POST /auth/logout operation.
// ログアウト
func (h *Handler) AuthLogoutPost(ctx context.Context) error {
// TODO: APIの処理を実装
return nil
}

// AuthMeGet implements GET /auth/me operation.
// 現在のユーザー情報を取得
func (h *Handler) AuthMeGet(ctx context.Context) (api.AuthMeGetRes, error) {
// TODO: APIの処理を実装
return &api.User{}, nil
}
