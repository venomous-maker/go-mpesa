package Services

import (
	"errors"
	"github.com/venomous-maker/go-mpesa/Abstracts"
)

// TransactionStatusService handles transaction status inquiry operations.
// This service allows businesses to check the status of any M-Pesa transaction
// using the transaction ID to get detailed information about the transaction.
type TransactionStatusService struct {
	*AbstractService

	initiator      string // Username of the M-Pesa API operator
	transactionID  string // ID of the transaction to check status for
	identifierType string // Type of organization checking the transaction
	remarks        string // Comments for the status inquiry
	occasion       string // Occasion or reason for the status check
}

// NewTransactionStatusService creates a new transaction status service instance with the provided configuration and client.
// This is the constructor for creating transaction status service instances that can be used to query transaction status.
//
// Parameters:
//   - cfg: M-Pesa configuration containing credentials and settings
//   - client: HTTP client interface for making API requests
//
// Returns:
//   - *TransactionStatusService: A configured transaction status service ready for status operations
//
// Example:
//
//	cfg := createMpesaConfig()
//	client := Abstracts.NewApiClient(cfg)
//	statusService := NewTransactionStatusService(cfg, client)
func NewTransactionStatusService(cfg *Abstracts.MpesaConfig, client Abstracts.MpesaInterface) *TransactionStatusService {
	return &TransactionStatusService{
		AbstractService: NewAbstractService(cfg, client),
	}
}

// SetInitiator sets the username of the M-Pesa API operator initiating the status inquiry.
// This is typically the username provided by Safaricom for API operations.
//
// Parameters:
//   - initiator: The initiator username/name
//
// Returns:
//   - *TransactionStatusService: Returns self for method chaining
//
// Example:
//
//	statusService.SetInitiator("testapi")
func (s *TransactionStatusService) SetInitiator(initiator string) *TransactionStatusService {
	s.initiator = initiator
	return s
}

// SetTransactionID sets the ID of the transaction to check the status for.
// This should be the original transaction ID from the M-Pesa transaction.
//
// Parameters:
//   - id: The transaction ID to check status for
//
// Returns:
//   - *TransactionStatusService: Returns self for method chaining
//
// Example:
//
//	statusService.SetTransactionID("OEI2AK4Q16")
//	statusService.SetTransactionID("ABC123XYZ789")
func (s *TransactionStatusService) SetTransactionID(id string) *TransactionStatusService {
	s.transactionID = id
	return s
}

// SetIdentifierType sets the type of organization checking the transaction status.
// This identifies the type of shortcode making the inquiry.
//
// Parameters:
//   - idType: The identifier type for the organization
//
// Returns:
//   - *TransactionStatusService: Returns self for method chaining
//
// Common Identifier Types:
//   - "1": MSISDN (phone number)
//   - "2": Till Number
//   - "4": Organization shortcode
//   - "11": Paybill
//
// Example:
//
//	statusService.SetIdentifierType("4")  // For organization shortcode
func (s *TransactionStatusService) SetIdentifierType(idType string) *TransactionStatusService {
	s.identifierType = idType
	return s
}

// SetRemarks sets comments or additional information for the status inquiry.
// This helps identify the purpose of the status check in transaction records.
//
// Parameters:
//   - remarks: A descriptive string for the status inquiry
//
// Returns:
//   - *TransactionStatusService: Returns self for method chaining
//
// Example:
//
//	statusService.SetRemarks("Status check for customer inquiry")
//	statusService.SetRemarks("Reconciliation status verification")
func (s *TransactionStatusService) SetRemarks(remarks string) *TransactionStatusService {
	s.remarks = remarks
	return s
}

// SetOccasion sets the occasion or specific reason for the status inquiry.
// This provides additional context for the status check operation.
//
// Parameters:
//   - occasion: A string describing the occasion for the status inquiry
//
// Returns:
//   - *TransactionStatusService: Returns self for method chaining
//
// Example:
//
//	statusService.SetOccasion("Customer complaint investigation")
//	statusService.SetOccasion("Daily reconciliation process")
func (s *TransactionStatusService) SetOccasion(occasion string) *TransactionStatusService {
	s.occasion = occasion
	return s
}

// Query initiates a transaction status inquiry to check the current status of a transaction.
// This method validates all required parameters and sends the status request to M-Pesa.
//
// Returns:
//   - map[string]any: The response from the M-Pesa API containing transaction status information
//   - error: An error if validation fails or the API request encounters issues
//
// Example:
//
//	response, err := statusService.
//	    SetInitiator("testapi").
//	    SetTransactionID("OEI2AK4Q16").
//	    SetIdentifierType("4").
//	    SetRemarks("Status inquiry").
//	    SetOccasion("Customer inquiry").
//	    Query()
//	if err != nil {
//	    log.Printf("Transaction status inquiry failed: %v", err)
//	    return
//	}
//	fmt.Printf("Transaction status: %+v", response)
func (s *TransactionStatusService) Query() (map[string]any, error) {
	// Validate required fields
	if s.initiator == "" {
		return nil, errors.New("initiator is required")
	}
	if s.transactionID == "" {
		return nil, errors.New("transaction ID is required")
	}
	if s.identifierType == "" {
		return nil, errors.New("identifier type is required")
	}

	data := map[string]any{
		"Initiator":          s.initiator,
		"SecurityCredential": s.Config.GetSecurityCredential(),
		"CommandID":          "TransactionStatusQuery",
		"TransactionID":      s.transactionID,
		"PartyA":             s.Config.GetBusinessCode(),
		"IdentifierType":     s.identifierType,
		"Remarks":            s.remarks,
		"QueueTimeOutURL":    s.Config.GetQueueTimeoutURL(),
		"ResultURL":          s.Config.GetResultURL(),
		"Occasion":           s.occasion,
	}

	response, err := s.Client.ExecuteRequest(data, "/mpesa/transactionstatus/v1/query")
	if err != nil {
		return nil, err
	}

	s.setResponse(response)
	return response, nil
}
