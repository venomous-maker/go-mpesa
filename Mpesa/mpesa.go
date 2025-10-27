// Package Mpesa provides a comprehensive Go SDK for integrating with Safaricom's M-Pesa API.
// It offers a clean, type-safe interface for M-Pesa services including STK Push, B2C, C2B,
// Account Balance, Transaction Status, and Reversals.
package Mpesa

import (
	"github.com/venomous-maker/go-mpesa/Abstracts"
	"github.com/venomous-maker/go-mpesa/Services"
)

// Mpesa represents the main client for interacting with M-Pesa API services.
// It encapsulates the configuration and API client needed for making requests.
type Mpesa struct {
	Config *Abstracts.MpesaConfig // Configuration for M-Pesa API credentials and settings
	Client *Abstracts.ApiClient   // HTTP client for making API requests
}

// New creates a new Mpesa instance with the provided credentials and environment.
// This is the primary constructor for initializing the M-Pesa SDK.
//
// Parameters:
//   - consumerKey: The consumer key obtained from Safaricom Developer Portal
//   - consumerSecret: The consumer secret obtained from Safaricom Developer Portal
//   - environment: The target environment ("sandbox" or "production")
//
// Returns:
//   - *Mpesa: A configured Mpesa instance ready for API calls
//   - error: An error if configuration fails
//
// Example:
//
//	mpesa, err := Mpesa.New("your_consumer_key", "your_consumer_secret", "sandbox")
//	if err != nil {
//	    log.Fatal(err)
//	}
func New(consumerKey, consumerSecret, environment string) (*Mpesa, error) {
	cfg, err := Abstracts.NewMpesaConfig(
		consumerKey,
		consumerSecret,
		Abstracts.Environment(environment),
		nil, nil, nil, nil, nil,
	)
	if err != nil {
		return nil, err
	}

	client := Abstracts.NewApiClient(cfg)

	return &Mpesa{
		Config: cfg,
		Client: client,
	}, nil
}

// SetCredentials updates the API credentials and environment for the Mpesa instance.
// This method allows changing credentials without creating a new instance.
//
// Parameters:
//   - consumerKey: The new consumer key
//   - consumerSecret: The new consumer secret
//   - environment: The new target environment ("sandbox" or "production")
//
// Returns:
//   - error: An error if credential update fails
//
// Example:
//
//	err := mpesa.SetCredentials("new_key", "new_secret", "production")
//	if err != nil {
//	    log.Printf("Failed to update credentials: %v", err)
//	}
func (m *Mpesa) SetCredentials(consumerKey, consumerSecret, environment string) error {
	cfg, err := Abstracts.NewMpesaConfig(
		consumerKey,
		consumerSecret,
		Abstracts.Environment(environment),
		nil, nil, nil, nil, nil,
	)
	if err != nil {
		return err
	}
	m.Config = cfg
	m.Client = Abstracts.NewApiClient(cfg)
	return nil
}

// SetBusinessCode sets the business shortcode for M-Pesa transactions.
// The business shortcode is required for most M-Pesa operations and identifies
// your business in the M-Pesa system.
//
// Parameters:
//   - code: The business shortcode (e.g., "174379" for sandbox)
//
// Example:
//
//	mpesa.SetBusinessCode("174379") // Sandbox shortcode
//	mpesa.SetBusinessCode("123456") // Production shortcode
func (m *Mpesa) SetBusinessCode(code string) {
	m.Config.SetBusinessCode(code)
}

// SetPassKey sets the Lipa na M-Pesa Online passkey used for STK Push transactions.
// The passkey is obtained from Safaricom and is required for STK Push operations.
//
// Parameters:
//   - passkey: The Lipa na M-Pesa Online passkey
//
// Example:
//
//	mpesa.SetPassKey("bfb279f9aa9bdbcf158e97dd71a467cd2e0c893059b10f78e6b72ada1ed2c919")
func (m *Mpesa) SetPassKey(passkey string) {
	m.Config.SetPassKey(passkey)
}

