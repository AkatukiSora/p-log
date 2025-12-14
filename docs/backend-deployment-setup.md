# Backend デプロイ設定ガイド

## 必要な GitHub Secrets

### WIF（Workload Identity Federation）認証

| Secret名 | 説明 | 例 |
|---------|------|-----|
| `GCP_WORKLOAD_IDENTITY_PROVIDER` | Workload Identity Provider のリソース名 | `projects/PROJECT_NUMBER/locations/global/workloadIdentityPools/POOL_NAME/providers/PROVIDER_NAME` |
| `GCP_SERVICE_ACCOUNT` | デプロイ用サービスアカウント | `deploy-sa@p-log-481112.iam.gserviceaccount.com` |

### Cloud SQL

| Secret名 | 説明 | 例 |
|---------|------|-----|
| `CLOUDSQL_INSTANCE` | Cloud SQL インスタンスの接続名 | `p-log-481112:asia-northeast1:p-log-db` |

### Artifact Registry（オプション：デフォルト値を変更する場合）

| Secret名 | 説明 | デフォルト値 |
|---------|------|------------|
| `ARTIFACT_REGISTRY_REPO` | Artifact Registry リポジトリ名 | `p-log` |

## GitHub Secrets 設定手順

1. GitHubリポジトリの **Settings** → **Secrets and variables** → **Actions** にアクセス
2. **New repository secret** をクリック
3. 上記の各Secretを追加

## 必要な IAM 権限

### デプロイ用サービスアカウント（例：deploy-sa@p-log-481112.iam.gserviceaccount.com）

以下の権限が必要です：

```bash
# Cloud Run デプロイ権限
gcloud projects add-iam-policy-binding p-log-481112 \
  --member="serviceAccount:deploy-sa@p-log-481112.iam.gserviceaccount.com" \
  --role="roles/run.admin"

# Artifact Registry へのアクセス
gcloud projects add-iam-policy-binding p-log-481112 \
  --member="serviceAccount:deploy-sa@p-log-481112.iam.gserviceaccount.com" \
  --role="roles/artifactregistry.writer"

# Runtime SA を指定するための権限（backend-runner に対して）
gcloud iam service-accounts add-iam-policy-binding \
  backend-runner@p-log-481112.iam.gserviceaccount.com \
  --member="serviceAccount:deploy-sa@p-log-481112.iam.gserviceaccount.com" \
  --role="roles/iam.serviceAccountUser"

# サービスアカウント情報の取得
gcloud projects add-iam-policy-binding p-log-481112 \
  --member="serviceAccount:deploy-sa@p-log-481112.iam.gserviceaccount.com" \
  --role="roles/iam.serviceAccountViewer"
```

### Runtime サービスアカウント（backend-runner@p-log-481112.iam.gserviceaccount.com）

アプリケーション実行時に必要な権限：

```bash
# Cloud SQL 接続
gcloud projects add-iam-policy-binding p-log-481112 \
  --member="serviceAccount:backend-runner@p-log-481112.iam.gserviceaccount.com" \
  --role="roles/cloudsql.client"

# Cloud Storage（画像アップロード用など）
# 例：特定のバケットに対して objectCreator 権限を付与
gsutil iam ch serviceAccount:backend-runner@p-log-481112.iam.gserviceaccount.com:objectCreator \
  gs://YOUR_BUCKET_NAME

# または、バケット全体の管理権限が必要な場合
gsutil iam ch serviceAccount:backend-runner@p-log-481112.iam.gserviceaccount.com:objectAdmin \
  gs://YOUR_BUCKET_NAME
```

## Workload Identity Federation（WIF）セットアップ

フロントエンドと同じ WIF プールを使用する場合、追加の設定は不要です。  
新規にセットアップする場合は以下を実行：

