package Services

import (
	"errors"
	"fmt"
	abstracts "github.com/venomous-maker/go-mpesa/Abstracts"
)

// CustomerToBusinessService handles Customer to Business (C2B) payment operations.
// C2B allows customers to make payments to businesses and enables businesses to register
// validation and confirmation URLs for payment notifications.
type CustomerToBusinessService struct {
	Config          *abstracts.MpesaConfig   // M-Pesa configuration containing credentials and settings
	Client          abstracts.MpesaInterface // HTTP client interface for making API requests
	ConfirmationURL string                   // URL to receive payment confirmation notifications
	ValidationURL   string                   // URL to validate payment requests (optional)
	ResponseType    string                   // Response type for URL registration ("Completed" or "Cancelled")
	CommandID       string                   // Command ID for the transaction
	BillRefNumber   string                   // Bill reference number for the payment
	Amount          string                   // Amount for the payment simulation
	PhoneNumber     string                   // Customer's phone number for payment simulation
	Response        map[string]interface{}   // Response from the last API call
}

// NewCustomerToBusinessService creates a new C2B service instance with the provided configuration and client.
// This is the constructor for creating C2B service instances for handling customer payments.
//
// Parameters:
//   - cfg: M-Pesa configuration containing credentials and settings
//   - client: HTTP client interface for making API requests
//
// Returns:
//   - *CustomerToBusinessService: A configured C2B service ready for payment operations
//
// Example:
//
//	cfg := createMpesaConfig()
//	client := Abstracts.NewApiClient(cfg)
//	c2bService := NewCustomerToBusinessService(cfg, client)
func NewCustomerToBusinessService(cfg *abstracts.MpesaConfig, client abstracts.MpesaInterface) *CustomerToBusinessService {
	return &CustomerToBusinessService{
		Config: cfg,
		Client: client,
	}
}

// SetConfirmationURL sets the URL where M-Pesa will send payment confirmation notifications.
// This URL will receive POST requests when payments are successfully completed.
//
// Parameters:
//   - url: The fully qualified confirmation URL (must be HTTPS in production)
//
// Returns:
//   - *CustomerToBusinessService: Returns self for method chaining
//
// Example:
//
//	c2bService.SetConfirmationURL("https://yourdomain.com/mpesa/confirmation")
func (s *CustomerToBusinessService) SetConfirmationURL(url string) *CustomerToBusinessService {
	s.ConfirmationURL = url
	return s
}

// SetValidationURL sets the URL where M-Pesa will send payment validation requests.
// This URL allows you to validate payments before they are processed (optional).
//
// Parameters:
//   - url: The fully qualified validation URL (must be HTTPS in production)
//
// Returns:
//   - *CustomerToBusinessService: Returns self for method chaining
//
// Example:
//
//	c2bService.SetValidationURL("https://yourdomain.com/mpesa/validation")
func (s *CustomerToBusinessService) SetValidationURL(url string) *CustomerToBusinessService {
	s.ValidationURL = url
	return s
}

// SetResponseType sets the response type for URL registration.
// This determines how M-Pesa handles the validation response.
//
// Parameters:
//   - t: The response type ("Completed" or "Cancelled")
//
// Returns:
//   - *CustomerToBusinessService: Returns self for method chaining
//
// Response Types:
//   - "Completed": Default response - completes the transaction
//   - "Cancelled": Cancels the transaction if validation fails
//
// Example:
//
//	c2bService.SetResponseType("Completed")
//	c2bService.SetResponseType("Cancelled")
func (s *CustomerToBusinessService) SetResponseType(t string) *CustomerToBusinessService {
	s.ResponseType = t
	return s
}

// SetCommandID sets the command ID for C2B transactions.
// This identifies the type of transaction being performed.
//
// Parameters:
//   - cmd: The command ID (typically "CustomerPayBillOnline" or "CustomerBuyGoodsOnline")
//
// Returns:
//   - *CustomerToBusinessService: Returns self for method chaining
//
// Example:
//
//	c2bService.SetCommandID("CustomerPayBillOnline")
//	c2bService.SetCommandID("CustomerBuyGoodsOnline")
func (s *CustomerToBusinessService) SetCommandID(cmd string) *CustomerToBusinessService {
	s.CommandID = cmd
	return s
}

// SetBillRefNumber sets the bill reference number for the payment.
// This helps identify what the customer is paying for.
//
// Parameters:
//   - ref: The bill reference number or identifier
//
// Returns:
//   - *CustomerToBusinessService: Returns self for method chaining
//
// Example:
//
//	c2bService.SetBillRefNumber("INVOICE123")
//	c2bService.SetBillRefNumber("ACCOUNT456")
func (s *CustomerToBusinessService) SetBillRefNumber(ref string) *CustomerToBusinessService {
	s.BillRefNumber = ref
	return s
}

