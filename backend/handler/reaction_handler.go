package handler

import (
	"context"

	"backend/api"
	"backend/ent"
	"backend/ent/post"
	"backend/ent/reaction"
	"backend/ent/user"
	"backend/security"

	"github.com/google/uuid"
)

// PostsPostIDReactionsPost implements POST /posts/{post_id}/reactions operation.
// 投稿にリアクション（いいね）を追加
func (h *Handler) PostsPostIDReactionsPost(ctx context.Context, params api.PostsPostIDReactionsPostParams) (api.PostsPostIDReactionsPostRes, error) {
	// コンテキストからユーザーIDを取得
	userIDStr, err := security.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ErrBadRequest
	}

	// 投稿が存在するか確認
	exists, err := h.client.Post.Query().Where(post.IDEQ(params.PostID)).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound
	}

	// ユーザーが存在するか確認
	exists, err = h.client.User.Query().Where(user.IDEQ(userID)).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound
	}

	// 既にリアクションしているか確認
	exists, err = h.client.Reaction.Query().
		Where(
			reaction.HasUserWith(user.IDEQ(userID)),
			reaction.HasPostWith(post.IDEQ(params.PostID)),
		).
		Exist(ctx)
	if err != nil {
		return nil, err
	}

	// 既にリアクションしている場合はそのまま成功を返す（冪等性）
	if exists {
		return &api.PostsPostIDReactionsPostCreated{}, nil
	}

	// リアクションを作成
	_, err = h.client.Reaction.
		Create().
		SetUserID(userID).
		SetPostID(params.PostID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &api.PostsPostIDReactionsPostCreated{}, nil
}

// PostsPostIDReactionsDelete implements DELETE /posts/{post_id}/reactions operation.
// 現在のユーザー自身のリアクションを削除
func (h *Handler) PostsPostIDReactionsDelete(ctx context.Context, params api.PostsPostIDReactionsDeleteParams) (api.PostsPostIDReactionsDeleteRes, error) {
	// コンテキストからユーザーIDを取得
	userIDStr, err := security.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ErrBadRequest
	}

	// 投稿が存在するか確認
	exists, err := h.client.Post.Query().Where(post.IDEQ(params.PostID)).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound
	}

	// 自分のリアクションを取得
	r, err := h.client.Reaction.Query().
		Where(
			reaction.HasUserWith(user.IDEQ(userID)),
			reaction.HasPostWith(post.IDEQ(params.PostID)),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			// リアクションが存在しない場合も成功を返す（冪等性）
			return &api.PostsPostIDReactionsDeleteNoContent{}, nil
		}
		return nil, err
	}

	// リアクションを削除
	err = h.client.Reaction.DeleteOne(r).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &api.PostsPostIDReactionsDeleteNoContent{}, nil
}

// PostsPostIDReactionsGet implements GET /posts/{post_id}/reactions operation.
// 投稿にリアクションしたユーザーの一覧取得
func (h *Handler) PostsPostIDReactionsGet(ctx context.Context, params api.PostsPostIDReactionsGetParams) (api.PostsPostIDReactionsGetRes, error) {
	// 投稿が存在するか確認
	exists, err := h.client.Post.Query().Where(post.IDEQ(params.PostID)).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound
	}

	// 投稿のリアクション一覧を取得
	reactions, err := h.client.Reaction.Query().
		Where(reaction.HasPostWith(post.IDEQ(params.PostID))).
		WithUser().
		Order(ent.Desc(reaction.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	// entのReactionをapiのReactionに変換
	apiReactions := make([]api.Reaction, 0, len(reactions))
	for _, r := range reactions {
		reactionUser, err := r.Edges.UserOrErr()
		if err != nil {
			return nil, err
		}

		apiReaction := api.Reaction{
			UserID:    reactionUser.ID,
			CreatedAt: r.CreatedAt,
		}
		apiReactions = append(apiReactions, apiReaction)
	}

	response := api.PostsPostIDReactionsGetOKApplicationJSON(apiReactions)
	return &response, nil
}
