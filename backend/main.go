package main

//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --target api --clean ../docs/api.yaml

import (
	"log"

	"backend/api"
	"backend/handler"
	"backend/internal/db"
	"backend/internal/jwt"
	"backend/security"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	// データベースクライアントの作成
	client, err := db.CreateClient()
	if err != nil {
		log.Fatalf("failed to create database client: %v", err)
	}

	// ハンドラーとセキュリティハンドラーの作成
	jwtConfig := jwt.NewJWTConfig()
	jwtHandler := jwt.NewJwtHandler(jwtConfig, client)
	h, err := handler.NewHandler(client, jwtHandler)
	if err != nil {
		log.Fatalf("failed to create handler: %v", err)
	}

	sec := security.NewSecurityHandler(jwtHandler, client)

	srv, err := api.NewServer(h, sec)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}
	defer client.Close() // 関数終了時にデータベース接続を閉じる

	// データベースのマイグレーション
	db.Migrate(client)

	// サーバーの起動
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
