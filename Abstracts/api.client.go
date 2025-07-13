package Abstracts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ApiClient struct {
	Config       *MpesaConfig
	TokenManager *TokenManager
}

// NewApiClient initializes an API client using config and token manager
func NewApiClient(config *MpesaConfig) *ApiClient {
	return &ApiClient{
		Config:       config,
		TokenManager: NewTokenManager(config),
	}
}

// ExecuteRequest performs an authenticated POST request to an M-Pesa endpoint
func (client *ApiClient) ExecuteRequest(payload any, endpoint string) (map[string]any, error) {
	token, err := client.TokenManager.GetToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	return client.sendRequest(payload, endpoint, token, false)
}

// sendRequest performs the actual HTTP request logic and retries if needed
func (client *ApiClient) sendRequest(payload any, endpoint, token string, isRetry bool) (map[string]any, error) {
	url := client.Config.GetBaseURL() + endpoint

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("json encode error: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 401 && !isRetry {
		// Clear token and retry once
		client.TokenManager.ClearCache()
		newToken, err := client.TokenManager.GetToken()
		if err != nil {
			return nil, fmt.Errorf("token refresh failed: %w", err)
		}
		return client.sendRequest(payload, endpoint, newToken, true)
	}

	var response map[string]any
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("response decode error: %w", err)
	}

	if resp.StatusCode >= 400 {
		msg := "Unknown error"
		if val, ok := response["errorMessage"]; ok {
			msg = fmt.Sprint(val)
		}
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, msg)
	}

	return response, nil
}
