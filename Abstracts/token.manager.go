package Abstracts

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

// TokenManager handles OAuth token acquisition and caching for M-Pesa API authentication.
// It automatically manages token lifecycle, including caching valid tokens and refreshing expired ones.
type TokenManager struct {
	ConsumerKey    string // Consumer key for OAuth authentication
	ConsumerSecret string // Consumer secret for OAuth authentication
	BaseURL        string // Base URL for M-Pesa API
	TokenURL       string // OAuth token endpoint path
	CachePath      string // File path for token cache storage

	mu       sync.Mutex  // protects memCache + file operations
	memCache *tokenCache // in-memory cache to avoid frequent FS reads / duplicate requests
}

// tokenCache represents the structure for storing cached tokens.
type tokenCache struct {
	Token     string `json:"token"`      // The cached access token
	ExpiresAt int64  `json:"expires_at"` // Unix timestamp when token expires
	CreatedAt int64  `json:"created_at"` // Unix timestamp when token was created
}

// tokenResponse represents the OAuth token response from M-Pesa API.
type tokenResponse struct {
	AccessToken string `json:"access_token"` // The access token from OAuth response
	ExpiresIn   string `json:"expires_in"`   // Token expiration time in seconds (as string)
}

// NewTokenManager creates a new token manager instance from the provided configuration.
// The token manager handles OAuth authentication and token caching automatically.
//
// Parameters:
//   - cfg: M-Pesa configuration containing consumer credentials and environment settings
//
// Returns:
//   - *TokenManager: A configured token manager ready for token operations
//
// Example:
//
//	cfg := createMpesaConfig()
//	tokenManager := NewTokenManager(cfg)
func NewTokenManager(cfg *MpesaConfig) *TokenManager {
	manager := &TokenManager{
		ConsumerKey:    cfg.GetConsumerKey(),
		ConsumerSecret: cfg.GetConsumerSecret(),
		BaseURL:        cfg.GetBaseURL(),
		TokenURL:       "/oauth/v1/generate?grant_type=client_credentials",
		CachePath:      filepath.Join(os.TempDir(), "mpesa_api_token_cache.json"),
	}
	manager.CachePath = filepath.Join(os.TempDir(), manager.EncryptedCacheFileName())
	return manager
}

// EncryptedCacheFileName (unchanged) ...
func (tm *TokenManager) EncryptedCacheFileName() string {
	_ = "AES-256-CBC"
	password := []byte("mypassword")
	iv := []byte("passwordpassword")
	plaintext := []byte(tm.ConsumerKey + tm.ConsumerSecret + " + Certificate")

	// Ensure key length is 32 bytes for AES-256
	key := make([]byte, 32)
	copy(key, password)

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Error creating cipher:", err)
		return ""
	}

	// CBC mode requires plaintext to be padded to block size
	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	plaintext = append(plaintext, padtext...)

	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	// Base64 encode the ciphertext and append ".json"
	return base64.StdEncoding.EncodeToString(ciphertext) + ".json"
}

// SetCachePath sets the path for the token cache file.
// This method allows customizing the location where the token cache is stored.
//
// Parameters:
//   - path: The new path for the token cache file
//
// Returns:
//   - *TokenManager: The token manager instance for method chaining
//
// Example:
//
//	tokenManager.SetCachePath("/path/to/custom/cache.json")
func (tm *TokenManager) SetCachePath(path string) *TokenManager {
	tm.CachePath = path
	return tm
}

// GetToken returns a valid OAuth access token. Uses in-memory cache first and
// falls back to file cache. Serializes requests to avoid duplicate token calls.
func (tm *TokenManager) GetToken() (string, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// check in-memory cache
	if tm.memCache != nil && time.Now().Unix() < tm.memCache.ExpiresAt {
		return tm.memCache.Token, nil
	}

	// try file cache (and populate in-memory if valid)
	if token := tm.getCachedToken(); token != "" {
		return token, nil
	}

	// no valid cache -> request new token (protected by mutex to avoid duplicate requests)
	token, err := tm.requestNewToken()
	if err != nil {
		return "", err
	}

	return token, nil
}

// getCachedToken reads and checks the cached token for validity and populates in-memory cache.
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

	// valid -> set in-memory cache
	tm.memCache = &cached
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

	fmt.Println("🔐 Requesting token...")
	fmt.Println("🔗 URL:", url)
	fmt.Println("🧾 Auth:", "Basic "+credentials)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("📦 Raw Token Response (%d): %s\n", resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-200 response: %s", resp.Status)
	}

	var tokenResp tokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("token decode failed: %w - body: %s", err, string(body))
	}
	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("no access token returned. Body: %s", string(body))
	}

	expiresInInt, err := strconv.ParseInt(tokenResp.ExpiresIn, 10, 64)
	if err != nil {
		// try to be tolerant: if parse fails, default to 300 seconds
		fmt.Println("warning: invalid expires_in, defaulting to 300s:", err)
		expiresInInt = 300
	}

	// safe buffer handling: only subtract buffer when expires_in is larger than buffer
	var effectiveExpires int64
	const buffer = int64(60)
	if expiresInInt > buffer {
		effectiveExpires = expiresInInt - buffer
	} else {
		// avoid negative expiry; use half of the returned TTL or at least 1 second
		if expiresInInt > 1 {
			effectiveExpires = expiresInInt / 2
		} else {
			effectiveExpires = 1
		}
	}

	expiresAt := time.Now().Unix() + effectiveExpires

	// update memory cache first then persist
	tm.memCache = &tokenCache{
		Token:     tokenResp.AccessToken,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now().Unix(),
	}
	tm.cacheToken(tokenResp.AccessToken, expiresAt)

	return tokenResp.AccessToken, nil
}

// cacheToken writes token details to file atomically and logs errors
func (tm *TokenManager) cacheToken(token string, expiresAt int64) {
	cache := tokenCache{
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now().Unix(),
	}

	data, _ := json.Marshal(cache)

	// atomic write: write to temp file in same dir then rename
	dir := filepath.Dir(tm.CachePath)
	tmpf, err := os.CreateTemp(dir, "mpesa-token-*.tmp")
	if err != nil {
		fmt.Println("failed to create temp file for token cache:", err)
		// try non-atomic fallback
		_ = os.WriteFile(tm.CachePath, data, os.ModePerm)
		return
	}
	_, err = tmpf.Write(data)
	tmpf.Close()
	if err != nil {
		fmt.Println("failed writing token cache temp file:", err)
		_ = os.Remove(tmpf.Name())
		_ = os.WriteFile(tm.CachePath, data, os.ModePerm)
		return
	}
	_ = os.Chmod(tmpf.Name(), os.ModePerm)
	if err := os.Rename(tmpf.Name(), tm.CachePath); err != nil {
		fmt.Println("failed to rename token cache temp file:", err)
		// fallback
		_ = os.WriteFile(tm.CachePath, data, os.ModePerm)
	}
}

// ClearCache deletes the token cache file and resets in-memory cache
func (tm *TokenManager) ClearCache() {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.memCache = nil
	if _, err := os.Stat(tm.CachePath); err == nil {
		_ = os.Remove(tm.CachePath)
	}
}
