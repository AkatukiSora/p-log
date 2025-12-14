package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"backend/api"
	"backend/ent"
	"backend/internal/jwt"
	"backend/security"
	"github.com/ogen-go/ogen/ogenerrors"
)

var (
	// ErrClientRequired はent.Clientが必須であることを示すエラーです。
	ErrClientRequired = errors.New("ent client is required")

	// ErrJWTHandlerRequired はJWTHandlerが必須であることを示すエラーです。
	ErrJWTHandlerRequired = errors.New("JWT handler is required")

	// ErrImageServiceRequired はImageServiceが必須であることを示すエラーです。
	ErrImageServiceRequired = errors.New("image service is required")

	// ErrNotFound はリソースが見つからない場合のエラーです。
	ErrNotFound = errors.New("resource not found")

	// ErrUnauthorized は認証が必要な場合のエラーです。
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden は権限がない場合のエラーです。
	ErrForbidden = errors.New("forbidden")

	// ErrBadRequest は不正なリクエストの場合のエラーです。
	ErrBadRequest = errors.New("bad request")

	// ErrConflict はリソースの競合が発生した場合のエラーです。
	ErrConflict = errors.New("conflict")

	// ErrInternalServer は内部サーバーエラーです。
	ErrInternalServer = errors.New("internal server error")
)

// NewError は発生したエラーを適切なHTTPステータスコードとメッセージに変換します。
func NewError(ctx context.Context, err error) *api.GeneralErrorStatusCode {
	if err == nil {
		return &api.GeneralErrorStatusCode{
			StatusCode: http.StatusInternalServerError,
			Response: api.Error{
				Message: "unknown error",
			},
		}
	}

	// エラーログを記録
	slog.ErrorContext(ctx, "handler error occurred",
		"error", err.Error(),
	)

	// ogenerrors.SecurityErrorの処理
	var secErr *ogenerrors.SecurityError
	if errors.As(err, &secErr) {
		// SecurityErrorの中身を確認
		if secErr.Err != nil {
			errMsg := secErr.Err.Error()
			// "security requirement is not satisfied"を含む場合は401
			if strings.Contains(errMsg, "security requirement is not satisfied") {
				return &api.GeneralErrorStatusCode{
					StatusCode: http.StatusUnauthorized,
					Response: api.Error{
						Message: "authentication required",
					},
				}
			}
			// JWT関連のエラーをチェック
			if errors.Is(secErr.Err, jwt.ErrInvalidToken) ||
				errors.Is(secErr.Err, jwt.ErrExpiredToken) ||
				errors.Is(secErr.Err, jwt.ErrInvalidClaims) ||
				errors.Is(secErr.Err, jwt.ErrRevokedToken) ||
				errors.Is(secErr.Err, security.ErrMissingToken) {
				return &api.GeneralErrorStatusCode{
					StatusCode: http.StatusUnauthorized,
					Response: api.Error{
						Message: "authentication required",
					},
				}
			}
		}
	}

	// エラーの種類に応じて適切なレスポンスを返す
	switch {
	case errors.Is(err, ErrNotFound) || ent.IsNotFound(err):
		return &api.GeneralErrorStatusCode{
			StatusCode: http.StatusNotFound,
			Response: api.Error{
				Message: "resource not found",
			},
		}
	case errors.Is(err, ErrUnauthorized) ||
		errors.Is(err, jwt.ErrInvalidToken) ||
		errors.Is(err, jwt.ErrExpiredToken) ||
		errors.Is(err, jwt.ErrInvalidClaims) ||
		errors.Is(err, jwt.ErrRevokedToken) ||
		errors.Is(err, security.ErrMissingToken):
		return &api.GeneralErrorStatusCode{
			StatusCode: http.StatusUnauthorized,
			Response: api.Error{
				Message: "authentication required",
			},
		}
	case errors.Is(err, ErrForbidden):
		return &api.GeneralErrorStatusCode{
			StatusCode: http.StatusForbidden,
			Response: api.Error{
				Message: "access forbidden",
			},
		}
	case errors.Is(err, ErrBadRequest) || ent.IsValidationError(err):
		return &api.GeneralErrorStatusCode{
			StatusCode: http.StatusBadRequest,
			Response: api.Error{
				Message: err.Error(),
			},
		}
	case errors.Is(err, ErrConflict) || ent.IsConstraintError(err):
		return &api.GeneralErrorStatusCode{
			StatusCode: http.StatusConflict,
			Response: api.Error{
				Message: "resource conflict occurred",
			},
		}
	default:
		// その他のエラーは500 Internal Server Errorとして扱う
		return &api.GeneralErrorStatusCode{
			StatusCode: http.StatusInternalServerError,
			Response: api.Error{
				Message: "internal server error",
			},
		}
	}
}
