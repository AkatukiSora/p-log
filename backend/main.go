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

	if err := dev(client, jwtHandler); err != nil {
		log.Fatalf("failed to run dev function: %v", err)
	}

	// サーバーの起動
	mux := http.NewServeMux()
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", srv))

	// 開発環境用: CORSをすべて許可
	handler := enableCORSForDevelopment(mux)

	log.Println("Starting server on :8080 with prefix /api/v1")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

// enableCORSForDevelopment は開発環境用にすべてのCORSリクエストを許可するミドルウェアです
// 本番環境では使用しないでください
func enableCORSForDevelopment(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 開発環境用: すべてのオリジンを許可
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type, Set-Cookie, Location")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// プリフライトリクエストの処理
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
