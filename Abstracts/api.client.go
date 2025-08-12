package Abstracts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ApiClient handles HTTP communication with M-Pesa API endpoints.
// It manages authentication tokens and provides methods for making authenticated requests.
type ApiClient struct {
	Config       *MpesaConfig  // Configuration containing API credentials and settings
	TokenManager *TokenManager // Manager for handling OAuth tokens
}

// NewApiClient creates a new API client instance with the provided configuration.
// The client automatically manages OAuth tokens and handles request authentication.
//
// Parameters:
//   - config: M-Pesa configuration containing credentials and environment settings
//
// Returns:
//   - *ApiClient: A configured API client ready for making requests
//
// Example:
//
//	cfg := createMpesaConfig()
//	client := NewApiClient(cfg)
func NewApiClient(config *MpesaConfig) *ApiClient {
	return &ApiClient{
		Config:       config,
		TokenManager: NewTokenManager(config),
	}
}

// ExecuteRequest performs an authenticated POST request to the specified M-Pesa endpoint.
// This method automatically handles token acquisition, request formatting, and response parsing.
// It implements the MpesaInterface contract for making API calls.
//
// Parameters:
//   - payload: The request payload (typically a map[string]any with request data)
//   - endpoint: The API endpoint path (e.g., "/mpesa/stkpush/v1/processrequest")
//
// Returns:
//   - map[string]any: The parsed JSON response from the API
//   - error: An error if token acquisition, request execution, or response parsing fails
//
// Example:
//
//	data := map[string]any{
//	    "BusinessShortCode": "174379",
//	    "Amount": "100",
//	    "PhoneNumber": "254711223344",
//	    "CallBackURL": "https://example.com/callback",
//	}
//	response, err := client.ExecuteRequest(data, "/mpesa/stkpush/v1/processrequest")
//	if err != nil {
//	    log.Printf("Request failed: %v", err)
//	    return
//	}
func (client *ApiClient) ExecuteRequest(payload any, endpoint string) (map[string]any, error) {
	token, err := client.TokenManager.GetToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	return client.sendRequest(payload, endpoint, token, false)
}

// sendRequest performs the actual HTTP request with retry logic for token expiration.
// This internal method handles the low-level HTTP communication and automatic token refresh.
//
// Parameters:
//   - payload: The request payload to be JSON-encoded
//   - endpoint: The API endpoint path
//   - token: The OAuth bearer token for authentication
//   - isRetry: Flag indicating if this is a retry attempt after token refresh
//
// Returns:
//   - map[string]any: The parsed JSON response from the API
//   - error: An error if the request fails or response parsing fails
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

		// Handle "The transaction is being processed" specially (for STK Push Query)
		if msg == "The transaction is being processed" {
			// Return response without error to allow caller to handle this state
			return response, nil
		}

		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, msg)
	}

	return response, nil
}
