package handler

import (
	"context"

	"backend/api"
	"backend/security"
)

// UsersPost implements POST /users operation.
// UsersUserIDDelete implements DELETE /users/{user_id} operation.
// ユーザーアカウント削除
func (h *Handler) UsersUserIDDelete(ctx context.Context, params api.UsersUserIDDeleteParams) (api.UsersUserIDDeleteRes, error) {
	// TODO: APIの処理を実装
	err := h.client.User.DeleteOneID(params.UserID).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &api.UsersUserIDDeleteNoContent{}, nil
}

// UsersUserIDGet implements GET /users/{user_id} operation.
// ユーザープロフィール取得
func (h *Handler) UsersUserIDGet(ctx context.Context, params api.UsersUserIDGetParams) (api.UsersUserIDGetRes, error) {
	if _, err := security.GetUserIDFromContext(ctx); err != nil {
		return nil, ErrUnauthorized
	}
	userID := params.UserID

	user, err := h.client.User.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	genreIDs, err := user.QueryGenres().IDs(ctx)
	if err != nil {
		return nil, err
	}

	var birthday api.OptDate
	if user.Birthday != nil {
		birthday = api.NewOptDate(*user.Birthday)
	}

	var hometown api.OptString
	if user.Hometown != nil {
		hometown = api.NewOptString(*user.Hometown)
	}

	var bio api.OptString
	if user.Bio != nil {
		bio = api.NewOptString(*user.Bio)
	}
	res := api.User{
		ID:       user.ID,
		Name:     user.Name,
		Birthday: birthday,
		Genres:   genreIDs,
		Hometown: hometown,
		Bio:      bio,
	}

	return &res, nil
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
