package handler

import (
	"context"

	"backend/api"
)

// GoalsGet は現在のユーザーの目標一覧取得のモック実装です
func (h *Handler) GoalsGet(ctx context.Context, params api.GoalsGetParams) (api.GoalsGetRes, error) {
	// TODO: 実装
	goals := api.GoalsGetOKApplicationJSON([]api.Goal{})
	return &goals, nil
}

// GoalsPost は新規目標作成のモック実装です
func (h *Handler) GoalsPost(ctx context.Context, req *api.GoalRequest) (api.GoalsPostRes, error) {
	// TODO: 実装
	return &api.Goal{}, nil
}

// GoalsGoalIDDelete は目標削除のモック実装です
func (h *Handler) GoalsGoalIDDelete(ctx context.Context, params api.GoalsGoalIDDeleteParams) (api.GoalsGoalIDDeleteRes, error) {
	// TODO: 実装
	return &api.GoalsGoalIDDeleteNoContent{}, nil
}

// GoalsGoalIDGet は目標詳細取得のモック実装です
func (h *Handler) GoalsGoalIDGet(ctx context.Context, params api.GoalsGoalIDGetParams) (api.GoalsGoalIDGetRes, error) {
	// TODO: 実装
	return &api.Error{}, nil
}

// GoalsGoalIDPut は目標更新のモック実装です
func (h *Handler) GoalsGoalIDPut(ctx context.Context, req *api.GoalRequest, params api.GoalsGoalIDPutParams) (api.GoalsGoalIDPutRes, error) {
	// TODO: 実装
	return &api.Goal{}, nil
}

// UsersUserIDGoalsGet は指定ユーザーの目標一覧取得のモック実装です
func (h *Handler) UsersUserIDGoalsGet(ctx context.Context, params api.UsersUserIDGoalsGetParams) (api.UsersUserIDGoalsGetRes, error) {
	// TODO: 実装
	goals := api.UsersUserIDGoalsGetOKApplicationJSON([]api.Goal{})
	return &goals, nil
}
