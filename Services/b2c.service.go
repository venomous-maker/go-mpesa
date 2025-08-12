package Services

import (
	"errors"
	abstracts "github.com/venomous-maker/go-mpesa/Abstracts"
)

// BusinessToCustomerService handles Business to Customer (B2C) payment operations.
// B2C allows businesses to send money directly to customer M-Pesa accounts.
// This service supports various payment types including salary payments, business payments, and promotional payments.
type BusinessToCustomerService struct {
	Config        *abstracts.MpesaConfig   // M-Pesa configuration containing credentials and settings
	Client        abstracts.MpesaInterface // HTTP client interface for making API requests
	initiatorName string                   // Username of the M-Pesa API operator
	commandID     string                   // Type of B2C payment (SalaryPayment, BusinessPayment, etc.)
	remarks       string                   // Transaction remarks/description
	occasion      string                   // Occasion for the payment
	amount        int                      // Amount to be sent to the customer
	phoneNumber   string                   // Customer's phone number
}

// NewBusinessToCustomerService creates a new B2C service instance with the provided configuration and client.
// This is the constructor for creating B2C service instances that can be used to send money to customers.
//
// Parameters:
//   - cfg: M-Pesa configuration containing credentials and settings
//   - client: HTTP client interface for making API requests
//
// Returns:
//   - *BusinessToCustomerService: A configured B2C service ready for payment operations
//
// Example:
//
//	cfg := createMpesaConfig()
//	client := Abstracts.NewApiClient(cfg)
//	b2cService := NewBusinessToCustomerService(cfg, client)
func NewBusinessToCustomerService(cfg *abstracts.MpesaConfig, client abstracts.MpesaInterface) *BusinessToCustomerService {
	return &BusinessToCustomerService{
		Config: cfg,
		Client: client,
	}
}

// SetInitiatorName sets the username of the M-Pesa API operator initiating the transaction.
// This is typically the username provided by Safaricom for B2C operations.
//
// Parameters:
//   - name: The initiator username/name
//
// Returns:
//   - *BusinessToCustomerService: Returns self for method chaining
//
// Example:
//
//	b2cService.SetInitiatorName("testapi")
func (s *BusinessToCustomerService) SetInitiatorName(name string) *BusinessToCustomerService {
	s.initiatorName = name
	return s
}

// SetCommandID sets the type of B2C payment being made.
// Different command IDs are used for different types of payments.
//
// Parameters:
//   - cmd: The command ID for the payment type
//
// Returns:
//   - *BusinessToCustomerService: Returns self for method chaining
//
// Common Command IDs:
//   - "SalaryPayment": For salary disbursements
//   - "BusinessPayment": For general business payments
//   - "PromotionPayment": For promotional payments and rewards
//
// Example:
//
//	b2cService.SetCommandID("SalaryPayment")
//	b2cService.SetCommandID("BusinessPayment")
func (s *BusinessToCustomerService) SetCommandID(cmd string) *BusinessToCustomerService {
	s.commandID = cmd
	return s
}

// SetRemarks sets the remarks or description for the B2C transaction.
// This helps identify the purpose of the payment in transaction records.
//
// Parameters:
//   - remarks: A descriptive string for the transaction
//
// Returns:
//   - *BusinessToCustomerService: Returns self for method chaining
//
// Example:
//
//	b2cService.SetRemarks("Monthly salary payment")
//	b2cService.SetRemarks("Bonus payment for Q4 performance")
func (s *BusinessToCustomerService) SetRemarks(remarks string) *BusinessToCustomerService {
	s.remarks = remarks
	return s
}

// SetOccasion sets the occasion or reason for the B2C payment.
// This provides additional context for the transaction.
//
// Parameters:
//   - occasion: A string describing the occasion for the payment
//
// Returns:
//   - *BusinessToCustomerService: Returns self for method chaining
//
// Example:
//
//	b2cService.SetOccasion("December 2024 Salary")
//	b2cService.SetOccasion("Annual bonus distribution")
func (s *BusinessToCustomerService) SetOccasion(occasion string) *BusinessToCustomerService {
	s.occasion = occasion
	return s
}

// SetAmount sets the amount to be sent to the customer.
// The amount should be in Kenyan Shillings.
//
// Parameters:
//   - amount: The amount in KES to send to the customer
//
// Returns:
//   - *BusinessToCustomerService: Returns self for method chaining
//
// Example:
//
//	b2cService.SetAmount(50000)  // Send KES 50,000
//	b2cService.SetAmount(1000)   // Send KES 1,000
func (s *BusinessToCustomerService) SetAmount(amount int) *BusinessToCustomerService {
	s.amount = amount
	return s
}

// SetPhoneNumber sets the customer's phone number for the B2C payment.
// The phone number should be in international format.
//
// Parameters:
//   - phone: The customer's phone number (e.g., "254711223344")
//
// Returns:
//   - *BusinessToCustomerService: Returns self for method chaining
//
// Example:
//
//	b2cService.SetPhoneNumber("254711223344")
//	b2cService.SetPhoneNumber("254722000000")
func (s *BusinessToCustomerService) SetPhoneNumber(phone string) *BusinessToCustomerService {
	s.phoneNumber = phone
	return s
}

