# p-log Mock Backend

p-log APIの簡易的なモックサーバーです。固定のレスポンスと基本的なバリデーションのみを実装しています。

## 特徴

- 🚀 Node.js + Expressベースの軽量サーバー
- 📦 固定レスポンスで素早く動作確認
- ✅ 基本的なUUIDやrequiredフィールドのバリデーション
- 🔐 簡易的な認証（Bearerトークンのヘッダーチェックのみ）
- 📝 全エンドポイントをサポート

## セットアップ

### 依存関係のインストール

```bash
npm install
```

## 起動方法

### 通常起動

```bash
npm start
```

### 開発モード（ファイル変更を自動検知）

```bash
npm run dev
```

サーバーは `http://localhost:8080` で起動します。

## 使用方法

### 認証

ほとんどのエンドポイントで認証が必要ですが、モックサーバーでは以下のように任意のトークンで認証できます：

```bash
curl -H "Authorization: Bearer mock_token" http://localhost:8080/auth/me
```

### エンドポイント例

#### ジャンル一覧取得

```bash
curl http://localhost:8080/genres
```

#### ログイン（認証トークン取得）

```bash
curl "http://localhost:8080/auth/callback?code=test&state=test"
```

#### 現在のユーザー情報取得

```bash
curl -H "Authorization: Bearer mock_token" http://localhost:8080/auth/me
```

#### 目標一覧取得

```bash
curl -H "Authorization: Bearer mock_token" http://localhost:8080/goals
```

#### 新規目標作成

```bash
curl -X POST http://localhost:8080/goals \
  -H "Authorization: Bearer mock_token" \
  -H "Content-Type: application/json" \
  -d '{"title":"新しい目標","deadline":"2025-12-31"}'
```

#### タイムライン取得

```bash
curl -H "Authorization: Bearer mock_token" http://localhost:8080/timeline
```

#### 投稿作成

```bash
curl -X POST http://localhost:8080/posts \
  -H "Authorization: Bearer mock_token" \
  -H "Content-Type: application/json" \
  -d '{"goal_id":"323e4567-e89b-12d3-a456-426614174002","content":"今日も頑張りました！"}'
```

#### フレンド一覧取得

```bash
curl -H "Authorization: Bearer mock_token" http://localhost:8080/friends
```

## モックデータ

`mock-data.js`に定義された固定データを使用しています。主なIDは以下の通り：

- ユーザーID: `123e4567-e89b-12d3-a456-426614174000`
- 目標ID: `323e4567-e89b-12d3-a456-426614174002`
- 投稿ID: `523e4567-e89b-12d3-a456-426614174004`
- ジャンルID（プログラミング）: `823e4567-e89b-12d3-a456-426614174007`

## バリデーション

以下の簡易的なバリデーションを実装しています：

- ✅ UUID形式の検証
- ✅ 必須フィールドの存在確認
- ✅ Authorizationヘッダーの存在確認（Bearer形式）

## 制限事項

- データは永続化されません（サーバー再起動で初期化）
- 実際のトークン検証は行いません
- 画像は1x1の透明PNGダミーを返します
- ページネーションは実装していません（パラメータは受け取るが無視）
- エラーメッセージは簡易的です

## ポート変更

デフォルトは8080ですが、環境変数で変更できます：

```bash
PORT=3000 npm start
```

## OpenAPI仕様

詳細なAPI仕様は `/docs/api.yaml` を参照してください。

## トラブルシューティング

### ポートが既に使用されている

```bash
# プロセスを確認
lsof -i :8080

# プロセスを終了
kill -9 <PID>
```

### 依存関係のエラー

```bash
# node_modulesを削除して再インストール
rm -rf node_modules package-lock.json
npm install
```

## ライセンス

MIT
