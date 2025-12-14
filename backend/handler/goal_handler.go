package handler

import (
	"context"

	"backend/api"
	"backend/ent"
	"backend/ent/goal"
	"backend/ent/user"
	"backend/security"

	"github.com/google/uuid"
)

// GoalsGet は現在のユーザーの目標一覧取得を実装します
func (h *Handler) GoalsGet(ctx context.Context, params api.GoalsGetParams) (api.GoalsGetRes, error) {
	// コンテキストからユーザーIDを取得
	userIDStr, err := security.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ErrBadRequest
	}

	// ユーザーの目標一覧を取得
	goals, err := h.client.Goal.
		Query().
		Where(goal.HasUserWith(user.IDEQ(userID))).
		WithUser().
		All(ctx)
	if err != nil {
		return nil, err
	}

	// entのGoalをapiのGoalに変換
	apiGoals := make([]api.Goal, 0, len(goals))
	for _, g := range goals {
		apiGoal := api.Goal{
			ID:        g.ID,
			UserID:    userID,
			Title:     g.Title,
			CreatedAt: g.CreatedAt,
		}
		if g.Deadline != nil {
			apiGoal.Deadline.SetTo(*g.Deadline)
		}
		apiGoals = append(apiGoals, apiGoal)
	}

	response := api.GoalsGetOKApplicationJSON(apiGoals)
	return &response, nil
}

// GoalsPost は新規目標作成を実装します
func (h *Handler) GoalsPost(ctx context.Context, req *api.GoalRequest) (api.GoalsPostRes, error) {
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

	// 目標を作成
	goalBuilder := h.client.Goal.
		Create().
		SetTitle(req.Title).
		SetUserID(userID)

	// デッドラインが設定されている場合は追加
	if deadline, ok := req.Deadline.Get(); ok {
		goalBuilder.SetDeadline(deadline)
	}

	g, err := goalBuilder.Save(ctx)
	if err != nil {
		return nil, err
	}

	// レスポンスを構築
	response := api.Goal{
		ID:        g.ID,
		UserID:    userID,
		Title:     g.Title,
		CreatedAt: g.CreatedAt,
	}
	if g.Deadline != nil {
		response.Deadline.SetTo(*g.Deadline)
	}

	return &response, nil
}

// GoalsGoalIDDelete は目標削除を実装します
func (h *Handler) GoalsGoalIDDelete(ctx context.Context, params api.GoalsGoalIDDeleteParams) (api.GoalsGoalIDDeleteRes, error) {
	// コンテキストからユーザーIDを取得
	userIDStr, err := security.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ErrBadRequest
	}

	// 目標が存在し、かつ自分の目標であることを確認
	g, err := h.client.Goal.Query().
		Where(goal.IDEQ(params.GoalID)).
		WithUser().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// 目標の所有者が現在のユーザーか確認
	goalUser, err := g.Edges.UserOrErr()
	if err != nil {
		return nil, err
	}
	if goalUser.ID != userID {
		return nil, ErrForbidden
	}

	// 目標を削除
	err = h.client.Goal.DeleteOneID(params.GoalID).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &api.GoalsGoalIDDeleteNoContent{}, nil
}

// GoalsGoalIDGet は目標詳細取得を実装します
func (h *Handler) GoalsGoalIDGet(ctx context.Context, params api.GoalsGoalIDGetParams) (api.GoalsGoalIDGetRes, error) {

	// 目標を取得
	g, err := h.client.Goal.Query().
		Where(goal.IDEQ(params.GoalID)).
		WithUser().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// 目標の所有者を取得
	goalUser, err := g.Edges.UserOrErr()
	if err != nil {
		return nil, err
	}

	// レスポンスを構築
	response := api.Goal{
		ID:        g.ID,
		UserID:    goalUser.ID,
		Title:     g.Title,
		CreatedAt: g.CreatedAt,
	}
	if g.Deadline != nil {
		response.Deadline.SetTo(*g.Deadline)
	}

	return &response, nil
}

// GoalsGoalIDPut は目標更新を実装します
func (h *Handler) GoalsGoalIDPut(ctx context.Context, req *api.GoalRequest, params api.GoalsGoalIDPutParams) (api.GoalsGoalIDPutRes, error) {
	// コンテキストからユーザーIDを取得
	userIDStr, err := security.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ErrBadRequest
	}

	// 目標が存在し、かつ自分の目標であることを確認
	g, err := h.client.Goal.Query().
		Where(goal.IDEQ(params.GoalID)).
		WithUser().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// 目標の所有者が現在のユーザーか確認
	goalUser, err := g.Edges.UserOrErr()
	if err != nil {
		return nil, err
	}
	if goalUser.ID != userID {
		return nil, ErrForbidden
	}

	// 目標を更新
	updateBuilder := h.client.Goal.
		UpdateOneID(params.GoalID).
		SetTitle(req.Title)

	// デッドラインが設定されている場合は更新
	if deadline, ok := req.Deadline.Get(); ok {
		updateBuilder.SetDeadline(deadline)
	} else {
		updateBuilder.ClearDeadline()
	}

	updatedGoal, err := updateBuilder.Save(ctx)
	if err != nil {
		return nil, err
	}

	// レスポンスを構築
	response := api.Goal{
		ID:        updatedGoal.ID,
		UserID:    userID,
		Title:     updatedGoal.Title,
		CreatedAt: updatedGoal.CreatedAt,
	}
	if updatedGoal.Deadline != nil {
		response.Deadline.SetTo(*updatedGoal.Deadline)
	}

	return &response, nil
}

// UsersUserIDGoalsGet は指定ユーザーの目標一覧取得を実装します
func (h *Handler) UsersUserIDGoalsGet(ctx context.Context, params api.UsersUserIDGoalsGetParams) (api.UsersUserIDGoalsGetRes, error) {
	// 指定されたユーザーが存在するか確認
	exists, err := h.client.User.Query().Where(user.IDEQ(params.UserID)).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound
	}

	// 指定ユーザーの目標一覧を取得
	goals, err := h.client.Goal.
		Query().
		Where(goal.HasUserWith(user.IDEQ(params.UserID))).
		WithUser().
		All(ctx)
	if err != nil {
		return nil, err
	}

	// entのGoalをapiのGoalに変換
	apiGoals := make([]api.Goal, 0, len(goals))
	for _, g := range goals {
		apiGoal := api.Goal{
			ID:        g.ID,
			UserID:    params.UserID,
			Title:     g.Title,
			CreatedAt: g.CreatedAt,
		}
		if g.Deadline != nil {
			apiGoal.Deadline.SetTo(*g.Deadline)
		}
		apiGoals = append(apiGoals, apiGoal)
	}

	response := api.UsersUserIDGoalsGetOKApplicationJSON(apiGoals)
	return &response, nil
}