// PaymentRequest sends a business to customer payment request to the M-Pesa API.
// All parameters are optional. If provided, they override the existing fields.
//
// Parameters:
//   - initiatorName: Optional initiator username
//   - initiatorPassword: Optional initiator password for security credential
//   - commandID: Optional command ID (e.g. SalaryPayment, BusinessPayment)
//   - amount: Optional amount for the transaction
//   - partyA: Optional business short code
//   - phoneNumber: Optional customer's phone number
//   - remarks: Optional transaction remarks
//   - queueTimeoutURL: Optional URL for queue timeout notification
//   - resultURL: Optional URL for result notification
//   - occasion: Optional occasion for the transaction
//
// Returns:
//   - map[string]interface{}: The response from the M-Pesa API
//   - error: An error if validation fails or the API request encounters issues
//
// Example:
//
//	response, err := b2cService.PaymentRequest(
//	    abstracts.String("testapi"),
//	    abstracts.String("password"),
//	    abstracts.String("SalaryPayment"),
//	    abstracts.Int(5000),
//	    abstracts.String("123456"),
//	    abstracts.String("254711223344"),
//	    abstracts.String("Monthly salary"),
//	    abstracts.String("https://example.com/timeout"),
//	    abstracts.String("https://example.com/result"),
//	    abstracts.String("December Salary"),
//	)
//	if err != nil {
//	    log.Printf("Payment request failed: %v", err)
//	    return
//	}
//	fmt.Printf("Payment request sent: %+v", response)
func (s *BusinessToCustomerService) PaymentRequest(
	initiatorName, initiatorPassword, commandID *string,
	amount *int,
	partyA, phoneNumber, remarks, queueTimeoutURL, resultURL, occasion *string,
) (map[string]interface{}, error) {
	if initiatorName != nil {
		s.SetInitiatorName(*initiatorName)
	}
	if commandID != nil {
		s.SetCommandID(*commandID)
	}
	if amount != nil {
		s.SetAmount(*amount)
	}
	if partyA != nil {
		s.Config.SetBusinessCode(*partyA)
	}
	if phoneNumber != nil {
		s.SetPhoneNumber(*phoneNumber)
	}
	if remarks != nil {
		s.SetRemarks(*remarks)
	}
	if queueTimeoutURL != nil {
		s.Config.SetQueueTimeoutURL(*queueTimeoutURL)
	}
	if resultURL != nil {
		s.Config.SetResultURL(*resultURL)
	}
	if occasion != nil {
		s.SetOccasion(*occasion)
	}
	if initiatorPassword != nil {
		if err := s.Config.SetSecurityCredential(*initiatorPassword); err != nil {
			return nil, err
		}
	}

	// Validate required fields
	if s.initiatorName == "" || s.commandID == "" || s.amount == 0 || s.phoneNumber == "" || s.Config.GetBusinessCode() == "" {
		return nil, errors.New("initiator name, command ID, amount, phone number, and business code are required")
	}

	requestData := map[string]interface{}{
		"InitiatorName":      s.initiatorName,
		"SecurityCredential": s.Config.GetSecurityCredential(),
		"CommandID":          s.commandID,
		"Amount":             s.amount,
		"PartyA":             s.Config.GetBusinessCode(),
		"PartyB":             s.phoneNumber,
		"Remarks":            s.remarks,
		"QueueTimeOutURL":    s.Config.GetQueueTimeoutURL(),
		"ResultURL":          s.Config.GetResultURL(),
		"Occassion":          s.occasion,
	}

	response, err := s.Client.ExecuteRequest(requestData, "/mpesa/b2c/v1/paymentrequest")
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Send initiates the B2C payment to the customer.
// This method validates all required parameters and sends the payment request to M-Pesa.
//
// Returns:
//   - map[string]any: The response from the M-Pesa API
//   - error: An error if validation fails or the API request encounters issues
//
// Example:
//
//	response, err := b2cService.
//	    SetInitiatorName("testapi").
//	    SetCommandID("SalaryPayment").
//	    SetAmount(50000).
//	    SetPhoneNumber("254711223344").
//	    SetRemarks("Monthly salary").
//	    SetOccasion("December 2024").
//	    Send()
//	if err != nil {
//	    log.Printf("B2C payment failed: %v", err)
//	    return
//	}
//	fmt.Printf("Payment initiated: %+v", response)
func (s *BusinessToCustomerService) Send() (map[string]any, error) {
	// Validate required fields
	if s.initiatorName == "" {
		return nil, errors.New("initiator name is required")
	}
	if s.commandID == "" {
		return nil, errors.New("command ID is required")
	}
	if s.amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	if s.phoneNumber == "" {
		return nil, errors.New("phone number is required")
	}

	data := map[string]any{
		"InitiatorName":      s.initiatorName,
		"SecurityCredential": s.Config.GetSecurityCredential(),
		"CommandID":          s.commandID,
		"Amount":             s.amount,
		"PartyA":             s.Config.GetBusinessCode(),
		"PartyB":             s.phoneNumber,
		"Remarks":            s.remarks,
		"QueueTimeOutURL":    s.Config.GetQueueTimeoutURL(),
		"ResultURL":          s.Config.GetResultURL(),
		"Occasion":           s.occasion,
	}

	return s.Client.ExecuteRequest(data, "/mpesa/b2c/v1/paymentrequest")
}