// STK creates and returns a new STK Push service instance.
// STK Push allows initiating M-Pesa payments directly from a customer's phone.
//
// Returns:
//   - *Services.StkService: A configured STK service for payment operations
//
// Example:
//
//	stkService := mpesa.STK()
//	response, err := stkService.
//	    SetAmount(100).
//	    SetPhoneNumber("254711223344").
//	    SetCallbackUrl("https://example.com/callback").
//	    Push()
func (m *Mpesa) STK() *Services.StkService {
	return Services.NewStkService(m.Config, m.Client)
}

// B2C creates and returns a new Business to Customer service instance.
// B2C allows sending money from your business account to customer accounts.
//
// Returns:
//   - *Services.B2cService: A configured B2C service for money transfers
//
// Example:
//
//	b2cService := mpesa.B2C()
//	response, err := b2cService.
//	    SetAmount(1000).
//	    SetPhoneNumber("254711223344").
//	    SetCommandID("BusinessPayment").
//	    Send()
func (m *Mpesa) B2C() *Services.BusinessToCustomerService {
	return Services.NewBusinessToCustomerService(m.Config, m.Client)
}

// C2B creates and returns a new Customer to Business service instance.
// C2B allows registering URLs and simulating customer payments to your business.
//
// Returns:
//   - *Services.C2bService: A configured C2B service for receiving payments
//
// Example:
//
//	c2bService := mpesa.C2B()
//	err := c2bService.
//	    SetValidationURL("https://example.com/validation").
//	    SetConfirmationURL("https://example.com/confirmation").
//	    RegisterURLs()
func (m *Mpesa) C2B() *Services.CustomerToBusinessService {
	return Services.NewCustomerToBusinessService(m.Config, m.Client)
}

// AccountBalance creates and returns a new Account Balance service instance.
// This service allows querying the balance of your M-Pesa business account.
//
// Returns:
//   - *Services.AccountBalanceService: A configured service for balance inquiries
//
// Example:
//
//	balanceService := mpesa.AccountBalance()
//	balance, err := balanceService.
//	    SetCommandID("AccountBalance").
//	    SetRemarks("Balance inquiry").
//	    Query()
func (m *Mpesa) AccountBalance() *Services.AccountBalanceService {
	return Services.NewAccountBalanceService(m.Config, m.Client)
}

// TransactionStatus creates and returns a new Transaction Status service instance.
// This service allows querying the status of any M-Pesa transaction.
//
// Returns:
//   - *Services.TransactionStatusService: A configured service for status queries
//
// Example:
//
//	statusService := mpesa.TransactionStatus()
//	status, err := statusService.
//	    SetTransactionID("ABC123XYZ").
//	    SetCommandID("TransactionStatusQuery").
//	    Query()
func (m *Mpesa) TransactionStatus() *Services.TransactionStatusService {
	return Services.NewTransactionStatusService(m.Config, m.Client)
}

// Reversal creates and returns a new Reversal service instance.
// This service allows reversing completed M-Pesa transactions.
//
// Returns:
//   - *Services.ReversalService: A configured service for transaction reversals
//
// Example:
//
//	reversalService := mpesa.Reversal()
//	response, err := reversalService.
//	    SetTransactionID("ABC123XYZ").
//	    SetAmount(1000).
//	    SetCommandID("TransactionReversal").
//	    Reverse()
func (m *Mpesa) Reversal() *Services.ReversalService {
	return Services.NewReversalService(m.Config, m.Client)
}

// B2PayBill creates and returns a new Business-to-PayBill service instance.
// This service allows a business to pay directly to a PayBill number or store on behalf of a consumer.
//
// Returns:
//   - *Services.BusinessToPayBillService: A configured service for B2B PayBill payments
//
// Example:
//
//	b2paybillService := mpesa.B2PayBill()
//	response, err := b2paybillService.
//	    SetInitiator("testapi").
//	    SetSecurityCredential("your_security_credential").
//	    SetAmount(1000).
//	    SetPartyA("174379").
//	    SetPartyB("123456").
//	    SetAccountReference("ABC123").
//	    SetRequester("254711223344").
//	    SetRemarks("Payment for goods").
//	    SetOccasion("Payment").
//	    SetQueueTimeoutURL("https://example.com/timeout").
//	    SetResultURL("https://example.com/result").
//	    Send()
func (m *Mpesa) B2PayBill() *Services.BusinessToPayBillService {
	return Services.NewBusinessToPayBillService(m.Config, m.Client)
}
