# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

p-logのバックエンドAPI。Go言語で構築されており、ogen（OpenAPI Generator for Go）による自動生成コードとentフレームワークによるORMを利用。

## 開発コマンド

```bash
# 開発サーバ起動（Docker Compose）
docker compose up

# Swagger UI（APIドキュメント）
# http://localhost:8081 でアクセス可能

# OpenAPI仕様からAPIコード生成
go generate ./...

# entスキーマからORMコード生成
go generate ./ent/...
```

## アーキテクチャ

```
backend/
├── api/              # ogen自動生成コード（編集禁止）
├── ent/              # ent ORM
│   └── schema/       # エンティティスキーマ定義（編集対象）
├── handler/          # APIエンドポイント実装
├── internal/
│   ├── db/           # DB接続・マイグレーション
│   ├── image/        # GCS画像サービス
│   ├── jwt/          # JWT認証処理
│   └── oauth/        # Google OAuth2
├── security/         # JWT Bearer認証ハンドラ
├── main.go           # エントリポイント
└── dev.go            # 開発用初期データ
```

### コード生成フロー

1. **API変更時**: `../docs/api.yaml`を編集 → `go generate ./...` → `api/`が再生成
2. **スキーマ変更時**: `ent/schema/`を編集 → `go generate ./ent/...` → `ent/`の他ファイルが再生成

### エンティティ構造

- User: ユーザー（フォロー/フォロワー関係あり）
- Post: 投稿（User/Goalへの関連必須）
- Goal: 目標
- Genre: ジャンル（Userと多対多）
- Image: 画像（GCS連携）
- Reaction: リアクション
- RefreshToken: リフレッシュトークン

## 実装パターン

### 認証済みユーザーIDの取得

```go
userID, err := security.GetUserIDFromContext(ctx)
if err != nil {
    return nil, ErrUnauthorized
}
```

### ogenのパラメータ取得

- パスパラメータ: `params.UserID`
- オプショナルパラメータ: `params.Limit.Value`（`params.Limit.IsSet()`で存在確認）
- リクエストボディ: `req.Title.Value`

### エラーハンドリング

`handler/error.go`でエラー型を定義。`NewError()`でHTTPステータスコードへ変換。

## 外部サービス

- **PostgreSQL**: ローカルまたはCloud SQL
- **Google Cloud Storage**: 画像保存
- **Google OAuth2**: 認証
