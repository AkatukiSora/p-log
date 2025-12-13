package handler

import (
	"context"

	"backend/api"
)

// PostsPostIDReactionsDelete implements DELETE /posts/{post_id}/reactions operation.
// 現在のユーザー自身のリアクションを削除
func (h *Handler) PostsPostIDReactionsDelete(ctx context.Context, params api.PostsPostIDReactionsDeleteParams) (api.PostsPostIDReactionsDeleteRes, error) {
	// TODO: APIの処理を実装
	return &api.PostsPostIDReactionsDeleteNoContent{}, nil
}

// PostsPostIDReactionsGet implements GET /posts/{post_id}/reactions operation.
// 投稿にリアクションしたユーザーの一覧取得
func (h *Handler) PostsPostIDReactionsGet(ctx context.Context, params api.PostsPostIDReactionsGetParams) (api.PostsPostIDReactionsGetRes, error) {
	// TODO: APIの処理を実装
	return &api.PostsPostIDReactionsGetOKApplicationJSON{}, nil
}

// PostsPostIDReactionsPost implements POST /posts/{post_id}/reactions operation.
// 投稿にリアクション（いいね）を追加
func (h *Handler) PostsPostIDReactionsPost(ctx context.Context, params api.PostsPostIDReactionsPostParams) (api.PostsPostIDReactionsPostRes, error) {
	// TODO: APIの処理を実装
	return &api.PostsPostIDReactionsPostCreated{}, nil
}
