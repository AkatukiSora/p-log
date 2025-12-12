package handler

import (
	"context"

	"backend/api"
)

// PostsGet implements GET /posts operation.
// 現在のユーザーの投稿一覧取得
func (h *Handler) PostsGet(ctx context.Context, params api.PostsGetParams) (api.PostsGetRes, error) {
	// TODO: APIの処理を実装
	return &api.PostsGetOKApplicationJSON{}, nil
}

// PostsPost implements POST /posts operation.
// 進捗投稿作成
func (h *Handler) PostsPost(ctx context.Context, req *api.PostRequest) (api.PostsPostRes, error) {
	// TODO: APIの処理を実装
	return &api.Post{}, nil
}

// PostsPostIDDelete implements DELETE /posts/{post_id} operation.
// 投稿を削除（紐づいている画像も同時に削除）
func (h *Handler) PostsPostIDDelete(ctx context.Context, params api.PostsPostIDDeleteParams) (api.PostsPostIDDeleteRes, error) {
	// TODO: APIの処理を実装
	return &api.PostsPostIDDeleteNoContent{}, nil
}

// PostsPostIDGet implements GET /posts/{post_id} operation.
// 投稿詳細取得
func (h *Handler) PostsPostIDGet(ctx context.Context, params api.PostsPostIDGetParams) (api.PostsPostIDGetRes, error) {
	// TODO: APIの処理を実装
	return &api.Post{}, nil
}

// PostsPostIDPut implements PUT /posts/{post_id} operation.
// 投稿更新
func (h *Handler) PostsPostIDPut(ctx context.Context, req *api.PostRequest, params api.PostsPostIDPutParams) (api.PostsPostIDPutRes, error) {
	// TODO: APIの処理を実装
	return &api.Post{}, nil
}

// UsersUserIDPostsGet implements GET /users/{user_id}/posts operation.
// 指定ユーザーの投稿一覧取得
func (h *Handler) UsersUserIDPostsGet(ctx context.Context, params api.UsersUserIDPostsGetParams) (api.UsersUserIDPostsGetRes, error) {
	// TODO: APIの処理を実装
	return &api.UsersUserIDPostsGetOKApplicationJSON{}, nil
}
