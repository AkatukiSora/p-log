# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

p-log は目標管理とソーシャル機能を持つWebアプリケーション。フロントエンド(React SPA)、バックエンド(Go API)、モックサーバー(Node.js)の3つのコンポーネントで構成。

## Development Commands

### Frontend (react/)

```bash
cd react
npm install
npm run dev      # 開発サーバー起動 (Vite)
npm run build    # プロダクションビルド
npm run lint     # ESLint実行
```

### Backend (backend/)

```bash
cd backend
cp .env.example .env           # 環境変数設定
docker compose up -d           # PostgreSQL起動
go run main.go dev.go          # サーバー起動
go generate ./...              # APIコード再生成 (docs/api.yaml変更後)
```

### Mock Backend (mock_backend/)

```bash
cd mock_backend
npm install
npm run dev      # 開発モード (nodemon)
npm start        # 通常起動
```

## Architecture

### Frontend: React 19 + Vite + React Router v7

- `react/src/App.jsx` - ルーティング定義 (BrowserRouter)
- `react/src/pages/pathData/path.js` - ルート定数 (ROUTES)
- `react/src/pages/` - ページコンポーネント (CSS Modules使用)
- `react/src/components/` - 再利用可能コンポーネント

### Backend: Go + ogen + ent

- `docs/api.yaml` - OpenAPI仕様 (APIの単一信頼源)
- `backend/api/` - ogenによる自動生成コード (`go generate`で再生成)
- `backend/ent/schema/` - entスキーマ定義 (User, Post, Goal, Image, Reaction, Genre)
- `backend/handler/` - APIハンドラー実装
- `backend/security/` - JWT認証 (`security.GetUserIDFromContext(ctx)`でユーザーID取得)
- `backend/internal/` - 内部パッケージ (db, image, jwt)

### Database Schema

6エンティティ: USER, POST, GOAL, IMAGE, REACTION, GENRE
詳細は `docs/er-diagram.md` 参照。entフレームワークによる外部キー命名規則 (`user_posts`, `goal_posts`など)。

## Key Patterns

### Backend: ogen APIパラメーター

```go
// パスパラメーター + クエリパラメーター
func (h *Handler) XXXGet(ctx context.Context, params api.XXXGetParams) (api.XXXGetRes, error) {
    id := params.ID                    // パスパラメーター
    limit := params.Limit.Value        // Optionalは .IsSet() で確認
}

// リクエストボディ
func (h *Handler) XXXPost(ctx context.Context, req *api.XXXRequest) (api.XXXPostRes, error) {
    title := req.Title.Value
}
```

### Backend: 認証済みユーザー取得

```go
userID, ok := security.GetUserIDFromContext(ctx)
```

## External Services

- Google OAuth2: 認証
- Google Cloud Storage: 画像保存
- Google Cloud SQL: 本番PostgreSQL

## Language

PR、コード内コメント、GitHub上のコミュニケーションは日本語で記述。
