package handler

import (
	"context"

	"backend/api"
	"backend/ent"
	"backend/ent/goal"
	"backend/ent/image"
	"backend/ent/post"
	"backend/ent/user"
	"backend/security"

	"github.com/google/uuid"
)

// PostsGet implements GET /posts operation.
// 現在のユーザーの投稿一覧取得
func (h *Handler) PostsGet(ctx context.Context, params api.PostsGetParams) (api.PostsGetRes, error) {
	// コンテキストからユーザーIDを取得
	userIDStr, err := security.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ErrBadRequest
	}

	// クエリを構築
	query := h.client.Post.
		Query().
		Where(post.HasUserWith(user.IDEQ(userID))).
		WithUser().
		WithGoal().
		WithImages().
		WithReactions()

	// goal_idフィルタを適用
	if goalID, ok := params.GoalID.Get(); ok {
		query = query.Where(post.HasGoalWith(goal.IDEQ(goalID)))
	}

	// 作成日時の降順でソート
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

	// 投稿一覧を取得
	posts, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	// entのPostをapiのPostに変換
	apiPosts, err := h.convertPostsToAPI(posts)
	if err != nil {
		return nil, err
	}

	response := api.PostsGetOKApplicationJSON(apiPosts)
	return &response, nil
}

// PostsPost implements POST /posts operation.
// 進捗投稿作成
func (h *Handler) PostsPost(ctx context.Context, req *api.PostRequest) (api.PostsPostRes, error) {
	// コンテキストからユーザーIDを取得
	userIDStr, err := security.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ErrBadRequest
	}

	// ユーザーが存在するか確認
	exists, err := h.client.User.Query().Where(user.IDEQ(userID)).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound
	}

	// ゴールが存在するか確認
	exists, err = h.client.Goal.Query().Where(goal.IDEQ(req.GoalID)).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound
	}

	// 投稿を作成
	postBuilder := h.client.Post.
		Create().
		SetContent(req.Content).
		SetUserID(userID).
		SetGoalID(req.GoalID)

	p, err := postBuilder.Save(ctx)
	if err != nil {
		return nil, err
	}

	// 画像IDが指定されている場合は、画像と投稿を紐付け
	if len(req.ImageIds) > 0 {
		for _, imageID := range req.ImageIds {
			// 画像が存在するか確認
			exists, err := h.client.Image.Query().Where(image.IDEQ(imageID)).Exist(ctx)
			if err != nil {
				return nil, err
			}
			if !exists {
				// 作成した投稿をロールバック
				h.client.Post.DeleteOne(p).Exec(ctx)
				return nil, ErrNotFound
			}

			// 画像を投稿に紐付け
			err = h.client.Image.UpdateOneID(imageID).SetPostID(p.ID).Exec(ctx)
			if err != nil {
				// 作成した投稿をロールバック
				h.client.Post.DeleteOne(p).Exec(ctx)
				return nil, err
			}
		}
	}

	// レスポンス用に投稿を再取得（関連データを含む）
	p, err = h.client.Post.Query().
		Where(post.IDEQ(p.ID)).
		WithUser().
		WithGoal().
		WithImages().
		WithReactions().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// レスポンスを構築
	response, err := h.convertPostToAPI(p)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// PostsPostIDDelete implements DELETE /posts/{post_id} operation.
// 投稿を削除（紐づいている画像も同時に削除）
func (h *Handler) PostsPostIDDelete(ctx context.Context, params api.PostsPostIDDeleteParams) (api.PostsPostIDDeleteRes, error) {
	// コンテキストからユーザーIDを取得
	userIDStr, err := security.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ErrBadRequest
	}

	// 投稿が存在し、かつ自分の投稿であることを確認
	p, err := h.client.Post.Query().
		Where(post.IDEQ(params.PostID)).
		WithUser().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// 投稿の所有者が現在のユーザーか確認
	postUser, err := p.Edges.UserOrErr()
	if err != nil {
		return nil, err
	}
	if postUser.ID != userID {
		return nil, ErrForbidden
	}

	// 投稿を削除（Cascadeで画像とリアクションも削除される）
	err = h.client.Post.DeleteOneID(params.PostID).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &api.PostsPostIDDeleteNoContent{}, nil
}

