// Package Services provides M-Pesa API service implementations for various operations
// including STK Push, B2C, C2B, Account Balance, Transaction Status, and Reversals.
package Services

import (
	"errors"
	"fmt"
	"github.com/venomous-maker/go-mpesa/Abstracts"
	"strconv"
)

// StkService handles STK Push (Lipa na M-Pesa Online) operations.
// STK Push allows initiating M-Pesa payments directly from a customer's mobile phone.
type StkService struct {
	*BaseService

	transactionType  string // The type of transaction (e.g., "CustomerPayBillOnline")
	amount           string // The amount to be charged from the customer
	phoneNumber      string // The customer's mobile phone number
	callbackUrl      string // URL to receive payment notifications
	accountReference string // Reference for the account being paid
	transactionDesc  string // Description of the transaction

	response map[string]any // Response from the last STK push request
}

// NewStkService creates a new STK Push service instance with the provided configuration and client.
// This is the constructor for creating STK service instances that can be used to initiate payments.
//
// Parameters:
//   - cfg: M-Pesa configuration containing credentials and settings
//   - client: HTTP client interface for making API requests
//
// Returns:
//   - *StkService: A configured STK service ready for payment operations
//
// Example:
//
//	cfg := createMpesaConfig()
//	client := Abstracts.NewApiClient(cfg)
//	stkService := NewStkService(cfg, client)
func NewStkService(cfg *Abstracts.MpesaConfig, client Abstracts.MpesaInterface) *StkService {
	return &StkService{
		BaseService: NewBaseService(cfg, client),
	}
}

// SetTransactionType sets the type of STK Push transaction.
// Common transaction types include "CustomerPayBillOnline" for pay bill transactions
// and "CustomerBuyGoodsOnline" for buy goods transactions.
//
// Parameters:
//   - t: The transaction type string
//
// Returns:
//   - *StkService: Returns self for method chaining
//
// Example:
//
//	stkService.SetTransactionType("CustomerPayBillOnline")
//	stkService.SetTransactionType("CustomerBuyGoodsOnline")
func (s *StkService) SetTransactionType(t string) *StkService {
	s.transactionType = t
	return s
}

// SetAmount sets the amount to be charged from the customer's M-Pesa account.
// The method accepts various numeric types and converts them to the required string format.
//
// Parameters:
//   - a: The amount as int, int64, string, float64, or any other type
//
// Returns:
//   - *StkService: Returns self for method chaining
//
// Example:
//
//	stkService.SetAmount(100)        // int
//	stkService.SetAmount("250.50")   // string
//	stkService.SetAmount(99.99)      // float64
//	stkService.SetAmount(int64(500)) // int64
func (s *StkService) SetAmount(a any) *StkService {
	switch v := a.(type) {
	case int:
		s.amount = strconv.Itoa(v)
	case int64:
		s.amount = strconv.FormatInt(v, 10)
	case string:
		s.amount = v
	default:
		s.amount = fmt.Sprint(v)
	}
	return s
}

// SetPhoneNumber sets and validates the customer's phone number for the STK Push.
// The method automatically formats the phone number to the correct international format
// and validates its length and format.
//
// Parameters:
//   - phone: The customer's phone number in various formats (254711223344, 0711223344, +254711223344)
//
// Returns:
//   - *StkService: Returns self for method chaining
//   - error: An error if the phone number is invalid or improperly formatted
//
// Example:
//
//	stkService, err := stkService.SetPhoneNumber("254711223344")  // With country code
//	stkService, err := stkService.SetPhoneNumber("0711223344")    // Without country code
//	stkService, err := stkService.SetPhoneNumber("+254711223344") // International format
func (s *StkService) SetPhoneNumber(phone string) (*StkService, error) {
	cleaned, err := s.CleanPhoneNumber(phone, "254")
	if err != nil {
		return s, err
	}
	s.phoneNumber = cleaned
	return s, nil
}

