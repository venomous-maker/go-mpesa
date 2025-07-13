package Abstracts

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type TokenManager struct {
	ConsumerKey    string
	ConsumerSecret string
	BaseURL        string
	TokenURL       string
	CachePath      string
}

type tokenCache struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	CreatedAt int64  `json:"created_at"`
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// NewTokenManager initializes a token manager from config
func NewTokenManager(cfg *MpesaConfig) *TokenManager {
	return &TokenManager{
		ConsumerKey:    cfg.ConsumerKey,
		ConsumerSecret: cfg.ConsumerSecret,
		BaseURL:        cfg.BaseURL,
		TokenURL:       "/oauth/v1/generate?grant_type=client_credentials",
		CachePath:      filepath.Join(os.TempDir(), "mpesa_api_token_cache.json"),
	}
}

// GetToken returns a valid access token (cached or freshly requested)
func (tm *TokenManager) GetToken() (string, error) {
	if token := tm.getCachedToken(); token != "" {
		return token, nil
	}

	token, err := tm.requestNewToken()
	if err != nil {
		return "", err
	}

	return token, nil
}

// getCachedToken reads and checks the cached token for validity
func (tm *TokenManager) getCachedToken() string {
	data, err := os.ReadFile(tm.CachePath)
	if err != nil {
		return ""
	}

	var cached tokenCache
	if err := json.Unmarshal(data, &cached); err != nil {
		return ""
	}

	if time.Now().Unix() > cached.ExpiresAt {
		return ""
	}

	return cached.Token
}

// requestNewToken requests a new token and caches it
func (tm *TokenManager) requestNewToken() (string, error) {
	url := tm.BaseURL + tm.TokenURL
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	credentials := base64.StdEncoding.EncodeToString([]byte(tm.ConsumerKey + ":" + tm.ConsumerSecret))
	req.Header.Set("Authorization", "Basic "+credentials)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var tokenResp tokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil || tokenResp.AccessToken == "" {
		return "", errors.New("failed to get access token from M-Pesa API")
	}

	expiresIn := tokenResp.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = 3600 // fallback
	}

	// Subtract buffer
	expiresIn -= 60

	// Cache the token
	tm.cacheToken(tokenResp.AccessToken, time.Now().Unix()+expiresIn)

	return tokenResp.AccessToken, nil
}

// cacheToken writes token details to file
func (tm *TokenManager) cacheToken(token string, expiresAt int64) {
	cache := tokenCache{
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now().Unix(),
	}

	data, _ := json.Marshal(cache)
	_ = os.WriteFile(tm.CachePath, data, 0600)
}

// ClearCache deletes the token cache file
func (tm *TokenManager) ClearCache() {
	if _, err := os.Stat(tm.CachePath); err == nil {
		_ = os.Remove(tm.CachePath)
	}
}