// SetAmount sets the amount for C2B payment simulation.
// The amount should be in Kenyan Shillings.
//
// Parameters:
//   - amount: The amount as a string
//
// Returns:
//   - *CustomerToBusinessService: Returns self for method chaining
//
// Example:
//
//	c2bService.SetAmount("100")
//	c2bService.SetAmount("1500")
func (s *CustomerToBusinessService) SetAmount(amount string) *CustomerToBusinessService {
	s.Amount = amount
	return s
}

// SetPhoneNumber sets the customer's phone number for payment simulation.
// The phone number should be in international format.
//
// Parameters:
//   - phone: The customer's phone number (e.g., "254711223344")
//
// Returns:
//   - *CustomerToBusinessService: Returns self for method chaining
//
// Example:
//
//	c2bService.SetPhoneNumber("254711223344")
func (s *CustomerToBusinessService) SetPhoneNumber(phone string) *CustomerToBusinessService {
	s.PhoneNumber = phone
	return s
}

// RegisterURLs registers the validation and confirmation URLs with M-Pesa.
// This must be done before customers can make C2B payments to your business.
//
// Returns:
//   - error: An error if URL registration fails
//
// Example:
//
//	err := c2bService.
//	    SetConfirmationURL("https://yourdomain.com/mpesa/confirmation").
//	    SetValidationURL("https://yourdomain.com/mpesa/validation").
//	    SetResponseType("Completed").
//	    RegisterURLs()
//	if err != nil {
//	    log.Printf("URL registration failed: %v", err)
//	}
func (s *CustomerToBusinessService) RegisterURLs() error {
	if s.ConfirmationURL == "" {
		return errors.New("confirmation URL is required")
	}

	data := map[string]interface{}{
		"ShortCode":       s.Config.GetBusinessCode(),
		"ResponseType":    s.getResponseType(),
		"ConfirmationURL": s.ConfirmationURL,
		"ValidationURL":   s.ValidationURL,
	}

	response, err := s.Client.ExecuteRequest(data, "/mpesa/c2b/v1/registerurl")
	if err != nil {
		return fmt.Errorf("URL registration failed: %w", err)
	}

	s.Response = response
	return nil
}

// Simulate simulates a C2B payment for testing purposes.
// This is useful for testing your C2B integration in sandbox environment.
//
// Returns:
//   - map[string]interface{}: The simulation response from M-Pesa
//   - error: An error if the simulation fails
//
// Example:
//
//	response, err := c2bService.
//	    SetCommandID("CustomerPayBillOnline").
//	    SetAmount("100").
//	    SetPhoneNumber("254711223344").
//	    SetBillRefNumber("INVOICE123").
//	    Simulate()
//	if err != nil {
//	    log.Printf("C2B simulation failed: %v", err)
//	    return
//	}
//	fmt.Printf("Simulation response: %+v", response)
func (s *CustomerToBusinessService) Simulate() (map[string]interface{}, error) {
	if s.CommandID == "" {
		return nil, errors.New("command ID is required")
	}
	if s.Amount == "" {
		return nil, errors.New("amount is required")
	}
	if s.PhoneNumber == "" {
		return nil, errors.New("phone number is required")
	}

	data := map[string]interface{}{
		"ShortCode":     s.Config.GetBusinessCode(),
		"CommandID":     s.CommandID,
		"Amount":        s.Amount,
		"Msisdn":        s.PhoneNumber,
		"BillRefNumber": s.getBillRefNumber(),
	}

	response, err := s.Client.ExecuteRequest(data, "/mpesa/c2b/v1/simulate")
	if err != nil {
		return nil, fmt.Errorf("C2B simulation failed: %w", err)
	}

	s.Response = response
	return response, nil
}

// GetResponse returns the response from the last API call.
//
// Returns:
//   - map[string]interface{}: The response data, or nil if no call has been made
func (s *CustomerToBusinessService) GetResponse() map[string]interface{} {
	return s.Response
}

// getResponseType returns the response type, defaulting to "Completed" if not set.
func (s *CustomerToBusinessService) getResponseType() string {
	if s.ResponseType == "" {
		return "Completed"
	}
	return s.ResponseType
}

// getBillRefNumber returns the bill reference number, defaulting to "default" if not set.
func (s *CustomerToBusinessService) getBillRefNumber() string {
	if s.BillRefNumber == "" {
		return "default"
	}
	return s.BillRefNumber
}
