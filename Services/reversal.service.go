package Services

import (
	"errors"
	abstracts "github.com/venomous-maker/go-mpesa/Abstracts"
)

// ReversalService handles M-Pesa transaction reversal operations.
// This service allows businesses to reverse completed M-Pesa transactions when necessary,
// such as in cases of customer refunds or transaction errors.
type ReversalService struct {
	Config                 *abstracts.MpesaConfig   // M-Pesa configuration containing credentials and settings
	Client                 abstracts.MpesaInterface // HTTP client interface for making API requests
	Initiator              string                   // Username of the M-Pesa API operator
	TransactionID          string                   // ID of the transaction to be reversed
	ReceiverIdentifierType string                   // Type of identifier for the transaction receiver
	Remarks                string                   // Comments for the reversal transaction
	Occasion               string                   // Occasion or reason for the reversal
	Response               map[string]interface{}   // Response from the last API call
}

// NewReversalService creates a new reversal service instance with the provided configuration and client.
// This is the constructor for creating reversal service instances that can be used to reverse transactions.
//
// Parameters:
//   - cfg: M-Pesa configuration containing credentials and settings
//   - client: HTTP client interface for making API requests
//
// Returns:
//   - *ReversalService: A configured reversal service ready for transaction reversal operations
//
// Example:
//
//	cfg := createMpesaConfig()
//	client := Abstracts.NewApiClient(cfg)
//	reversalService := NewReversalService(cfg, client)
func NewReversalService(cfg *abstracts.MpesaConfig, client abstracts.MpesaInterface) *ReversalService {
	return &ReversalService{
		Config: cfg,
		Client: client,
	}
}

// SetInitiator sets the username of the M-Pesa API operator initiating the reversal.
// This is typically the username provided by Safaricom for transaction operations.
//
// Parameters:
//   - initiator: The initiator username/name
//
// Returns:
//   - *ReversalService: Returns self for method chaining
//
// Example:
//
//	reversalService.SetInitiator("testapi")
func (s *ReversalService) SetInitiator(initiator string) *ReversalService {
	s.Initiator = initiator
	return s
}

// SetTransactionID sets the ID of the transaction to be reversed.
// This should be the original transaction ID from the M-Pesa transaction that needs to be reversed.
//
// Parameters:
//   - txID: The transaction ID to reverse
//
// Returns:
//   - *ReversalService: Returns self for method chaining
//
// Example:
//
//	reversalService.SetTransactionID("OEI2AK4Q16")
//	reversalService.SetTransactionID("ABC123XYZ789")
func (s *ReversalService) SetTransactionID(txID string) *ReversalService {
	s.TransactionID = txID
	return s
}

// SetReceiverIdentifierType sets the type of identifier for the transaction receiver.
// This identifies the type of account that received the original transaction.
//
// Parameters:
//   - identifierType: The identifier type for the receiver
//
// Returns:
//   - *ReversalService: Returns self for method chaining
//
// Common Identifier Types:
//   - "1": MSISDN (phone number)
//   - "2": Till Number
//   - "4": Organization shortcode
//   - "11": Paybill
//
// Example:
//
//	reversalService.SetReceiverIdentifierType("1")  // For phone number
//	reversalService.SetReceiverIdentifierType("4")  // For shortcode
func (s *ReversalService) SetReceiverIdentifierType(identifierType string) *ReversalService {
	s.ReceiverIdentifierType = identifierType
	return s
}

// SetRemarks sets comments or additional information for the reversal transaction.
// This helps identify the reason for the reversal in transaction records.
//
// Parameters:
//   - remarks: A descriptive string for the reversal
//
// Returns:
//   - *ReversalService: Returns self for method chaining
//
// Example:
//
//	reversalService.SetRemarks("Customer refund requested")
//	reversalService.SetRemarks("Duplicate transaction reversal")
func (s *ReversalService) SetRemarks(remarks string) *ReversalService {
	s.Remarks = remarks
	return s
}

// SetOccasion sets the occasion or specific reason for the transaction reversal.
// This provides additional context for the reversal operation.
//
// Parameters:
//   - occasion: A string describing the occasion for the reversal
//
// Returns:
//   - *ReversalService: Returns self for method chaining
//
// Example:
//
//	reversalService.SetOccasion("Customer complaint")
//	reversalService.SetOccasion("System error correction")
func (s *ReversalService) SetOccasion(occasion string) *ReversalService {
	s.Occasion = occasion
	return s
}

// Reverse initiates the transaction reversal process.
// This method validates all required parameters and sends the reversal request to M-Pesa.
//
// Returns:
//   - map[string]interface{}: The response from the M-Pesa API
//   - error: An error if validation fails or the API request encounters issues
//
// Example:
//
//	response, err := reversalService.
//	    SetInitiator("testapi").
//	    SetTransactionID("OEI2AK4Q16").
//	    SetReceiverIdentifierType("1").
//	    SetRemarks("Customer refund").
//	    SetOccasion("Customer complaint").
//	    Reverse()
//	if err != nil {
//	    log.Printf("Transaction reversal failed: %v", err)
//	    return
//	}
//	fmt.Printf("Reversal initiated: %+v", response)
func (s *ReversalService) Reverse() (map[string]interface{}, error) {
	// Validate required fields
	if s.Initiator == "" {
		return nil, errors.New("initiator is required")
	}
	if s.TransactionID == "" {
		return nil, errors.New("transaction ID is required")
	}
	if s.ReceiverIdentifierType == "" {
		return nil, errors.New("receiver identifier type is required")
	}

	data := map[string]interface{}{
		"Initiator":              s.Initiator,
		"SecurityCredential":     s.Config.GetSecurityCredential(),
		"CommandID":              "TransactionReversal",
		"TransactionID":          s.TransactionID,
		"Amount":                 "1", // Amount is required but not used in reversals
		"ReceiverParty":          s.Config.GetBusinessCode(),
		"RecieverIdentifierType": s.ReceiverIdentifierType,
		"Remarks":                s.Remarks,
		"QueueTimeOutURL":        s.Config.GetQueueTimeoutURL(),
		"ResultURL":              s.Config.GetResultURL(),
		"Occasion":               s.Occasion,
	}

	response, err := s.Client.ExecuteRequest(data, "/mpesa/reversal/v1/request")
	if err != nil {
		return nil, err
	}

	s.Response = response
	return response, nil
}

// GetResponse returns the response from the last reversal operation.
//
// Returns:
//   - map[string]interface{}: The response data, or nil if no reversal has been made
func (s *ReversalService) GetResponse() map[string]interface{} {
	return s.Response
}
