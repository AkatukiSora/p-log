package handler

import (
	"context"

	"backend/api"
	"backend/ent"
	"backend/util"
)

// Handler は api.Handler インターフェースを実装するメイン構造体です。
// 各ドメインごとのハンドラーを保持し、メソッド呼び出しを委譲します。
type Handler struct {
	client    *ent.Client
	jwtConfig *util.JWTConfig
}

// NewHandler は新しいHandlerインスタンスを作成します。
// 各ドメインハンドラーの初期化が必要な場合は、ここで行います。
func NewHandler(client *ent.Client, jwtConfig *util.JWTConfig) (*Handler, error) {
	if client == nil {
		return nil, ErrClientRequired
	}
	if jwtConfig == nil {
		return nil, ErrJWTConfigRequired
	}

	h := &Handler{
		client:    client,
		jwtConfig: jwtConfig,
	}

	return h, nil
}

// NewError は api.Handler インターフェースで必要なエラーハンドラーを実装します。
func (h *Handler) NewError(ctx context.Context, err error) *api.GeneralErrorStatusCode {
	return NewError(ctx, err)
}
