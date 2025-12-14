package handler

import (
	"context"

	"backend/api"
	"backend/ent"
	"backend/internal/image"
	"backend/internal/jwt"
)

// Handler は api.Handler インターフェースを実装するメイン構造体です。
// 各ドメインごとのハンドラーを保持し、メソッド呼び出しを委譲します。
type Handler struct {
	client       *ent.Client
	jwtHandler   *jwt.JwtHandler
	imageService *image.ImageService
}

// NewHandler は新しいHandlerインスタンスを作成します。
// 各ドメインハンドラーの初期化が必要な場合は、ここで行います。
func NewHandler(client *ent.Client, jwtHandler *jwt.JwtHandler, imageService *image.ImageService) (*Handler, error) {
	if client == nil {
		return nil, ErrClientRequired
	}
	if jwtHandler == nil {
		return nil, ErrJWTHandlerRequired
	}
	if imageService == nil {
		return nil, ErrImageServiceRequired
	}

	h := &Handler{
		client:       client,
		jwtHandler:   jwtHandler,
		imageService: imageService,
	}

	return h, nil
}

// NewError は api.Handler インターフェースで必要なエラーハンドラーを実装します。
func (h *Handler) NewError(ctx context.Context, err error) *api.GeneralErrorStatusCode {
	return NewError(ctx, err)
}
