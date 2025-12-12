package handler

import (
	"context"

	"backend/api"
)

// ImagesImageIDGet implements GET /images/{image_id} operation.
// 画像取得
func (h *Handler) ImagesImageIDGet(ctx context.Context, params api.ImagesImageIDGetParams) (api.ImagesImageIDGetRes, error) {
	// TODO: APIの処理を実装
	return &api.ImagesImageIDGetOK{}, nil
}

// ImagesPost implements POST /images operation.
// 画像をアップロードしてIDを発行
func (h *Handler) ImagesPost(ctx context.Context, req api.OptImagesPostReq) (api.ImagesPostRes, error) {
	// TODO: APIの処理を実装
	return &api.Image{}, nil
}
