package image

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

// アップロードされた画像をGoogle cloud storageに保存し、画像メタデータをデータベースに記録する処理
// 画像の取得処理もここに含まれる

// ImageService は画像の保存・取得を管理する構造体
type ImageService struct {
	bucket *storage.BucketHandle
}

// NewImageService は新しいImageServiceインスタンスを作成します
func NewImageService(ctx context.Context, bucketName string) (*ImageService, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}
	bucket := client.Bucket(bucketName)
	return &ImageService{
		bucket: bucket,
	}, nil
}

// UploadImage は画像をGCSにアップロードし、オブジェクト名とコンテンツタイプを返します
func (s *ImageService) UploadImage(ctx context.Context, file io.ReadSeeker, filename string, header map[string][]string) (objectName string, contentType string, err error) {
	// ファイルのMIMEタイプを検出
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", "", fmt.Errorf("failed to read file for content type detection: %w", err)
	}
	contentType = http.DetectContentType(buffer[:n])

	// 画像ファイルかどうかチェック
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" && contentType != "image/webp" {
		return "", "", fmt.Errorf("unsupported file type: %s", contentType)
	}

	// ファイルの読み込み位置をリセット
	if _, err := file.Seek(0, 0); err != nil {
		return "", "", fmt.Errorf("failed to reset file position: %w", err)
	}

	// UUIDでオブジェクト名を生成
	imageID := uuid.New()
	ext := path.Ext(filename)
	objectName = fmt.Sprintf("images/%s%s", imageID.String(), ext)

	// GCSにアップロード
	obj := s.bucket.Object(objectName)
	writer := obj.NewWriter(ctx)
	writer.ContentType = contentType

	if _, err := io.Copy(writer, file); err != nil {
		writer.Close()
		return "", "", fmt.Errorf("failed to upload to GCS: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", "", fmt.Errorf("failed to close GCS writer: %w", err)
	}

	return objectName, contentType, nil
}

// GetImage はGCSから画像を取得します
func (s *ImageService) GetImage(ctx context.Context, objectName string) (io.ReadCloser, string, error) {
	obj := s.bucket.Object(objectName)

	// オブジェクトの属性を取得してContent-Typeを確認
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return nil, "", fmt.Errorf("image not found: %w", err)
		}
		return nil, "", fmt.Errorf("failed to get object attributes: %w", err)
	}

	// オブジェクトを読み込む
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create reader: %w", err)
	}

	return reader, attrs.ContentType, nil
}

// DeleteImage はGCSから画像を削除します
func (s *ImageService) DeleteImage(ctx context.Context, objectName string) error {
	obj := s.bucket.Object(objectName)
	if err := obj.Delete(ctx); err != nil {
		if err == storage.ErrObjectNotExist {
			return fmt.Errorf("image not found: %w", err)
		}
		return fmt.Errorf("failed to delete image: %w", err)
	}
	return nil
}
