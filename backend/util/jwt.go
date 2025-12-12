package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrInvalidToken はトークンが無効な場合のエラーです。
	ErrInvalidToken = errors.New("invalid token")

	// ErrExpiredToken はトークンの有効期限が切れている場合のエラーです。
	ErrExpiredToken = errors.New("token expired")

	// ErrInvalidClaims はクレームが不正な場合のエラーです。
	ErrInvalidClaims = errors.New("invalid claims")
)

// JWTClaims はJWTに含まれるカスタムクレームです。
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// JWTConfig はJWTの設定を保持します。
type JWTConfig struct {
	SecretKey            []byte
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	Issuer               string
	Audience             string
}

// NewJWTConfig は環境変数から設定を読み込んでJWTConfigを作成します。
func NewJWTConfig() *JWTConfig {
	secretKey := getEnv("JWT_SECRET", "your-secret-key-change-in-production")
	issuer := getEnv("JWT_ISSUER", "p-log")
	audience := getEnv("JWT_AUDIENCE", "p-log-users")

	return &JWTConfig{
		SecretKey:            []byte(secretKey),
		AccessTokenDuration:  15 * time.Minute,   // アクセストークンは15分
		RefreshTokenDuration: 7 * 24 * time.Hour, // リフレッシュトークンは7日
		Issuer:               issuer,
		Audience:             audience,
	}
}

// GenerateAccessToken はアクセストークンを生成します。
func (c *JWTConfig) GenerateAccessToken(userID, email string) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(c.AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    c.Issuer,
			Audience:  []string{c.Audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(c.SecretKey)
}

// GenerateRefreshToken はリフレッシュトークンを生成します。
func (c *JWTConfig) GenerateRefreshToken(userID string) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(c.RefreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    c.Issuer,
			Audience:  []string{c.Audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(c.SecretKey)
}

// ValidateToken はJWTトークンを検証し、クレームを返します。
func (c *JWTConfig) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 署名方法の検証
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return c.SecretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidClaims
	}

	// 発行者の検証
	if claims.Issuer != c.Issuer {
		return nil, ErrInvalidClaims
	}

	// オーディエンスの検証
	validAudience := false
	for _, aud := range claims.Audience {
		if aud == c.Audience {
			validAudience = true
			break
		}
	}
	if !validAudience {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

// RefreshAccessToken はリフレッシュトークンを検証し、新しいアクセストークンを生成します。
func (c *JWTConfig) RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := c.ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}

	// リフレッシュトークンには通常Emailが含まれていないため、
	// ここではUserIDのみを使用して新しいアクセストークンを生成
	// 実際の実装では、UserIDからEmailを取得する必要があります
	return c.GenerateAccessToken(claims.UserID, claims.Email)
}
