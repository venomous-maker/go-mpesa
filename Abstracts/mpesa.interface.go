package Abstracts

// MpesaInterface defines the contract for executing API requests
type MpesaInterface interface {
	ExecuteRequest(payload any, endpoint string) (map[string]any, error)
}
