package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"backend/internal/other"
)

var (
	// ErrInvalidState はstateパラメータが無効な場合のエラー
	ErrInvalidState = errors.New("invalid oauth state")
	// ErrCodeExchangeFailed はコード交換に失敗した場合のエラー
	ErrCodeExchangeFailed = errors.New("failed to exchange code for token")
	// ErrUserInfoFailed はユーザー情報の取得に失敗した場合のエラー
	ErrUserInfoFailed = errors.New("failed to get user info")
)

// GoogleUserInfo はGoogleから取得するユーザー情報
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
}

// OAuthConfig はOAuth2の設定を保持します
type OAuthConfig struct {
	config      *oauth2.Config
	stateString string
}

// NewOAuthConfig は新しいOAuthConfigを作成します
func NewOAuthConfig() *OAuthConfig {
	clientID := other.GetEnv("GOOGLE_CLIENT_ID", "")
	clientSecret := other.GetEnv("GOOGLE_CLIENT_SECRET", "")
	redirectURL := other.GetEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/auth/callback")

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &OAuthConfig{
		config:      config,
		stateString: generateStateString(),
	}
}

// generateStateString はCSRF対策用のランダムな文字列を生成します
func generateStateString() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// GetAuthURL は認証用のURLを返します
func (o *OAuthConfig) GetAuthURL() string {
	return o.config.AuthCodeURL(o.stateString, oauth2.AccessTypeOffline)
}

// GetState はstateパラメータを返します
func (o *OAuthConfig) GetState() string {
	return o.stateString
}

// ValidateState はstateパラメータを検証します
func (o *OAuthConfig) ValidateState(state string) bool {
	return state == o.stateString
}

// ExchangeCode は認可コードをトークンに交換します
func (o *OAuthConfig) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := o.config.Exchange(ctx, code)
	if err != nil {
		return nil, ErrCodeExchangeFailed
	}
	return token, nil
}

// GetUserInfo はアクセストークンを使ってユーザー情報を取得します
func (o *OAuthConfig) GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := o.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, ErrUserInfoFailed
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrUserInfoFailed
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrUserInfoFailed
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, ErrUserInfoFailed
	}

	return &userInfo, nil
}
