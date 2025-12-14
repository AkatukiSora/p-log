package jwt

import (
	"backend/ent"
	"backend/ent/refreshtoken"
	"backend/internal/other"
	"context"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	// ErrInvalidToken はトークンが無効な場合のエラーです。
	ErrInvalidToken = errors.New("invalid token")

	// ErrExpiredToken はトークンの有効期限が切れている場合のエラーです。
	ErrExpiredToken = errors.New("token expired")

	// ErrInvalidClaims はクレームが不正な場合のエラーです。
	ErrInvalidClaims = errors.New("invalid claims")

	// ErrRevokedToken はトークンが無効化されている場合のエラーです。
	ErrRevokedToken = errors.New("token revoked")
)

// JWTClaims はJWTに含まれるカスタムクレームです。
type JWTClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	IsRefresh bool   `json:"is_refresh"`
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

type JwtHandler struct {
	jwtConfig *JWTConfig
	client    *ent.Client
}

func NewJWTConfig() *JWTConfig {
	secretKey := other.GetEnv("JWT_SECRET", "your-secret-key-change-in-production")
	issuer := other.GetEnv("JWT_ISSUER", "p-log")
	audience := other.GetEnv("JWT_AUDIENCE", "p-log-users")

	return &JWTConfig{
		SecretKey:            []byte(secretKey),
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 7 * 24 * time.Hour,
		Issuer:               issuer,
		Audience:             audience,
	}
}

// NewJwtHandlerは新しいJwtHandlerインスタンスを作成します。
func NewJwtHandler(jwtConfig *JWTConfig, client *ent.Client) *JwtHandler {

	return &JwtHandler{
		jwtConfig: jwtConfig,
		client:    client,
	}
}

func (c *JwtHandler) storeRefreshToken(token string, userID uuid.UUID, expiresAt time.Time, ctx context.Context) error {
	tokenHashByte := sha512.Sum512([]byte(token))
	tokenHash := hex.EncodeToString(tokenHashByte[:])

	// ユーザーが存在するか事前に検証（外部キー制約違反の早期検出）
	if _, err := c.client.User.Get(ctx, userID); err != nil {
		return err
	}

	// RefreshTokenエンティティを作成してデータベースに保存
	_, err := c.client.RefreshToken.
		Create().
		SetTokenHash(tokenHash).
		SetExpiresAt(expiresAt).
		SetUserID(userID).
		Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *JwtHandler) isRefreshTokenValid(token string, ctx context.Context) (bool, error) {
	// トークンを検証してクレームを取得
	jwtClaims, err := c.ValidateToken(token)
	if err != nil || !jwtClaims.IsRefresh {
		return false, nil
	}

	// トークンのハッシュを計算
	tokenHashByte := sha512.Sum512([]byte(token))
	tokenHash := hex.EncodeToString(tokenHashByte[:])

	// 有効期限の切れていないtokenがあることを確認
	count, err := c.client.RefreshToken.Query().
		Where(
			refreshtoken.TokenHashEQ(tokenHash),
			refreshtoken.Or(
				refreshtoken.ExpiresAtGT(time.Now()),
				refreshtoken.RevokedEQ(false),
			),
		).
		Count(ctx)
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

// GenerateAccessToken はアクセストークンを生成します。
func (c *JwtHandler) generateAccessToken(userID uuid.UUID, email string) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID:    userID.String(),
		Email:     email,
		IsRefresh: false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(c.jwtConfig.AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    c.jwtConfig.Issuer,
			Audience:  []string{c.jwtConfig.Audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(c.jwtConfig.SecretKey)
}

// GenerateRefreshToken はリフレッシュトークンを生成します。
func (c *JwtHandler) generateRefreshToken(userID uuid.UUID, email string, ctx context.Context) (string, error) {
	now := time.Now()

	claims := JWTClaims{
		UserID:    userID.String(),
		Email:     email,
		IsRefresh: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(c.jwtConfig.RefreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    c.jwtConfig.Issuer,
			Audience:  []string{c.jwtConfig.Audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(c.jwtConfig.SecretKey)
	if err != nil {
		return "", err
	}

	// リフレッシュトークンをデータベースに保存
	err = c.storeRefreshToken(tokenString, userID, now.Add(c.jwtConfig.RefreshTokenDuration), ctx)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateToken はJWTトークンを検証し、クレームを返します。
func (c *JwtHandler) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 署名方法の検証
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return c.jwtConfig.SecretKey, nil
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
	if claims.Issuer != c.jwtConfig.Issuer {
		return nil, ErrInvalidClaims
	}

	// オーディエンスの検証
	validAudience := false
	for _, aud := range claims.Audience {
		if aud == c.jwtConfig.Audience {
			validAudience = true
			break
		}
	}
	if !validAudience {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

// GenerateTokens はアクセストークンとリフレッシュトークンを生成します。
func (c *JwtHandler) GenerateTokens(userID uuid.UUID, email string, ctx context.Context) (accessToken string, refreshToken string, err error) {
	accessToken, err = c.generateAccessToken(userID, email)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = c.generateRefreshToken(userID, email, ctx)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// RefreshAccessToken はリフレッシュトークンを検証し、新しいアクセストークンを生成します。
func (c *JwtHandler) RefreshAccessToken(refreshToken string, ctx context.Context) (string, error) {
	// リフレッシュトークンが有効か確認
	isValid, err := c.isRefreshTokenValid(refreshToken, ctx)
	if err != nil {
		return "", err
	}
	if !isValid {
		return "", ErrRevokedToken
	}
	// クレームを取得
	claims, err := c.ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}

	return c.generateAccessToken(uuid.MustParse(claims.UserID), claims.Email)
}

// RevokeRefreshToken はリフレッシュトークンを無効化します。
func (c *JwtHandler) RevokeRefreshToken(refreshToken string, ctx context.Context) error {
	// トークンのハッシュを計算
	tokenHashByte := sha512.Sum512([]byte(refreshToken))
	tokenHash := hex.EncodeToString(tokenHashByte[:])

	// データベースからリフレッシュトークンを削除
	err := c.client.RefreshToken.
		Update().
		Where(refreshtoken.TokenHashEQ(tokenHash)).
		SetRevoked(true).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
