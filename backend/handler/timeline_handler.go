package handler

import (
	"context"

	"backend/api"
	"backend/ent"
	"backend/ent/goal"
	"backend/ent/post"
	"backend/ent/user"
	"backend/security"

	"github.com/google/uuid"
)

// TimelineGet implements GET /timeline operation.
// タイムライン取得（投稿は新しい順で返される）
func (h *Handler) TimelineGet(ctx context.Context, params api.TimelineGetParams) (api.TimelineGetRes, error) {
	// コンテキストからユーザーIDを取得
	userIDStr, err := security.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ErrBadRequest
	}

	// 現在のユーザーがフォローしているユーザーのIDを取得
	currentUser, err := h.client.User.Query().
		Where(user.IDEQ(userID)).
		WithFollowing().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// フォローしているユーザーのIDリストを作成
	followingIDs := []uuid.UUID{userID} // 自分自身も含める
	if following, err := currentUser.Edges.FollowingOrErr(); err == nil {
		for _, followedUser := range following {
			followingIDs = append(followingIDs, followedUser.ID)
		}
	}

	// タイムラインクエリを構築
	query := h.client.Post.
		Query().
		Where(post.HasUserWith(user.IDIn(followingIDs...))).
		WithUser().
		WithGoal().
		WithImages().
		WithReactions()

	// goal_idフィルタを適用
	if goalID, ok := params.GoalID.Get(); ok {
		query = query.Where(post.HasGoalWith(goal.IDEQ(goalID)))
	}

	// 作成日時の降順でソート（新しい投稿が最初）
	query = query.Order(ent.Desc(post.FieldCreatedAt))

	// ページネーション
	if limit, ok := params.Limit.Get(); ok {
		query = query.Limit(limit)
	}
	if page, ok := params.Page.Get(); ok {
		if limit, ok := params.Limit.Get(); ok {
			offset := (page - 1) * limit
			query = query.Offset(offset)
		}
	}

	// タイムラインの投稿を取得
	posts, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	// entのPostをapiのPostに変換
	apiPosts := make([]api.Post, 0, len(posts))
	for _, p := range posts {
		apiPost, err := h.convertPostToAPI(p)
		if err != nil {
			return nil, err
		}
		apiPosts = append(apiPosts, *apiPost)
	}

	response := api.TimelineGetOKApplicationJSON(apiPosts)
	return &response, nil
}
