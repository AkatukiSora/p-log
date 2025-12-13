package db

import (
	"backend/ent"
	"backend/internal/other"
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func CreateClient() (*ent.Client, error) {
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
