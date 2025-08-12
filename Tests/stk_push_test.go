package Tests

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/venomous-maker/go-mpesa/Abstracts"
	"github.com/venomous-maker/go-mpesa/Services"
)

// MockMpesaInterface is a mock implementation of MpesaInterface
type MockMpesaInterface struct {
	mock.Mock
}

func (m *MockMpesaInterface) ExecuteRequest(payload any, endpoint string) (map[string]any, error) {
	args := m.Called(payload, endpoint)
	return args.Get(0).(map[string]any), args.Error(1)
}

func (m *MockMpesaInterface) GetAccessToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

// Helper function to create test config
func createTestConfig() *Abstracts.MpesaConfig {
	cfg, _ := Abstracts.NewMpesaConfig(
		"test_consumer_key",
		"test_consumer_secret",
		Abstracts.Sandbox,
		ptr("174379"),
		ptr("test_passkey"),
		nil, nil, nil,
	)
	return cfg
}

func ptr(s string) *string {
	return &s
}

func TestNewStkService(t *testing.T) {
	cfg := createTestConfig()
	mockClient := &MockMpesaInterface{}

	service := Services.NewStkService(cfg, mockClient)

	assert.NotNil(t, service)
	assert.Equal(t, cfg, service.Config)
	assert.Equal(t, mockClient, service.Client)
}

func TestStkService_SetTransactionType(t *testing.T) {
	cfg := createTestConfig()
	mockClient := &MockMpesaInterface{}
	service := Services.NewStkService(cfg, mockClient)

	result := service.SetTransactionType("CustomerPayBillOnline")

	assert.Equal(t, service, result) // Should return self for chaining
}

func TestStkService_SetAmount(t *testing.T) {
	cfg := createTestConfig()
	mockClient := &MockMpesaInterface{}
	service := Services.NewStkService(cfg, mockClient)

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"Integer amount", 100, "100"},
		{"String amount", "250", "250"},
		{"Int64 amount", int64(500), "500"},
		{"Float amount", 99.99, "99.99"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.SetAmount(tt.input)
			assert.Equal(t, service, result)
		})
	}
}

func TestStkService_SetPhoneNumber(t *testing.T) {
	cfg := createTestConfig()
	mockClient := &MockMpesaInterface{}
	service := Services.NewStkService(cfg, mockClient)

	tests := []struct {
		name        string
		phoneNumber string
		expectError bool
	}{
		{"Valid phone number with country code", "254111844429", false},
		{"Valid phone number without country code", "0111844429", false},
		{"Empty phone number", "", true},
		{"Short phone number", "123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.SetPhoneNumber(tt.phoneNumber)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, service, result)
			}
		})
	}
}

func TestStkService_SetCallbackUrl(t *testing.T) {
	cfg := createTestConfig()
	mockClient := &MockMpesaInterface{}
	service := Services.NewStkService(cfg, mockClient)

	url := "https://example.com/callback"
	result := service.SetCallbackUrl(url)

	assert.Equal(t, service, result)
}

func TestStkService_SetAccountReference(t *testing.T) {
	cfg := createTestConfig()
	mockClient := &MockMpesaInterface{}
	service := Services.NewStkService(cfg, mockClient)

	ref := "TEST_REF_123"
	result := service.SetAccountReference(ref)

	assert.Equal(t, service, result)
}

func TestStkService_SetTransactionDesc(t *testing.T) {
	cfg := createTestConfig()
	mockClient := &MockMpesaInterface{}
	service := Services.NewStkService(cfg, mockClient)

	desc := "Test transaction description"
	result := service.SetTransactionDesc(desc)

	assert.Equal(t, service, result)
}

func TestStkService_Push_Success(t *testing.T) {
	cfg := createTestConfig()
	mockClient := &MockMpesaInterface{}
	service := Services.NewStkService(cfg, mockClient)

	// Setup mock response
	mockResponse := map[string]any{
		"CheckoutRequestID":   "ws_CO_12345678",
		"ResponseCode":        "0",
		"ResponseDescription": "Success",
		"MerchantRequestID":   "29115-34620561-1",
	}

	mockClient.On("ExecuteRequest", mock.AnythingOfType("map[string]interface {}"), "/mpesa/stkpush/v1/processrequest").Return(mockResponse, nil)

	// Configure service
	service.SetTransactionType("CustomerPayBillOnline")
	service.SetAmount("100")
	service, _ = service.SetPhoneNumber("254111844429")
	service.SetCallbackUrl("https://example.com/callback")
	service.SetAccountReference("TEST_REF")
	service.SetTransactionDesc("Test transaction")

	// Execute push
	result, err := service.Push()

	assert.NoError(t, err)
	assert.Equal(t, service, result)
	mockClient.AssertExpectations(t)
}