// PostsPostIDGet implements GET /posts/{post_id} operation.
// 投稿詳細取得
func (h *Handler) PostsPostIDGet(ctx context.Context, params api.PostsPostIDGetParams) (api.PostsPostIDGetRes, error) {
	// 投稿を取得
	p, err := h.client.Post.Query().
		Where(post.IDEQ(params.PostID)).
		WithUser().
		WithGoal().
		WithImages().
		WithReactions().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// レスポンスを構築
	response, err := h.convertPostToAPI(p)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// PostsPostIDPut implements PUT /posts/{post_id} operation.
// 投稿更新
func (h *Handler) PostsPostIDPut(ctx context.Context, req *api.PostRequest, params api.PostsPostIDPutParams) (api.PostsPostIDPutRes, error) {
	// コンテキストからユーザーIDを取得
	userIDStr, err := security.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ErrBadRequest
	}

	// 投稿が存在し、かつ自分の投稿であることを確認
	p, err := h.client.Post.Query().
		Where(post.IDEQ(params.PostID)).
		WithUser().
		WithImages().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// 投稿の所有者が現在のユーザーか確認
	postUser, err := p.Edges.UserOrErr()
	if err != nil {
		return nil, err
	}
	if postUser.ID != userID {
		return nil, ErrForbidden
	}

	// ゴールが存在するか確認
	exists, err := h.client.Goal.Query().Where(goal.IDEQ(req.GoalID)).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound
	}

	// 投稿を更新
	p, err = h.client.Post.UpdateOneID(params.PostID).
		SetContent(req.Content).
		SetGoalID(req.GoalID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// 既存の画像を取得
	existingImages, err := p.Edges.ImagesOrErr()
	if err != nil {
		existingImages = []*ent.Image{}
	}

	// 既存の画像IDのマップを作成
	existingImageIDs := make(map[uuid.UUID]bool)
	for _, img := range existingImages {
		existingImageIDs[img.ID] = true
	}

	// 新しい画像IDのマップを作成
	newImageIDs := make(map[uuid.UUID]bool)
	for _, imageID := range req.ImageIds {
		newImageIDs[imageID] = true
	}

	// 削除する画像を特定（既存にあって新しいリストにないもの）
	for _, img := range existingImages {
		if !newImageIDs[img.ID] {
			// 画像の投稿IDをnullに設定
			err = h.client.Image.UpdateOneID(img.ID).ClearPost().Exec(ctx)
			if err != nil {
				return nil, err
			}
		}
	}

	// 追加する画像を特定（新しいリストにあって既存にないもの）
	for _, imageID := range req.ImageIds {
		if !existingImageIDs[imageID] {
			// 画像が存在するか確認
			exists, err := h.client.Image.Query().Where(image.IDEQ(imageID)).Exist(ctx)
			if err != nil {
				return nil, err
			}
			if !exists {
				return nil, ErrNotFound
			}

			// 画像を投稿に紐付け
			err = h.client.Image.UpdateOneID(imageID).SetPostID(p.ID).Exec(ctx)
			if err != nil {
				return nil, err
			}
		}
	}

	// レスポンス用に投稿を再取得（関連データを含む）
	p, err = h.client.Post.Query().
		Where(post.IDEQ(p.ID)).
		WithUser().
		WithGoal().
		WithImages().
		WithReactions().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// レスポンスを構築
	response, err := h.convertPostToAPI(p)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UsersUserIDPostsGet implements GET /users/{user_id}/posts operation.
// 指定ユーザーの投稿一覧取得
func (h *Handler) UsersUserIDPostsGet(ctx context.Context, params api.UsersUserIDPostsGetParams) (api.UsersUserIDPostsGetRes, error) {
	// ユーザーが存在するか確認
	exists, err := h.client.User.Query().Where(user.IDEQ(params.UserID)).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound
	}

	// クエリを構築
	query := h.client.Post.
		Query().
		Where(post.HasUserWith(user.IDEQ(params.UserID))).
		WithUser().
		WithGoal().
		WithImages().
		WithReactions()

	// goal_idフィルタを適用
	if goalID, ok := params.GoalID.Get(); ok {
		query = query.Where(post.HasGoalWith(goal.IDEQ(goalID)))
	}

	// 作成日時の降順でソート
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

	// 投稿一覧を取得
	posts, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	// entのPostをapiのPostに変換
	apiPosts, err := h.convertPostsToAPI(posts)
	if err != nil {
		return nil, err
	}

	response := api.UsersUserIDPostsGetOKApplicationJSON(apiPosts)
	return &response, nil
}

// convertPostToAPI はent.Postをapi.Postに変換します
func (h *Handler) convertPostToAPI(p *ent.Post) (*api.Post, error) {
	// ユーザーを取得
	postUser, err := p.Edges.UserOrErr()
	if err != nil {
		return nil, err
	}

	// ゴールを取得
	postGoal, err := p.Edges.GoalOrErr()
	if err != nil {
		return nil, err
	}

	// 画像URLを取得
	var imageURLs []string
	if images, err := p.Edges.ImagesOrErr(); err == nil {
		for _, img := range images {
			// TODO: 実際の画像URLを生成（GCSのパスやCDNのURLなど）
			// 現状はobject_nameをそのまま使用
			imageURLs = append(imageURLs, img.ObjectName)
		}
	}

	// リアクション数を取得
	var reactionCount int
	if reactions, err := p.Edges.ReactionsOrErr(); err == nil {
		reactionCount = len(reactions)
	}

	return &api.Post{
		ID:            p.ID,
		UserID:        postUser.ID,
		GoalID:        postGoal.ID,
		Content:       p.Content,
		ImageUrls:     imageURLs,
		ReactionCount: reactionCount,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}, nil
}

// convertPostsToAPI はent.Postのスライスをapi.Postのスライスに変換します
func (h *Handler) convertPostsToAPI(posts []*ent.Post) ([]api.Post, error) {
	apiPosts := make([]api.Post, 0, len(posts))
	for _, p := range posts {
		apiPost, err := h.convertPostToAPI(p)
		if err != nil {
			return nil, err
		}
		apiPosts = append(apiPosts, *apiPost)
	}
	return apiPosts, nil
}
