package handler

import (
	"context"

	"backend/api"
)

// TimelineGet implements GET /timeline operation.
// タイムライン取得（投稿は新しい順で返される）
func (h *Handler) TimelineGet(ctx context.Context, params api.TimelineGetParams) (api.TimelineGetRes, error) {
	// TODO: APIの処理を実装
	return &api.TimelineGetOKApplicationJSON{}, nil
}
