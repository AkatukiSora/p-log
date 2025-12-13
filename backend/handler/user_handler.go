package handler

import (
	"context"

	"backend/api"
	"backend/ent"
	"backend/ent/genre"
	"backend/ent/user"
	"backend/security"

	"github.com/google/uuid"
)

// UsersUserIDDelete implements DELETE /users/{user_id} operation.
// ユーザーアカウント削除
func (h *Handler) UsersUserIDDelete(ctx context.Context, params api.UsersUserIDDeleteParams) (api.UsersUserIDDeleteRes, error) {
	// コンテキストからユーザーIDを取得
	userIDStr, err := security.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ErrBadRequest
	}

	// 自分自身のアカウントかチェック
	if userID != params.UserID {
		return nil, ErrForbidden
	}

	// ユーザーが存在するか確認
	exists, err := h.client.User.Query().Where(user.IDEQ(params.UserID)).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound
	}

	// ユーザーを削除
	err = h.client.User.DeleteOneID(params.UserID).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &api.UsersUserIDDeleteNoContent{}, nil
}

// UsersUserIDGet implements GET /users/{user_id} operation.
// ユーザープロフィール取得
func (h *Handler) UsersUserIDGet(ctx context.Context, params api.UsersUserIDGetParams) (api.UsersUserIDGetRes, error) {
	// 認証チェック
	if _, err := security.GetUserIDFromContext(ctx); err != nil {
		return nil, ErrUnauthorized
	}

	// ユーザーを取得
	user, err := h.client.User.Query().
		Where(user.IDEQ(params.UserID)).
		WithGenres().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// レスポンスを構築
	res, err := h.convertUserToAPI(user)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UsersUserIDPut implements PUT /users/{user_id} operation.
// ユーザープロフィール更新
func (h *Handler) UsersUserIDPut(ctx context.Context, req *api.UserRequest, params api.UsersUserIDPutParams) (api.UsersUserIDPutRes, error) {
	// コンテキストからユーザーIDを取得
	userIDStr, err := security.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ErrBadRequest
	}

	// 自分自身のプロフィールかチェック
	if userID != params.UserID {
		return nil, ErrForbidden
	}

	// ユーザーが存在するか確認
	exists, err := h.client.User.Query().Where(user.IDEQ(params.UserID)).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound
	}

	// 更新ビルダーを作成
	update := h.client.User.UpdateOneID(params.UserID).SetName(req.Name)

	// オプション型のフィールドを設定
	if req.Birthday.IsSet() {
		update.SetBirthday(req.Birthday.Value)
	}
	if req.Hometown.IsSet() {
		update.SetHometown(req.Hometown.Value)
	}
	if req.Bio.IsSet() {
		update.SetBio(req.Bio.Value)
	}

	// 更新を実行
	updatedUser, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	// ジャンルの更新（既存の関連を削除して新しく追加）
	if len(req.Genres) > 0 {
		// 既存のジャンル関連をクリア
		_, err = h.client.User.UpdateOneID(params.UserID).
			ClearGenres().
			Save(ctx)
		if err != nil {
			return nil, err
		}

		// ジャンルが存在するか確認
		existingGenres, err := h.client.Genre.Query().
			Where(genre.IDIn(req.Genres...)).
			IDs(ctx)
		if err != nil {
			return nil, err
		}

		// 存在しないジャンルがあればエラー
		if len(existingGenres) != len(req.Genres) {
			return nil, ErrBadRequest
		}

		// 新しいジャンルを追加
		_, err = h.client.User.UpdateOneID(params.UserID).
			AddGenreIDs(req.Genres...).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	// 更新後のユーザーを再取得（ジャンルを含む）
	updatedUser, err = h.client.User.Query().
		Where(user.IDEQ(params.UserID)).
		WithGenres().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// レスポンスを構築
	res, err := h.convertUserToAPI(updatedUser)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UsersUserIDIconDelete implements DELETE /users/{user_id}/icon operation.
// ユーザーアイコン削除
func (h *Handler) UsersUserIDIconDelete(ctx context.Context, params api.UsersUserIDIconDeleteParams) (api.UsersUserIDIconDeleteRes, error) {
	// コンテキストからユーザーIDを取得
	userIDStr, err := security.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ErrBadRequest
	}

	// 自分自身のアイコンかチェック
	if userID != params.UserID {
		return nil, ErrForbidden
	}

	// ユーザーが存在するか確認
	exists, err := h.client.User.Query().Where(user.IDEQ(params.UserID)).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound
	}

	// TODO: プロフィール画像を削除する処理を実装
	// 1. ユーザーのprofile_picture_idを取得
	// 2. 画像をストレージから削除
	// 3. profile_picture_idをnullに設定

	return &api.UsersUserIDIconDeleteNoContent{}, nil
}

// UsersUserIDIconGet implements GET /users/{user_id}/icon operation.
// ユーザーアイコン画像取得
func (h *Handler) UsersUserIDIconGet(ctx context.Context, params api.UsersUserIDIconGetParams) (api.UsersUserIDIconGetRes, error) {
	// ユーザーが存在するか確認
	exists, err := h.client.User.Query().Where(user.IDEQ(params.UserID)).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound
	}

	// TODO: プロフィール画像を取得する処理を実装
	// 1. ユーザーのprofile_picture_idを取得
	// 2. 画像をストレージから取得
	// 3. 画像バイナリを返す

	return &api.UsersUserIDIconGetOK{}, nil
}

// UsersUserIDIconPost implements POST /users/{user_id}/icon operation.
// ユーザーアイコンのアップロードまたは置換
func (h *Handler) UsersUserIDIconPost(ctx context.Context, req api.OptUsersUserIDIconPostReq, params api.UsersUserIDIconPostParams) (api.UsersUserIDIconPostRes, error) {
	// コンテキストからユーザーIDを取得
	userIDStr, err := security.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ErrBadRequest
	}

	// 自分自身のアイコンかチェック
	if userID != params.UserID {
		return nil, ErrForbidden
	}

	// ユーザーが存在するか確認
	exists, err := h.client.User.Query().Where(user.IDEQ(params.UserID)).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound
	}

	// TODO: プロフィール画像をアップロードする処理を実装
	// 1. リクエストから画像データを取得
	// 2. 画像をストレージにアップロード
	// 3. 既存のprofile_picture_idがあれば古い画像を削除
	// 4. ユーザーのprofile_picture_idを更新

	return &api.UsersUserIDIconPostNoContent{}, nil
}

// convertUserToAPI はent.Userをapi.Userに変換します
func (h *Handler) convertUserToAPI(u *ent.User) (*api.User, error) {
	// ジャンルIDを取得
	var genreIDs []uuid.UUID
	if genres, err := u.Edges.GenresOrErr(); err == nil {
		for _, g := range genres {
			genreIDs = append(genreIDs, g.ID)
		}
	}

	// オプション型のフィールドを設定
	var birthday api.OptDate
	if u.Birthday != nil {
		birthday = api.NewOptDate(*u.Birthday)
	}

	var hometown api.OptString
	if u.Hometown != nil {
		hometown = api.NewOptString(*u.Hometown)
	}

	var bio api.OptString
	if u.Bio != nil {
		bio = api.NewOptString(*u.Bio)
	}

	return &api.User{
		ID:       u.ID,
		Name:     u.Name,
		Birthday: birthday,
		Genres:   genreIDs,
		Hometown: hometown,
		Bio:      bio,
	}, nil
}