func TestStkService_Push_ValidationErrors(t *testing.T) {
	cfg := createTestConfig()
	mockClient := &MockMpesaInterface{}

	tests := []struct {
		name          string
		setupService  func(*Services.StkService) *Services.StkService
		expectedError string
	}{
		{
			name: "Missing transaction type",
			setupService: func(s *Services.StkService) *Services.StkService {
				s.SetAmount("100")
				s, _ = s.SetPhoneNumber("254111844429")
				s.SetCallbackUrl("https://example.com/callback")
				return s
			},
			expectedError: "transaction type is required",
		},
		{
			name: "Missing amount",
			setupService: func(s *Services.StkService) *Services.StkService {
				s.SetTransactionType("CustomerPayBillOnline")
				s, _ = s.SetPhoneNumber("254111844429")
				s.SetCallbackUrl("https://example.com/callback")
				return s
			},
			expectedError: "amount is required",
		},
		{
			name: "Missing phone number",
			setupService: func(s *Services.StkService) *Services.StkService {
				s.SetTransactionType("CustomerPayBillOnline")
				s.SetAmount("100")
				s.SetCallbackUrl("https://example.com/callback")
				return s
			},
			expectedError: "phone number is required",
		},
		{
			name: "Missing callback URL",
			setupService: func(s *Services.StkService) *Services.StkService {
				s.SetTransactionType("CustomerPayBillOnline")
				s.SetAmount("100")
				s, _ = s.SetPhoneNumber("254111844429")
				return s
			},
			expectedError: "callback URL is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := Services.NewStkService(cfg, mockClient)
			service = tt.setupService(service)

			_, err := service.Push()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

func TestStkService_Push_APIError(t *testing.T) {
	cfg := createTestConfig()
	mockClient := &MockMpesaInterface{}
	service := Services.NewStkService(cfg, mockClient)

	// Setup mock to return error
	mockClient.On("ExecuteRequest", mock.AnythingOfType("map[string]interface {}"), "/mpesa/stkpush/v1/processrequest").Return(map[string]any{}, errors.New("API error"))

	// Configure service
	service.SetTransactionType("CustomerPayBillOnline")
	service.SetAmount("100")
	service, _ = service.SetPhoneNumber("254111844429")
	service.SetCallbackUrl("https://example.com/callback")

	// Execute push
	_, err := service.Push()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API error")
	mockClient.AssertExpectations(t)
}

func TestStkService_GetCheckoutRequestID_Success(t *testing.T) {
	cfg := createTestConfig()
	mockClient := &MockMpesaInterface{}
	service := Services.NewStkService(cfg, mockClient)

	// Setup mock response
	mockResponse := map[string]any{
		"CheckoutRequestID": "ws_CO_12345678",
		"ResponseCode":      "0",
	}

	mockClient.On("ExecuteRequest", mock.AnythingOfType("map[string]interface {}"), "/mpesa/stkpush/v1/processrequest").Return(mockResponse, nil)

	// Configure and execute push first
	service.SetTransactionType("CustomerPayBillOnline")
	service.SetAmount("100")
	service, _ = service.SetPhoneNumber("254111844429")
	service.SetCallbackUrl("https://example.com/callback")
	_, err := service.Push()
	assert.NoError(t, err)

	// Get checkout request ID
	id, err := service.GetCheckoutRequestID()

	assert.NoError(t, err)
	assert.Equal(t, "ws_CO_12345678", id)
}

func TestStkService_GetCheckoutRequestID_NoResponse(t *testing.T) {
	cfg := createTestConfig()
	mockClient := &MockMpesaInterface{}
	service := Services.NewStkService(cfg, mockClient)

	// Try to get checkout request ID without pushing first
	_, err := service.GetCheckoutRequestID()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no STK push response available")
}

func TestStkService_Query_Success(t *testing.T) {
	cfg := createTestConfig()
	mockClient := &MockMpesaInterface{}
	service := Services.NewStkService(cfg, mockClient)

	// Setup mock response for query
	mockQueryResponse := map[string]any{
		"ResponseCode":        "0",
		"ResponseDescription": "The service request has been accepted successfully",
		"MerchantRequestID":   "29115-34620561-1",
		"CheckoutRequestID":   "ws_CO_12345678",
		"ResultCode":          "0",
		"ResultDesc":          "The service request is processed successfully.",
	}

	mockClient.On("ExecuteRequest", mock.AnythingOfType("map[string]interface {}"), "/mpesa/stkpushquery/v1/query").Return(mockQueryResponse, nil)

	// Execute query
	result, err := service.Query("ws_CO_12345678")

	assert.NoError(t, err)
	assert.Equal(t, mockQueryResponse, result)
	mockClient.AssertExpectations(t)
}

func TestStkService_Query_APIError(t *testing.T) {
	cfg := createTestConfig()
	mockClient := &MockMpesaInterface{}
	service := Services.NewStkService(cfg, mockClient)

	// Setup mock to return error
	mockClient.On("ExecuteRequest", mock.AnythingOfType("map[string]interface {}"), "/mpesa/stkpushquery/v1/query").Return(map[string]any{}, errors.New("Query API error"))

	// Execute query
	_, err := service.Query("ws_CO_12345678")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Query API error")
	mockClient.AssertExpectations(t)
}

// Integration test example (would require actual API credentials)
func TestStkService_Integration(t *testing.T) {
	t.Skip("Skipping integration test - requires valid API credentials")

	// This test would use real credentials and test against the sandbox environment
	cfg, err := Abstracts.NewMpesaConfig(
		"your_consumer_key",
		"your_consumer_secret",
		Abstracts.Sandbox,
		ptr("174379"),
		ptr("your_passkey"),
		nil, nil, nil,
	)
	assert.NoError(t, err)

	client := Abstracts.NewApiClient(cfg)
	service := Services.NewStkService(cfg, client)

	// Configure request
	service.SetTransactionType("CustomerPayBillOnline")
	service.SetAmount("1")
	service, err = service.SetPhoneNumber("254111844429")
	assert.NoError(t, err)
	service.SetCallbackUrl("https://example.com/callback")
	service.SetAccountReference("TEST_REF")
	service.SetTransactionDesc("Integration test")

	// Execute push
	_, err = service.Push()
	assert.NoError(t, err)

	// Get checkout request ID
	id, err := service.GetCheckoutRequestID()
	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	// Query status
	status, err := service.Query(id)
	assert.NoError(t, err)
	assert.NotNil(t, status)
}
