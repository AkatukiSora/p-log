package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"backend/api"
	"backend/ent"
)

var (
	// ErrClientRequired はent.Clientが必須であることを示すエラーです。
	ErrClientRequired = errors.New("ent client is required")

	// ErrJWTHandlerRequired はJWTHandlerが必須であることを示すエラーです。
	ErrJWTHandlerRequired = errors.New("JWT handler is required")

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

	// エラーの種類に応じて適切なレスポンスを返す
	switch {
	case errors.Is(err, ErrNotFound) || ent.IsNotFound(err):
		return &api.GeneralErrorStatusCode{
			StatusCode: http.StatusNotFound,
			Response: api.Error{
				Message: "resource not found",
			},
		}
	case errors.Is(err, ErrUnauthorized):
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
