package Services

import (
	"errors"
	abstracts "github.com/venomous-maker/go-mpesa/Abstracts"
)

// BusinessToCustomerService handles business-to-customer payment requests.
type BusinessToCustomerService struct {
	Config        *abstracts.MpesaConfig
	Client        abstracts.MpesaInterface
	initiatorName string
	commandID     string
	remarks       string
	occasion      string
	amount        int
	phoneNumber   string
}

// NewBusinessToCustomerService creates a new BusinessToCustomerService instance.
func NewBusinessToCustomerService(cfg *abstracts.MpesaConfig, client abstracts.MpesaInterface) *BusinessToCustomerService {
	return &BusinessToCustomerService{
		Config: cfg,
		Client: client,
	}
}

// SetInitiatorName sets the username of the M-Pesa API operator.
func (s *BusinessToCustomerService) SetInitiatorName(name string) *BusinessToCustomerService {
	s.initiatorName = name
	return s
}

// SetCommandID sets the command ID (e.g. SalaryPayment, BusinessPayment, PromotionPayment).
func (s *BusinessToCustomerService) SetCommandID(cmd string) *BusinessToCustomerService {
	s.commandID = cmd
	return s
}

// SetRemarks sets remarks for the transaction.
func (s *BusinessToCustomerService) SetRemarks(remarks string) *BusinessToCustomerService {
	s.remarks = remarks
	return s
}

// SetOccasion sets the occasion for the transaction.
func (s *BusinessToCustomerService) SetOccasion(occasion string) *BusinessToCustomerService {
	s.occasion = occasion
	return s
}

// SetAmount sets the amount for the transaction.
func (s *BusinessToCustomerService) SetAmount(amount int) *BusinessToCustomerService {
	s.amount = amount
	return s
}

// SetPhoneNumber sets the customer's phone number.
func (s *BusinessToCustomerService) SetPhoneNumber(phone string) *BusinessToCustomerService {
	s.phoneNumber = phone
	return s
}

// PaymentRequest sends a business to customer payment request to the M-Pesa API.
// All parameters are optional. If provided, they override the existing fields.
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
