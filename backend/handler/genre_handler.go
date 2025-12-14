package handler

import (
	"context"

	"backend/api"
)

// GenresGet implements GET /genres operation.
// 利用可能なジャンル一覧取得
func (h *Handler) GenresGet(ctx context.Context) ([]api.Genre, error) {
	// データベースから全てのジャンルを取得
	genres, err := h.client.Genre.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	// ent.Genreをapi.Genreに変換
	result := make([]api.Genre, len(genres))
	for i, g := range genres {
		result[i] = api.Genre{
			ID:   g.ID,
			Name: g.Name,
		}
	}

	return result, nil
}
