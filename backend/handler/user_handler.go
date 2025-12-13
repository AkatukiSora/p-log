package handler

import (
	"context"

	"backend/api"
)

// UsersPost implements POST /users operation.
// 新規ユーザー登録
func (h *Handler) UsersPost(ctx context.Context, req *api.UserRequest) (api.UsersPostRes, error) {
	// TODO: APIの処理を実装
	return &api.User{}, nil
}

// UsersUserIDDelete implements DELETE /users/{user_id} operation.
// ユーザーアカウント削除
func (h *Handler) UsersUserIDDelete(ctx context.Context, params api.UsersUserIDDeleteParams) (api.UsersUserIDDeleteRes, error) {
	// TODO: APIの処理を実装
	return &api.UsersUserIDDeleteNoContent{}, nil
}

// UsersUserIDGet implements GET /users/{user_id} operation.
// ユーザープロフィール取得
func (h *Handler) UsersUserIDGet(ctx context.Context, params api.UsersUserIDGetParams) (api.UsersUserIDGetRes, error) {
	// TODO: APIの処理を実装
	return &api.User{}, nil
}

// UsersUserIDPut implements PUT /users/{user_id} operation.
// ユーザープロフィール更新
func (h *Handler) UsersUserIDPut(ctx context.Context, req *api.UserRequest, params api.UsersUserIDPutParams) (api.UsersUserIDPutRes, error) {
	// TODO: APIの処理を実装
	return &api.User{}, nil
}

// UsersUserIDIconDelete implements DELETE /users/{user_id}/icon operation.
// ユーザーアイコン削除
func (h *Handler) UsersUserIDIconDelete(ctx context.Context, params api.UsersUserIDIconDeleteParams) (api.UsersUserIDIconDeleteRes, error) {
	// TODO: APIの処理を実装
	return &api.UsersUserIDIconDeleteNoContent{}, nil
}

// UsersUserIDIconGet implements GET /users/{user_id}/icon operation.
// ユーザーアイコン画像取得
func (h *Handler) UsersUserIDIconGet(ctx context.Context, params api.UsersUserIDIconGetParams) (api.UsersUserIDIconGetRes, error) {
	// TODO: APIの処理を実装
	return &api.UsersUserIDIconGetOK{}, nil
}

// UsersUserIDIconPost implements POST /users/{user_id}/icon operation.
// ユーザーアイコンのアップロードまたは置換
func (h *Handler) UsersUserIDIconPost(ctx context.Context, req api.OptUsersUserIDIconPostReq, params api.UsersUserIDIconPostParams) (api.UsersUserIDIconPostRes, error) {
	// TODO: APIの処理を実装
	return &api.UsersUserIDIconPostNoContent{}, nil
}
