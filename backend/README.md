# Backend API

## 概要

このバックエンドAPIは、ogen（OpenAPI Generator for Go）を使用して自動生成されたコードを基盤としており、entフレームワークによるORMを利用しています。

## セットアップ

### 1. 環境変数の設定

`.env.example`をコピーして`.env`ファイルを作成します：

```bash
cp .env.example .env
```

### 2. Google OAuth2の設定

1. [Google Cloud Console](https://console.cloud.google.com/)にアクセス
2. プロジェクトを作成または選択
3. 「APIとサービス」→「認証情報」に移動
4. 「認証情報を作成」→「OAuthクライアントID」を選択
5. アプリケーションの種類：「ウェブアプリケーション」を選択
6. 承認済みのリダイレクトURIに以下を追加：
   - `http://localhost:8080/api/v1/auth/callback`（開発環境）
   - 本番環境のURLも追加
7. 作成されたクライアントIDとクライアントシークレットを`.env`ファイルに設定：

```env
GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/callback
```

### 3. JWT秘密鍵の設定

`.env`ファイルで強力なランダム文字列を設定します：

```bash
# ランダムな秘密鍵を生成
openssl rand -base64 32
```

生成された文字列を`.env`に設定：

```env
JWT_SECRET=生成された秘密鍵
```

### 4. 認証フロー

#### ログインフロー：

1. クライアントが `/api/v1/auth/login` にアクセス
2. Google OAuth2の認証ページにリダイレクト
3. ユーザーがGoogleでログイン
4. `/api/v1/auth/callback` にリダイレクト
5. バックエンドがGoogleからユーザー情報を取得
6. メールアドレスでユーザーを検索、存在しない場合は新規作成
7. JWTトークン（access_token、refresh_token）をCookieに設定

#### トークンの保存方法：

- `access_token`: HttpOnly Cookie、有効期限15分、全パスでアクセス可能
- `refresh_token`: HttpOnly Cookie、有効期限7日間、`/api/v1/auth/refresh`のみでアクセス可能

## 開発ガイド

### 認証されたユーザーIDの取得方法

認証が必要なエンドポイントでは、JWT Bearer認証によってユーザーが認証されます。認証されたユーザーのIDは以下の方法で取得できます：

```go
import "backend/security"

func (h *Handler) SomeAuthenticatedEndpoint(ctx context.Context, ...) (..., error) {
    // contextから認証済みユーザーIDを取得
    userID, ok := security.GetUserIDFromContext(ctx)
    if !ok {
        // ユーザーIDが取得できない場合のエラーハンドリング
        return nil, errors.New("認証情報が見つかりません")
    }
    
    // userIDを使用してビジネスロジックを実装
    // ...
}
```

**ポイント：**
- `security.GetUserIDFromContext(ctx)` を使用してcontextからユーザーIDを取得します
- 戻り値は `(string, bool)` の形式で、第2引数でIDの取得成功/失敗を判定できます
- JWT認証が正常に完了すると、`security.SecurityHandler` によって自動的にcontextにユーザーIDが保存されます

### APIパラメーターの受け取り方

ogenによって自動生成されたハンドラーメソッドは、パラメーターの型に応じて異なる形式で受け取ります。

#### 1. パスパラメーター（Path Parameters）

パスパラメーター（例: `/users/{user_id}`）は、専用のParams構造体として受け取ります：

```go
// UsersUserIDGet implements GET /users/{user_id} operation.
func (h *Handler) UsersUserIDGet(ctx context.Context, params api.UsersUserIDGetParams) (api.UsersUserIDGetRes, error) {
    // パスパラメーターの取得
    userID := params.UserID
    
    // userIDを使用してデータベースクエリなどを実行
    // ...
}
```

#### 2. クエリパラメーター（Query Parameters）

クエリパラメーター（例: `?limit=10&offset=0`）も同じParams構造体のフィールドとして含まれます：

```go
// PostsGet implements GET /posts operation.
func (h *Handler) PostsGet(ctx context.Context, params api.PostsGetParams) (api.PostsGetRes, error) {
    // クエリパラメーターの取得（Optionalの場合）
    limit := params.Limit.Value  // api.OptInt32型
    if params.Limit.IsSet() {
        // limitが指定されている場合の処理
    }
    
    // 必須クエリパラメーターの場合
    // offset := params.Offset  // int32型
}
```

**Optionalパラメーターの扱い：**
- `api.OptXXX` 型（例: `api.OptInt32`, `api.OptString`）として定義されます
- `.IsSet()` メソッドでパラメーターが指定されているか確認できます
- `.Value` フィールドで実際の値を取得します

#### 3. リクエストボディ（Request Body）

POST/PUT等のリクエストボディは、専用のリクエスト構造体として受け取ります：

```go
// PostsPost implements POST /posts operation.
func (h *Handler) PostsPost(ctx context.Context, req *api.PostRequest) (api.PostsPostRes, error) {
    // リクエストボディの各フィールドにアクセス
    title := req.Title.Value
    content := req.Content.Value
    genreID := req.GenreID.Value
    
    // リクエストデータを使用してエンティティを作成
    // ...
}
```

#### 4. リクエストボディとパスパラメーターの組み合わせ

PUT/PATCH等では、両方を受け取る場合があります：

```go
// UsersUserIDPut implements PUT /users/{user_id} operation.
func (h *Handler) UsersUserIDPut(ctx context.Context, req *api.UserRequest, params api.UsersUserIDPutParams) (api.UsersUserIDPutRes, error) {
    // パスパラメーター
    userID := params.UserID
    
    // リクエストボディ
    displayName := req.DisplayName.Value
    bio := req.Bio.Value
    
    // 更新処理を実装
    // ...
}
```

### パラメーター型の命名規則

ogenは以下の命名規則でパラメーター構造体を生成します：

- **Params構造体**: `{エンドポイント名}Params`
  - 例: `GET /users/{user_id}` → `UsersUserIDGetParams`
  - 例: `PUT /posts/{post_id}` → `PostsPostIDPutParams`

- **Request構造体**: OpenAPI定義のスキーマ名に基づく
  - 例: `PostRequest`, `UserRequest`

### 実装例

以下は、認証とパラメーターを組み合わせた完全な実装例です：

```go
func (h *Handler) PostsPost(ctx context.Context, req *api.PostRequest) (api.PostsPostRes, error) {
    // 1. 認証されたユーザーIDを取得
    userID, ok := security.GetUserIDFromContext(ctx)
    if !ok {
        return nil, errors.New("認証情報が見つかりません")
    }
    
    // 2. リクエストボディからパラメーターを取得
    title := req.Title.Value
    content := req.Content.Value
    genreID := req.GenreID.Value
    
    // 3. データベース操作を実行
    post, err := h.client.Post.Create().
        SetTitle(title).
        SetContent(content).
        SetUserID(userID).
        SetGenreID(genreID).
        Save(ctx)
    if err != nil {
        return nil, err
    }
    
    // 4. レスポンスを返す
    return &api.Post{
        ID:        api.NewOptString(post.ID),
        Title:     api.NewOptString(post.Title),
        Content:   api.NewOptString(post.Content),
        // ...
    }, nil
}
```

## コード生成

APIスキーマ（`docs/api.yaml`）を変更した後、以下のコマンドでコードを再生成します：

```bash
go generate ./...
```

## 参考

- `security/handler.go` - JWT認証の実装
- `handler/` - 各エンドポイントのハンドラー実装
- `api/` - ogenによって自動生成されたコード
