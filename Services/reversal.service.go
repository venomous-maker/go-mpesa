package Services

import (
	"errors"
	abstracts "github.com/venomous-maker/go-mpesa/Abstracts"
	"strconv" // added for int to string conversion of amount
)

// ReversalService handles M-Pesa transaction reversal operations.
// This service allows businesses to reverse completed M-Pesa transactions when necessary,
// such as in cases of customer refunds or transaction errors.
type ReversalService struct {
	Config                 *abstracts.MpesaConfig   // M-Pesa configuration containing credentials and settings
	Client                 abstracts.MpesaInterface // HTTP client interface for making API requests
	Initiator              string                   // Username of the M-Pesa API operator
	TransactionID          string                   // ID of the transaction to be reversed
	Amount                 int                      // Original transaction amount to reverse
	ReceiverIdentifierType string                   // Type of identifier for the transaction receiver (e.g. 11 for Paybill)
	Remarks                string                   // Comments for the reversal transaction (2-100 chars, required)
	Occasion               string                   // Occasion or reason for the reversal (optional)
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
func (s *ReversalService) SetTransactionID(txID string) *ReversalService {
	s.TransactionID = txID
	return s
}

// SetAmount sets the original transaction amount that is being reversed (required).
//
// Parameters:
//   - amount: The numeric amount of the original transaction (must be > 0)
//
// Returns:
//   - *ReversalService: Returns self for method chaining
func (s *ReversalService) SetAmount(amount int) *ReversalService {
	s.Amount = amount
	return s
}

// SetReceiverIdentifierType sets the type of identifier for the transaction receiver.
// This identifies the type of account that received the original transaction.
// For reversals, Safaricom docs indicate Paybill reversals use identifier type "11".
//
// Parameters:
//   - identifierType: The identifier type for the receiver
//
// Returns:
//   - *ReversalService: Returns self for method chaining
func (s *ReversalService) SetReceiverIdentifierType(identifierType string) *ReversalService {
	s.ReceiverIdentifierType = identifierType
	return s
}

// SetRemarks sets comments or additional information for the reversal transaction.
// This helps identify the reason for the reversal in transaction records.
//
// Parameters:
//   - remarks: A descriptive string for the reversal (2-100 characters)
//
// Returns:
//   - *ReversalService: Returns self for method chaining
func (s *ReversalService) SetRemarks(remarks string) *ReversalService {
	s.Remarks = remarks
	return s
}

// SetOccasion sets the occasion or specific reason for the transaction reversal.
// This provides additional context for the reversal operation (optional).
//
// Parameters:
//   - occasion: A string describing the occasion for the reversal
//
// Returns:
//   - *ReversalService: Returns self for method chaining
func (s *ReversalService) SetOccasion(occasion string) *ReversalService {
	s.Occasion = occasion
	return s
}

// Reverse initiates the transaction reversal process.
// This method validates all required parameters and sends the reversal request to M-Pesa.
// Required fields (per Safaricom docs): Initiator, SecurityCredential, CommandID (TransactionReversal),
// TransactionID, Amount, ReceiverParty (Business Short Code), RecieverIdentifierType, Remarks,
// QueueTimeOutURL, ResultURL.
//
// Returns:
//   - map[string]interface{}: The response from the M-Pesa API
//   - error: An error if validation fails or the API request encounters issues
func (s *ReversalService) Reverse() (map[string]interface{}, error) {
	// Validate required fields
	if s.Initiator == "" {
		return nil, errors.New("initiator is required")
	}
	if s.TransactionID == "" {
		return nil, errors.New("transaction ID is required")
	}
	if s.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	if s.ReceiverIdentifierType == "" {
		return nil, errors.New("receiver identifier type is required")
	}
	if s.Remarks == "" {
		return nil, errors.New("remarks are required")
	}
	if s.Config.GetBusinessCode() == "" {
		return nil, errors.New("business shortcode (ReceiverParty) is required; call SetBusinessCode on mpesa config")
	}
	if s.Config.GetQueueTimeoutURL() == "" {
		return nil, errors.New("queue timeout URL is required; call SetQueueTimeoutURL on config")
	}
	if s.Config.GetResultURL() == "" {
		return nil, errors.New("result URL is required; call SetResultURL on config")
	}
	if s.Config.GetSecurityCredential() == "" {
		return nil, errors.New("security credential is required; set via SetSecurityCredential or OverrideSecurityCredential on config")
	}

	data := map[string]interface{}{
		"Initiator":              s.Initiator,
		"SecurityCredential":     s.Config.GetSecurityCredential(),
		"CommandID":              "TransactionReversal",
		"TransactionID":          s.TransactionID,
		"Amount":                 strconv.Itoa(s.Amount),
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