// SetCallbackUrl sets the URL where M-Pesa will send payment notifications.
// This URL will receive POST requests with the payment status and details.
// The callback URL must be publicly accessible and properly configured to handle M-Pesa callbacks.
//
// Parameters:
//   - url: The fully qualified callback URL (must be HTTPS in production)
//
// Returns:
//   - *StkService: Returns self for method chaining
//
// Example:
//
//	stkService.SetCallbackUrl("https://yourdomain.com/mpesa/callback")
//	stkService.SetCallbackUrl("https://api.example.com/webhooks/mpesa")
func (s *StkService) SetCallbackUrl(url string) *StkService {
	s.callbackUrl = url
	return s
}

// SetAccountReference sets the account reference for the transaction.
// This is typically used to identify the specific account, order, or invoice being paid.
// If not provided, a default value of "Account" will be used.
//
// Parameters:
//   - ref: The account reference string (e.g., order ID, invoice number, account number)
//
// Returns:
//   - *StkService: Returns self for method chaining
//
// Example:
//
//	stkService.SetAccountReference("ORDER123")
//	stkService.SetAccountReference("INV-2024-001")
//	stkService.SetAccountReference("ACCOUNT456789")
func (s *StkService) SetAccountReference(ref string) *StkService {
	s.accountReference = ref
	return s
}

// SetTransactionDesc sets a human-readable description for the transaction.
// This description helps identify the purpose of the payment and may be visible to the customer.
// If not provided, a default value of "Transaction" will be used.
//
// Parameters:
//   - desc: A descriptive string explaining the transaction purpose
//
// Returns:
//   - *StkService: Returns self for method chaining
//
// Example:
//
//	stkService.SetTransactionDesc("Payment for monthly subscription")
//	stkService.SetTransactionDesc("Purchase of product XYZ")
//	stkService.SetTransactionDesc("Invoice payment for services")
func (s *StkService) SetTransactionDesc(desc string) *StkService {
	s.transactionDesc = desc
	return s
}

// validatePushParams validates that all required parameters are set before initiating an STK Push.
// This internal method ensures that business code, transaction type, amount, phone number,
// and callback URL are all properly configured.
//
// Returns:
//   - error: An error describing which required parameter is missing, or nil if all are present
func (s *StkService) validatePushParams() error {
	if s.Config.GetBusinessCode() == "" {
		return errors.New("business code is required")
	}
	if s.transactionType == "" {
		return errors.New("transaction type is required")
	}
	if s.amount == "" {
		return errors.New("amount is required")
	}
	if s.phoneNumber == "" {
		return errors.New("phone number is required")
	}
	if s.callbackUrl == "" {
		return errors.New("callback URL is required")
	}
	return nil
}

// Push initiates an STK Push request to the customer's mobile phone.
// This method sends a payment request that will appear as a popup on the customer's phone,
// allowing them to authorize the payment using their M-Pesa PIN.
//
// Returns:
//   - *StkService: Returns self for method chaining and accessing response data
//   - error: An error if validation fails or the API request encounters issues
//
// Example:
//
//	stkService := mpesa.STK()
//	response, err := stkService.
//	    SetAmount(100).
//	    SetPhoneNumber("254711223344").
//	    SetTransactionType("CustomerPayBillOnline").
//	    SetCallbackUrl("https://example.com/callback").
//	    Push()
//	if err != nil {
//	    log.Printf("STK Push failed: %v", err)
//	    return
//	}
//
//	// Get checkout request ID for tracking
//	checkoutID, _ := response.GetCheckoutRequestID()
//	fmt.Printf("Payment initiated with ID: %s", checkoutID)
func (s *StkService) Push() (*StkService, error) {
	if err := s.validatePushParams(); err != nil {
		return s, err
	}

	data := map[string]any{
		"BusinessShortCode": s.Config.GetBusinessCode(),
		"Password":          s.GeneratePassword(),
		"Timestamp":         s.GenerateTimestamp(),
		"TransactionType":   s.transactionType,
		"Amount":            s.amount,
		"PartyA":            s.phoneNumber,
		"PartyB":            s.Config.GetBusinessCode(),
		"PhoneNumber":       s.phoneNumber,
		"CallBackURL":       s.callbackUrl,
		"AccountReference":  s.accountReference,
		"TransactionDesc":   s.transactionDesc,
	}

	// Provide defaults if empty
	if data["AccountReference"] == "" {
		data["AccountReference"] = "Account"
	}
	if data["TransactionDesc"] == "" {
		data["TransactionDesc"] = "Transaction"
	}

	resp, err := s.Client.ExecuteRequest(data, "/mpesa/stkpush/v1/processrequest")
	if err != nil {
		return s, err
	}

	s.response = resp
	return s, nil
}

