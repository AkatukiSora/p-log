package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"backend/api"
	"backend/ent/user"
	"backend/security"
	"github.com/google/uuid"
)

// ImagesImageIDGet implements GET /images/{image_id} operation.
// 画像取得
func (h *Handler) ImagesImageIDGet(ctx context.Context, params api.ImagesImageIDGetParams) (api.ImagesImageIDGetRes, error) {
	// DBから画像メタデータを取得
	img, err := h.client.Image.Get(ctx, params.ImageID)
	if err != nil {
		return nil, ErrNotFound
	}

	// GCSから画像を取得
	reader, _, err := h.imageService.GetImage(ctx, img.ObjectName)
	if err != nil {
		return nil, fmt.Errorf("failed to get image from GCS: %w", err)
	}
	defer reader.Close()

	return &api.ImagesImageIDGetOK{
		Data: reader,
	}, nil
}

// ImagesPost implements POST /images operation.
// 画像をアップロードしてIDを発行
func (h *Handler) ImagesPost(ctx context.Context, req api.OptImagesPostReq) (api.ImagesPostRes, error) {
	// リクエストボディが存在するかチェック
	if !req.Set {
		return nil, ErrBadRequest
	}

	// 画像ファイルが存在するかチェック
	imageOpt := req.Value.GetImage()
	if !imageOpt.Set {
		return nil, ErrBadRequest
	}

	multipartFile := imageOpt.Value

	// contextからユーザーIDを取得
	userIDStr, err := security.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ErrUnauthorized
	}

	// ユーザーが存在するか確認
	exists, err := h.client.User.Query().Where(user.ID(userID)).Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return nil, ErrUnauthorized
	}

	// ファイルの内容を読み込む（UploadImageでSeekが必要なため、バッファに読み込む）
	fileData, err := io.ReadAll(multipartFile.File)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	fileReader := bytes.NewReader(fileData)

	// GCSに画像をアップロード（ファイル名とヘッダーを渡す）
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
			// ログに記録するが、元のエラーを返す
			fmt.Printf("failed to cleanup uploaded image: %v\n", delErr)
		}
		return nil, fmt.Errorf("failed to save image metadata: %w", err)
	}

	// 画像URLを生成（/api/v1/images/{image_id}形式）
	imageURL := fmt.Sprintf("/api/v1/images/%s", img.ID.String())

	return &api.Image{
		ID:  img.ID,
		URL: imageURL,
	}, nil
}
