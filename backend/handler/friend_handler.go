package handler

import (
	"context"

	"backend/api"
)

// FriendsGet implements GET /friends operation.
// 自分のフレンド（フォロー）一覧取得
func (h *Handler) FriendsGet(ctx context.Context) (api.FriendsGetRes, error) {
	// TODO: APIの処理を実装
	return &api.FriendsGetOKApplicationJSON{}, nil
}

// FriendsPost implements POST /friends operation.
// フレンド追加（フォロー）
func (h *Handler) FriendsPost(ctx context.Context, req *api.FriendsPostReq) (api.FriendsPostRes, error) {
	// TODO: APIの処理を実装
	return &api.FriendsPostCreated{}, nil
}

// FriendsUserIDDelete implements DELETE /friends/{user_id} operation.
// フレンド削除（フォロー解除）
func (h *Handler) FriendsUserIDDelete(ctx context.Context, params api.FriendsUserIDDeleteParams) (api.FriendsUserIDDeleteRes, error) {
	// TODO: APIの処理を実装
	return &api.FriendsUserIDDeleteNoContent{}, nil
}

// UsersUserIDFriendsGet implements GET /users/{user_id}/friends operation.
// ユーザーのフレンド（フォロー）一覧取得
func (h *Handler) UsersUserIDFriendsGet(ctx context.Context, params api.UsersUserIDFriendsGetParams) (api.UsersUserIDFriendsGetRes, error) {
	// TODO: APIの処理を実装
	return &api.UsersUserIDFriendsGetOKApplicationJSON{}, nil
}