// GetCheckoutRequestID extracts and returns the CheckoutRequestID from the STK Push response.
// This ID is used to track the payment status and query the transaction later.
// The method should be called after a successful Push() operation.
//
// Returns:
//   - string: The CheckoutRequestID for tracking the payment
//   - error: An error if no response is available or the ID is not found/invalid
//
// Example:
//
//	response, err := stkService.Push()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	checkoutID, err := response.GetCheckoutRequestID()
//	if err != nil {
//	    log.Printf("Failed to get checkout ID: %v", err)
//	    return
//	}
//
//	fmt.Printf("Track your payment with ID: %s", checkoutID)
func (s *StkService) GetCheckoutRequestID() (string, error) {
	if s.response == nil {
		return "", errors.New("no STK push response available")
	}

	val, ok := s.response["CheckoutRequestID"]
	if !ok {
		return "", errors.New("CheckoutRequestID not found in response")
	}

	id, ok := val.(string)
	if !ok {
		return "", errors.New("CheckoutRequestID is not a string")
	}

	return id, nil
}

// Query checks the status of an STK Push transaction using the CheckoutRequestID.
// This method can be called with a specific CheckoutRequestID or without parameters
// to query the status of the last Push() operation performed by this service instance.
//
// Parameters:
//   - checkoutRequestId: Optional CheckoutRequestID to query. If not provided,
//     uses the ID from the last Push() operation
//
// Returns:
//   - map[string]any: The query response containing transaction status and details
//   - error: An error if the query fails or no CheckoutRequestID is available
//
// Example:
//
//	// Query using specific checkout ID
//	status, err := stkService.Query("ws_CO_123456789")
//	if err != nil {
//	    log.Printf("Query failed: %v", err)
//	    return
//	}
//
//	// Query last transaction from this service instance
//	status, err := stkService.Query()
//	if err != nil {
//	    log.Printf("Query failed: %v", err)
//	    return
//	}
//
//	// Check transaction result
//	if resultCode, ok := status["ResultCode"].(string); ok && resultCode == "0" {
//	    fmt.Println("Payment was successful!")
//	} else {
//	    fmt.Printf("Payment failed or pending. Status: %+v", status)
//	}
func (s *StkService) Query(checkoutRequestId ...string) (map[string]any, error) {
	reqID := ""
	if len(checkoutRequestId) > 0 {
		reqID = checkoutRequestId[0]
	} else {
		var err error
		reqID, err = s.GetCheckoutRequestID()
		if err != nil {
			return nil, err
		}
	}

	data := map[string]any{
		"BusinessShortCode": s.Config.GetBusinessCode(),
		"Password":          s.GeneratePassword(),
		"Timestamp":         s.GenerateTimestamp(),
		"CheckoutRequestID": reqID,
	}

	return s.Client.ExecuteRequest(data, "/mpesa/stkpushquery/v1/query")
}

// GetResponse returns the raw response map from the last STK Push operation.
// This method provides access to the complete API response for advanced use cases
// where you need to access fields not covered by other methods.
//
// Returns:
//   - map[string]any: The complete response map from the last API call, or nil if no request has been made
//
// Example:
//
//	response, err := stkService.Push()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	rawResponse := response.GetResponse()
//	if merchantRequestID, ok := rawResponse["MerchantRequestID"].(string); ok {
//	    fmt.Printf("Merchant Request ID: %s", merchantRequestID)
//	}
//
//	if responseCode, ok := rawResponse["ResponseCode"].(string); ok {
//	    fmt.Printf("Response Code: %s", responseCode)
//	}
func (s *StkService) GetResponse() map[string]any {
	return s.response
}
