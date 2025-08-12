package Services

import "github.com/venomous-maker/go-mpesa/Abstracts"

// AbstractService provides a generic service foundation with response management capabilities.
// This service extends BaseService to add response handling functionality that can be
// shared across different M-Pesa service implementations.
type AbstractService struct {
	*BaseService                // Embedded base service with common M-Pesa functionality
	response     map[string]any // Storage for the last API response
}

// NewAbstractService creates a new abstract service instance with the provided configuration and client.
// This service provides basic response management functionality for other services to extend.
//
// Parameters:
//   - cfg: M-Pesa configuration containing credentials and settings
//   - client: HTTP client interface for making API requests
//
// Returns:
//   - *AbstractService: A configured abstract service with response management capabilities
//
// Example:
//
//	cfg := createMpesaConfig()
//	client := Abstracts.NewApiClient(cfg)
//	abstractService := NewAbstractService(cfg, client)
func NewAbstractService(cfg *Abstracts.MpesaConfig, client Abstracts.MpesaInterface) *AbstractService {
	return &AbstractService{
		BaseService: NewBaseService(cfg, client),
	}
}

// GetResponse returns the stored response from the last API operation.
// This method provides access to the complete API response for advanced use cases
// where you need to access fields not covered by specific getter methods.
//
// Returns:
//   - map[string]any: The complete response map from the last API call, or nil if no call has been made
//
// Example:
//
//	response := abstractService.GetResponse()
//	if response != nil {
//	    if errorCode, ok := response["errorCode"].(string); ok {
//	        fmt.Printf("Error Code: %s", errorCode)
//	    }
//	    if message, ok := response["errorMessage"].(string); ok {
//	        fmt.Printf("Message: %s", message)
//	    }
//	}
func (s *AbstractService) GetResponse() map[string]any {
	return s.response
}

// setResponse is a helper method to store the API response for later retrieval.
// This internal method is used by service implementations to preserve response data
// after API calls for debugging and advanced processing.
//
// Parameters:
//   - resp: The response map to store
func (s *AbstractService) setResponse(resp map[string]any) {
	s.response = resp
}
