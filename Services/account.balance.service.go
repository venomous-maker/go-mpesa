package Services

import (
	"errors"
	abstracts "github.com/venomous-maker/go-mpesa/Abstracts"
)

// AccountBalanceService handles account balance inquiry operations.
// This service allows businesses to check their M-Pesa account balance programmatically.
type AccountBalanceService struct {
	Config         *abstracts.MpesaConfig   // M-Pesa configuration containing credentials and settings
	Client         abstracts.MpesaInterface // HTTP client interface for making API requests
	initiator      string                   // Username of the M-Pesa API operator
	identifierType string                   // Type of organization receiving the transaction
	remarks        string                   // Comments that are sent along with the transaction
}

// NewAccountBalanceService creates a new account balance service instance with the provided configuration and client.
// This is the constructor for creating account balance service instances that can be used to query account balances.
//
// Parameters:
//   - cfg: M-Pesa configuration containing credentials and settings
//   - client: HTTP client interface for making API requests
//
// Returns:
//   - *AccountBalanceService: A configured account balance service ready for balance operations
//
// Example:
//
//	cfg := createMpesaConfig()
//	client := Abstracts.NewApiClient(cfg)
//	balanceService := NewAccountBalanceService(cfg, client)
func NewAccountBalanceService(cfg *abstracts.MpesaConfig, client abstracts.MpesaInterface) *AccountBalanceService {
	return &AccountBalanceService{
		Config: cfg,
		Client: client,
	}
}

// SetInitiator sets the username of the M-Pesa API operator initiating the balance inquiry.
// This is typically the username provided by Safaricom for API operations.
//
// Parameters:
//   - initiator: The initiator username/name
//
// Returns:
//   - *AccountBalanceService: Returns self for method chaining
//
// Example:
//
//	balanceService.SetInitiator("testapi")
func (s *AccountBalanceService) SetInitiator(initiator string) *AccountBalanceService {
	s.initiator = initiator
	return s
}

// SetIdentifierType sets the type of organization receiving the transaction.
// This identifies the type of shortcode being queried for balance.
//
// Parameters:
//   - identifierType: The identifier type (e.g., "4" for organization shortcode)
//
// Returns:
//   - *AccountBalanceService: Returns self for method chaining
//
// Common Identifier Types:
//   - "1": MSISDN
//   - "2": Till Number
//   - "4": Organization shortcode
//
// Example:
//
//	balanceService.SetIdentifierType("4")  // For organization shortcode
func (s *AccountBalanceService) SetIdentifierType(identifierType string) *AccountBalanceService {
	s.identifierType = identifierType
	return s
}

// SetRemarks sets additional information to be associated with the balance inquiry.
// This helps identify the purpose of the balance check in transaction records.
//
// Parameters:
//   - remarks: A descriptive string for the balance inquiry
//
// Returns:
//   - *AccountBalanceService: Returns self for method chaining
//
// Example:
//
//	balanceService.SetRemarks("Daily balance check")
//	balanceService.SetRemarks("End of month reconciliation")
func (s *AccountBalanceService) SetRemarks(remarks string) *AccountBalanceService {
	s.remarks = remarks
	return s
}

// Query initiates an account balance inquiry to check the current account balance.
// This method validates all required parameters and sends the balance request to M-Pesa.
//
// Returns:
//   - map[string]any: The response from the M-Pesa API containing balance information
//   - error: An error if validation fails or the API request encounters issues
//
// Example:
//
//	response, err := balanceService.
//	    SetInitiator("testapi").
//	    SetIdentifierType("4").
//	    SetRemarks("Balance inquiry").
//	    Query()
//	if err != nil {
//	    log.Printf("Balance inquiry failed: %v", err)
//	    return
//	}
//	fmt.Printf("Balance response: %+v", response)
func (s *AccountBalanceService) Query() (map[string]any, error) {
	// Validate required fields
	if s.initiator == "" {
		return nil, errors.New("initiator is required")
	}
	if s.identifierType == "" {
		return nil, errors.New("identifier type is required")
	}

	data := map[string]any{
		"Initiator":          s.initiator,
		"SecurityCredential": s.Config.GetSecurityCredential(),
		"CommandID":          "AccountBalance",
		"PartyA":             s.Config.GetBusinessCode(),
		"IdentifierType":     s.identifierType,
		"Remarks":            s.remarks,
		"QueueTimeOutURL":    s.Config.GetQueueTimeoutURL(),
		"ResultURL":          s.Config.GetResultURL(),
	}

	return s.Client.ExecuteRequest(data, "/mpesa/accountbalance/v1/query")
}
