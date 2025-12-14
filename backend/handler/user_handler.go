package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"

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
	_, err = update.Save(ctx)
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
	updatedUser, err := h.client.User.Query().
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

	// ユーザーを取得
	u, err := h.client.User.Get(ctx, params.UserID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// プロフィール画像がない場合は404
	if u.ProfilePictureID == nil {
		return nil, ErrNotFound
	}

	// 画像メタデータを取得
	img, err := h.client.Image.Get(ctx, *u.ProfilePictureID)
	if err != nil {
		if ent.IsNotFound(err) {
			// DBに画像が見つからない場合は、profile_picture_idだけクリア
			_, err = h.client.User.UpdateOneID(params.UserID).
				ClearProfilePictureID().
				Save(ctx)
			return &api.UsersUserIDIconDeleteNoContent{}, err
		}
		return nil, err
	}

	// GCSから画像を削除
	if err := h.imageService.DeleteImage(ctx, img.ObjectName); err != nil {
		// GCS削除に失敗してもログに記録して続行
		fmt.Printf("failed to delete image from GCS: %v\n", err)
	}

	// DBから画像メタデータを削除
	if err := h.client.Image.DeleteOneID(*u.ProfilePictureID).Exec(ctx); err != nil {
		return nil, fmt.Errorf("failed to delete image metadata: %w", err)
	}

	// ユーザーのprofile_picture_idをクリア
	_, err = h.client.User.UpdateOneID(params.UserID).
		ClearProfilePictureID().
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &api.UsersUserIDIconDeleteNoContent{}, nil
}

// UsersUserIDIconGet implements GET /users/{user_id}/icon operation.
// ユーザーアイコン画像取得
func (h *Handler) UsersUserIDIconGet(ctx context.Context, params api.UsersUserIDIconGetParams) (api.UsersUserIDIconGetRes, error) {
	// ユーザーを取得
	u, err := h.client.User.Get(ctx, params.UserID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// プロフィール画像がない場合は404
	if u.ProfilePictureID == nil {
		return nil, ErrNotFound
	}

	// 画像メタデータを取得
	img, err := h.client.Image.Get(ctx, *u.ProfilePictureID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// GCSから画像を取得
	reader, _, err := h.imageService.GetImage(ctx, img.ObjectName)
	if err != nil {
		return nil, fmt.Errorf("failed to get image from GCS: %w", err)
	}

	return &api.UsersUserIDIconGetOK{
		Data: reader,
	}, nil
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

	// リクエストボディが存在するかチェック
	if !req.Set {
		return nil, ErrBadRequest
	}

	// アイコンファイルが存在するかチェック
	iconOpt := req.Value.GetIcon()
	if !iconOpt.Set {
		return nil, ErrBadRequest
	}

	multipartFile := iconOpt.Value

	// ユーザーを取得（既存のprofile_picture_idを確認）
	u, err := h.client.User.Get(ctx, params.UserID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// ファイルの内容を読み込む（UploadImageでSeekが必要なため、バッファに読み込む）
	fileData, err := io.ReadAll(multipartFile.File)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	fileReader := bytes.NewReader(fileData)

	// GCSに画像をアップロード
	objectName, contentType, err := h.imageService.UploadImage(ctx, fileReader, multipartFile.Name, multipartFile.Header)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image: %w", err)
	}

	// DBに画像メタデータを保存
	img, err := h.client.Image.Create().
		SetObjectName(objectName).
		SetContentType(contentType).
		SetUploadedByID(userID).
		Save(ctx)
	if err != nil {
		// アップロードに失敗した場合、GCSから画像を削除
		if delErr := h.imageService.DeleteImage(ctx, objectName); delErr != nil {
			fmt.Printf("failed to cleanup uploaded image: %v\n", delErr)
		}
		return nil, fmt.Errorf("failed to save image metadata: %w", err)
	}

	// 既存のプロフィール画像がある場合は削除
	if u.ProfilePictureID != nil {
		// 古い画像メタデータを取得
		oldImg, err := h.client.Image.Get(ctx, *u.ProfilePictureID)
		if err == nil {
			// GCSから古い画像を削除
			if delErr := h.imageService.DeleteImage(ctx, oldImg.ObjectName); delErr != nil {
				fmt.Printf("failed to delete old image from GCS: %v\n", delErr)
			}
			// DBから古い画像メタデータを削除
			if delErr := h.client.Image.DeleteOneID(*u.ProfilePictureID).Exec(ctx); delErr != nil {
				fmt.Printf("failed to delete old image metadata: %v\n", delErr)
			}
		}
	}

	// ユーザーのprofile_picture_idを更新
	_, err = h.client.User.UpdateOneID(params.UserID).
		SetProfilePictureID(img.ID).
		Save(ctx)
	if err != nil {
		// 更新に失敗した場合、新しくアップロードした画像を削除
		if delErr := h.imageService.DeleteImage(ctx, objectName); delErr != nil {
			fmt.Printf("failed to cleanup uploaded image: %v\n", delErr)
		}
		if delErr := h.client.Image.DeleteOneID(img.ID).Exec(ctx); delErr != nil {
			fmt.Printf("failed to cleanup image metadata: %v\n", delErr)
		}
		return nil, fmt.Errorf("failed to update user profile_picture_id: %w", err)
	}

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
