package handler

import (
	"context"

	"backend/api"
	"backend/security"

	"github.com/google/uuid"
)

// FriendsGet implements GET /friends operation.
// 自分のフレンド（フォロー）一覧取得
func (h *Handler) FriendsGet(ctx context.Context) (api.FriendsGetRes, error) {
	// 認証ユーザーIDを取得
	currentUserID, err := security.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	userID, err := uuid.Parse(currentUserID)
	if err != nil {
		return nil, ErrBadRequest
	}

	// ユーザーの存在確認
	user, err := h.client.User.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	// フォローしているユーザーのIDリストを取得
	followingIDs, err := user.QueryFollowing().IDs(ctx)
	if err != nil {
		return nil, err
	}

	result := api.FriendsGetOKApplicationJSON(followingIDs)
	return &result, nil
}

// FriendsPost implements POST /friends operation.
// フレンド追加（フォロー）
func (h *Handler) FriendsPost(ctx context.Context, req *api.FriendsPostReq) (api.FriendsPostRes, error) {
	// 認証ユーザーIDを取得
	currentUserID, err := security.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	userID, err := uuid.Parse(currentUserID)
	if err != nil {
		return nil, ErrBadRequest
	}

	// フォロー対象のユーザーIDを取得
	targetUserID := req.GetUserID()

	// 自分自身をフォローしようとしている場合はエラー
	if userID == targetUserID {
		return nil, ErrBadRequest
	}

	// フォロー対象のユーザーが存在するか確認
	_, err = h.client.User.Get(ctx, targetUserID)
	if err != nil {
		return nil, err
	}

	// ユーザーを取得してフォローを追加
	err = h.client.User.UpdateOneID(userID).
		AddFollowingIDs(targetUserID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &api.FriendsPostCreated{}, nil
}

// FriendsUserIDDelete implements DELETE /friends/{user_id} operation.
// フレンド削除（フォロー解除）
func (h *Handler) FriendsUserIDDelete(ctx context.Context, params api.FriendsUserIDDeleteParams) (api.FriendsUserIDDeleteRes, error) {
	// 認証ユーザーIDを取得
	currentUserID, err := security.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	userID, err := uuid.Parse(currentUserID)
	if err != nil {
		return nil, ErrBadRequest
	}

	// フォロー解除対象のユーザーID
	targetUserID := params.UserID

	// フォローを削除
	err = h.client.User.UpdateOneID(userID).
		RemoveFollowingIDs(targetUserID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &api.FriendsUserIDDeleteNoContent{}, nil
}

// UsersUserIDFriendsGet implements GET /users/{user_id}/friends operation.
// ユーザーのフレンド（フォロー）一覧取得
func (h *Handler) UsersUserIDFriendsGet(ctx context.Context, params api.UsersUserIDFriendsGetParams) (api.UsersUserIDFriendsGetRes, error) {
	// 認証確認（認証済みであれば誰でも他ユーザーのフォロー一覧を見られる）
	if _, err := security.GetUserIDFromContext(ctx); err != nil {
		return nil, ErrUnauthorized
	}

	// 対象ユーザーIDを取得
	targetUserID := params.UserID

	// ユーザーの存在確認
	user, err := h.client.User.Get(ctx, targetUserID)
	if err != nil {
		return nil, err
	}

	// フォローしているユーザーのIDリストを取得
	followingIDs, err := user.QueryFollowing().IDs(ctx)
	if err != nil {
		return nil, err
	}

	result := api.UsersUserIDFriendsGetOKApplicationJSON(followingIDs)
	return &result, nil
}
