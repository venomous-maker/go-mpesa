package Abstracts

// MpesaInterface defines the contract for executing M-Pesa API requests.
// This interface abstracts the HTTP client functionality, allowing for easy testing
// and different implementations of the API client.
type MpesaInterface interface {
	// ExecuteRequest sends an HTTP request to the specified M-Pesa API endpoint.
	// This method handles authentication, request formatting, and response parsing.
	//
	// Parameters:
	//   - payload: The request payload (typically a map[string]any with request data)
	//   - endpoint: The API endpoint path (e.g., "/mpesa/stkpush/v1/processrequest")
	//
	// Returns:
	//   - map[string]any: The parsed JSON response from the API
	//   - error: An error if the request fails or response parsing fails
	//
	// Example:
	//
	//	data := map[string]any{
	//	    "BusinessShortCode": "174379",
	//	    "Amount": "100",
	//	    "PhoneNumber": "254711223344",
	//	}
	//	response, err := client.ExecuteRequest(data, "/mpesa/stkpush/v1/processrequest")
	//	if err != nil {
	//	    log.Printf("Request failed: %v", err)
	//	    return
	//	}
	ExecuteRequest(payload any, endpoint string) (map[string]any, error)
}
