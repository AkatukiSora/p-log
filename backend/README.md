# Backend API

## 概要

このバックエンドAPIは、ogen（OpenAPI Generator for Go）を使用して自動生成されたコードを基盤としており、entフレームワークによるORMを利用しています。

## セットアップ

### 1. 環境変数の設定

`.env.example`をコピーして`.env`ファイルを作成します：

```bash
cp .env.example .env
```

### 2. データベース設定

#### ローカル開発環境

Docker Composeを使用してPostgreSQLを起動します：

```bash
docker compose up -d
```

`.env`ファイルでローカルデータベースの設定を確認：

```env
USE_CLOUD_SQL=false
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=p-log
DB_SSLMODE=disable
```

#### 本番環境（Google Cloud SQL）

Google Cloud SQLのPostgreSQLインスタンスを使用する場合の設定：

**前提条件：**
- Google Cloud SQLインスタンスが作成済みであること
- IAM認証を有効化していること（推奨）
- Cloud Runから接続する場合、サービスアカウントに必要な権限が付与されていること
  - `Cloud SQL Client`ロール

**環境変数の設定：**

```env
USE_CLOUD_SQL=true
CLOUD_SQL_INSTANCE=your-project-id:your-region:your-instance-name
DB_USER=your-database-user
DB_NAME=p-log
USE_IAM_AUTH=true
```

**注意事項：**
- `CLOUD_SQL_INSTANCE`は`project-id:region:instance-name`の形式で指定
- `USE_IAM_AUTH=true`の場合、パスワードは不要（IAM認証）
- Cloud Run上で実行する場合、サービスアカウントの認証情報が自動的に使用されます
- ローカルでCloud SQLに接続する場合は、Application Default Credentials（ADC）が必要：
  ```bash
  gcloud auth application-default login
  ```

**Cloud Runへのデプロイ時の設定：**

```bash
gcloud run deploy your-service \
  --add-cloudsql-instances=your-project-id:your-region:your-instance-name \
  --set-env-vars="USE_CLOUD_SQL=true,CLOUD_SQL_INSTANCE=your-project-id:your-region:your-instance-name,DB_USER=your-user,DB_NAME=p-log,USE_IAM_AUTH=true"
```

### 3. Google Cloud Storage設定

画像アップロード機能を使用するには、Google Cloud Storageバケットが必要です。

**前提条件：**
- Google Cloud Storageバケットが作成済みであること
- サービスアカウントにバケットへの書き込み権限があること
  - `Storage Object Admin`ロールまたは`Storage Object Creator`ロール

**環境変数の設定：**

```env
GCS_BUCKET_NAME=your-bucket-name
```

**IAM認証：**
- ローカル開発環境では`GOOGLE_APPLICATION_CREDENTIALS`環境変数でサービスアカウントキーを指定
- Cloud Run上では、サービスアカウントの認証情報が自動的に使用されます

**バケットの作成例：**

```bash
# バケット作成
gsutil mb -p your-project-id -c STANDARD -l asia-northeast1 gs://your-bucket-name

# サービスアカウントに権限を付与
gsutil iam ch serviceAccount:your-service-account@your-project.iam.gserviceaccount.com:objectAdmin gs://your-bucket-name
```

### 4. Google OAuth2の設定

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

### 5. JWT秘密鍵の設定

`.env`ファイルで強力なランダム文字列を設定します：

```bash
# ランダムな秘密鍵を生成
openssl rand -base64 32
```

生成された文字列を`.env`に設定：

```env
JWT_SECRET=生成された秘密鍵
```

### 6. 認証フロー

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

### 画像アップロード機能の使用方法

画像アップロード機能は、Google Cloud Storageを使用して実装されています。

#### 画像アップロードエンドポイント（POST /images）

```go
// ImagesPost implements POST /images operation.
func (h *Handler) ImagesPost(ctx context.Context, req api.OptImagesPostReq) (api.ImagesPostRes, error) {
    // 1. リクエストから画像ファイルを取得
    multipartFile := req.Value.GetImage().Value
    
    // 2. contextから認証済みユーザーIDを取得
    userIDStr, err := security.GetUserIDFromContext(ctx)
    if err != nil {
        return nil, ErrUnauthorized
    }
    
    // 3. GCSに画像をアップロード（UUID付きファイル名で保存）
    objectName, contentType, err := h.imageService.UploadImage(ctx, fileReader, multipartFile.Name, multipartFile.Header)
    
    // 4. DBに画像メタデータ（UUID、オブジェクト名、コンテンツタイプ）を保存
    img, err := h.client.Image.Create().
        SetObjectName(objectName).
        SetContentType(contentType).
        SetUploadedByID(userID).
        Save(ctx)
    
    // 5. 画像IDとURLを返す
    return &api.Image{
        ID:  img.ID,    // UUIDが返される
        URL: fmt.Sprintf("/api/v1/images/%s", img.ID.String()),
    }, nil
}
```

#### 画像取得エンドポイント（GET /images/{image_id}）

```go
// ImagesImageIDGet implements GET /images/{image_id} operation.
func (h *Handler) ImagesImageIDGet(ctx context.Context, params api.ImagesImageIDGetParams) (api.ImagesImageIDGetRes, error) {
    // 1. DBから画像メタデータを取得
    img, err := h.client.Image.Get(ctx, params.ImageID)
    
    // 2. GCSから画像データを取得
    reader, _, err := h.imageService.GetImage(ctx, img.ObjectName)
    
    // 3. 画像データをストリームとして返す
    return &api.ImagesImageIDGetOK{
        Data: reader,
    }, nil
}
```

**ポイント：**
- 画像ファイル名はUUIDで自動生成されるため、重複の心配がありません
- DBにはUUID、GCSオブジェクト名、コンテンツタイプが保存されます
- GCSへのアクセスはIAM認証で行われるため、安全です
- 画像取得時は、DBからメタデータを取得してからGCSから画像データを読み込みます

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
