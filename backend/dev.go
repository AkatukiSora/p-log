package main

import (
	"context"
	"log"

	"backend/ent"
	"backend/internal/jwt"

	"github.com/google/uuid"
)

// dev は開発用の初期データ投入関数です。
func dev(client *ent.Client, jwtHandler *jwt.JwtHandler) error {
	ctx := context.Background()
	userID := uuid.MustParse("3fa85f64-5717-4562-b3fc-2c963f66afa6")
	email := "aaaaa@aaa.aaa"
	name := "dev-user"

	// ユーザーが存在しない場合は作成する
	if _, err := client.User.Get(ctx, userID); err != nil {
		// Get で見つからない場合は作成を試みる
		if _, cerr := client.User.
			Create().
			SetID(userID).
			SetName(name).
			SetEmail(email).
			Save(ctx); cerr != nil {
			log.Fatalf("failed to create dev user: %v", cerr)
		}
	}

	accessToken, refreshToken, err := jwtHandler.GenerateTokens(userID, email, ctx)
	if err != nil {
		log.Fatalf("failed to generate tokens: %v", err)
	}
	log.Printf("AccessToken: %s\nRefreshToken: %s", accessToken, refreshToken)
	log.Printf("Dev user created with ID: %s, Email: %s", userID.String(), email)
	return nil
}