```bash
# Workload Identity Pool の作成
gcloud iam workload-identity-pools create "github-pool" \
  --project="p-log-481112" \
  --location="global" \
  --display-name="GitHub Actions Pool"

# GitHub Provider の作成
gcloud iam workload-identity-pools providers create-oidc "github-provider" \
  --project="p-log-481112" \
  --location="global" \
  --workload-identity-pool="github-pool" \
  --display-name="GitHub Provider" \
  --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.repository=assertion.repository" \
  --issuer-uri="https://token.actions.githubusercontent.com"

# サービスアカウントと WIF の紐付け
gcloud iam service-accounts add-iam-policy-binding "deploy-sa@p-log-481112.iam.gserviceaccount.com" \
  --project="p-log-481112" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/projects/PROJECT_NUMBER/locations/global/workloadIdentityPools/github-pool/attribute.repository/YOUR_GITHUB_USERNAME/p-log"
```

## Cloud SQL インスタンス接続名の確認方法

```bash
gcloud sql instances describe INSTANCE_NAME \
  --format='value(connectionName)'

# 出力例: p-log-481112:asia-northeast1:p-log-db
```

この値を `CLOUDSQL_INSTANCE` シークレットに設定します。

## トラブルシューティング

### デプロイエラー: "permission denied to use service account"

**原因**: デプロイ用SAに `iam.serviceAccountUser` 権限が付与されていない

**解決策**:
```bash
gcloud iam service-accounts add-iam-policy-binding \
  backend-runner@p-log-481112.iam.gserviceaccount.com \
  --member="serviceAccount:deploy-sa@p-log-481112.iam.gserviceaccount.com" \
  --role="roles/iam.serviceAccountUser"
```

### アプリケーションエラー: DB 接続タイムアウト

**原因**: Cloud SQL インスタンスが Cloud Run に関連付けられていない

**確認**:
```bash
gcloud run services describe p-log-backend \
  --region asia-northeast1 \
  --format='value(spec.template.metadata.annotations."run.googleapis.com/cloudsql-instances")'
```

**解決策**: ワークフローの `CLOUDSQL_INSTANCE` シークレットが正しく設定されているか確認

### アプリケーションエラー: Cloud Storage 403 エラー

**原因**: Runtime SA にバケットへのアクセス権限がない

**解決策**:
```bash
# 読み書き権限が必要な場合
gsutil iam ch serviceAccount:backend-runner@p-log-481112.iam.gserviceaccount.com:objectAdmin \
  gs://YOUR_BUCKET_NAME
```

## ロールバック手順

### 前のリビジョンに戻す

```bash
# 現在のリビジョンを確認
gcloud run revisions list \
  --service=p-log-backend \
  --region=asia-northeast1

# 特定のリビジョンにトラフィックを100%割り当て
gcloud run services update-traffic p-log-backend \
  --region=asia-northeast1 \
  --to-revisions=REVISION_NAME=100
```

### Runtime SA を変更する

```bash
gcloud run services update p-log-backend \
  --region=asia-northeast1 \
  --service-account=PREVIOUS_SA_EMAIL
```

## 手動デプロイコマンド

緊急時に手動でデプロイする場合：

```bash
# イメージのビルドとプッシュ
cd backend
IMAGE_URL="asia-northeast1-docker.pkg.dev/p-log-481112/p-log/p-log-backend:$(git rev-parse HEAD)"
docker build -f dockerfile.prod -t "${IMAGE_URL}" .
docker push "${IMAGE_URL}"

# Cloud Run にデプロイ
gcloud run deploy p-log-backend \
  --image "${IMAGE_URL}" \
  --region asia-northeast1 \
  --service-account backend-runner@p-log-481112.iam.gserviceaccount.com \
  --add-cloudsql-instances "YOUR_CLOUDSQL_INSTANCE" \
  --platform managed \
  --allow-unauthenticated \
  --quiet
```

## 環境変数の追加

アプリケーションに環境変数を追加する場合、ワークフローファイルの `Deploy to Cloud Run` ステップを編集：

```yaml
- name: Deploy to Cloud Run
  run: |
    gcloud run deploy "${{ env.SERVICE_NAME }}" \
      # ... existing flags ...
      --set-env-vars "VAR1=value1,VAR2=value2" \
      --quiet
```

または、機密情報の場合は Secret Manager を使用：

```yaml
--set-secrets "DB_PASSWORD=db-password:latest"
```
