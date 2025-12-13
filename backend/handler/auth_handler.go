package handler

import (
	"backend/api"
	"context"
)

// AuthCallbackGet implements GET /auth/callback operation.
// OIDCコールバック（コードをトークンに交換）
func (h *Handler) AuthCallbackGet(ctx context.Context, params api.AuthCallbackGetParams) (api.AuthCallbackGetRes, error) {
	// TODO: APIの処理を実装
	return &api.AuthToken{}, nil
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
