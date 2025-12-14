package db

import (
	"backend/ent"
	"backend/internal/other"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	"cloud.google.com/go/cloudsqlconn"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

func CreateClient() (*ent.Client, error) {
	// Cloud SQLを使用するかどうかを環境変数で判定
	useCloudSQL := other.GetEnv("USE_CLOUD_SQL", "false") == "true"

	if useCloudSQL {
		return createCloudSQLClient()
	}
	return createStandardClient()
}

// Cloud SQLへの接続（IAM認証対応）
func createCloudSQLClient() (*ent.Client, error) {
	// 環境変数から接続情報を取得
	instanceConnectionName := other.GetEnv("CLOUD_SQL_INSTANCE", "") // 例: project:region:instance
	dbUser := other.GetEnv("DB_USER", "postgres")
	dbName := other.GetEnv("DB_NAME", "p-log")
	useIAM := other.GetEnv("USE_IAM_AUTH", "true") == "true"

	if instanceConnectionName == "" {
		log.Fatal("CLOUD_SQL_INSTANCE environment variable is required when USE_CLOUD_SQL=true")
	}

	log.Printf("Connecting to Cloud SQL instance: %s (IAM: %v)", instanceConnectionName, useIAM)

	ctx := context.Background()

	// Cloud SQL Connectorの設定
	opts := []cloudsqlconn.Option{}
	if useIAM {
		// IAM認証を使用（Cloud Run上では自動でサービスアカウントの認証情報を使用）
		opts = append(opts, cloudsqlconn.WithIAMAuthN())
	}

	// Cloud SQL Dialerを作成
	d, err := cloudsqlconn.NewDialer(ctx, opts...)
	if err != nil {
		log.Fatalf("failed to create Cloud SQL dialer: %v", err)
		return nil, err
	}

	// pgx接続設定を構築
	dsn := fmt.Sprintf("user=%s database=%s", dbUser, dbName)
	if !useIAM {
		// IAM認証を使用しない場合はパスワードを追加
		dbPassword := other.GetEnv("DB_PASSWORD", "")
		if dbPassword != "" {
			dsn = fmt.Sprintf("%s password=%s", dsn, dbPassword)
		}
	}

	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("failed to parse pgx config: %v", err)
		return nil, err
	}

	// Cloud SQL Dialerを使用するようにダイヤル関数を設定
	config.DialFunc = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return d.Dial(ctx, instanceConnectionName)
	}

	// pgx driverを登録
	connStr := stdlib.RegisterConnConfig(config)
	sqlDB, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
		return nil, err
	}

	// 接続をテスト
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
		return nil, err
	}

	// Ent clientを作成
	drv := entsql.OpenDB(dialect.Postgres, sqlDB)
	client := ent.NewClient(ent.Driver(drv))

	log.Println("Successfully connected to Cloud SQL")
	return client, nil
}

// 標準的なPostgreSQLへの接続（ローカル開発用）
func createStandardClient() (*ent.Client, error) {
	// 環境変数から接続情報を取得（デフォルト値を設定）
	dbHost := other.GetEnv("DB_HOST", "localhost")
	dbPort := other.GetEnv("DB_PORT", "5432")
	dbUser := other.GetEnv("DB_USER", "postgres")
	dbPassword := other.GetEnv("DB_PASSWORD", "password")
	dbName := other.GetEnv("DB_NAME", "p-log")
	sslmode := other.GetEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbName, dbPassword, sslmode)

	log.Printf("Connecting to database at %s:%s", dbHost, dbPort)

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
		return nil, err
	}

	return client, nil
}

func Migrate(client *ent.Client) {
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	log.Println("Database migration completed successfully.")
}
