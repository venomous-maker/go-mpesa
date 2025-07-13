package Services

import (
	"errors"
	"fmt"
	abstracts "github.com/venomous-maker/go-mpesa/Abstracts"
)

// CustomerToBusinessService handles C2B operations.
type CustomerToBusinessService struct {
	Config          *abstracts.MpesaConfig
	Client          abstracts.MpesaInterface
	ConfirmationURL string
	ValidationURL   string
	ResponseType    string
	CommandID       string
	BillRefNumber   string
	Amount          string
	PhoneNumber     string
	Response        map[string]interface{}
}

// NewCustomerToBusinessService initializes a new CustomerToBusinessService.
func NewCustomerToBusinessService(cfg *abstracts.MpesaConfig, client abstracts.MpesaInterface) *CustomerToBusinessService {
	return &CustomerToBusinessService{
		Config: cfg,
		Client: client,
	}
}

// SetConfirmationURL sets the confirmation URL.
func (s *CustomerToBusinessService) SetConfirmationURL(url string) *CustomerToBusinessService {
	s.ConfirmationURL = url
	return s
}

// SetValidationURL sets the validation URL.
func (s *CustomerToBusinessService) SetValidationURL(url string) *CustomerToBusinessService {
	s.ValidationURL = url
	return s
}

// SetResponseType sets the response type.
func (s *CustomerToBusinessService) SetResponseType(t string) *CustomerToBusinessService {
	s.ResponseType = t
	return s
}

// SetCommandID sets the command ID.
func (s *CustomerToBusinessService) SetCommandID(cmd string) *CustomerToBusinessService {
	s.CommandID = cmd
	return s
}

// SetBillRefNumber sets the bill reference number.
func (s *CustomerToBusinessService) SetBillRefNumber(ref string) *CustomerToBusinessService {
	s.BillRefNumber = ref
	return s
}

// SetPhoneNumber sets the customer's phone number.
func (s *CustomerToBusinessService) SetPhoneNumber(phone string) *CustomerToBusinessService {
	s.PhoneNumber = phone
	return s
}

// SetAmount sets the transaction amount.
func (s *CustomerToBusinessService) SetAmount(amount interface{}) *CustomerToBusinessService {
	switch v := amount.(type) {
	case int:
		s.Amount = fmt.Sprintf("%d", v)
	case string:
		s.Amount = v
	}
	return s
}

// RegisterURL registers C2B confirmation and validation URLs.
func (s *CustomerToBusinessService) RegisterURL() (*CustomerToBusinessService, error) {
	if s.Config.GetBusinessCode() == "" {
		return nil, errors.New("business code is required")
	}
	if s.ConfirmationURL == "" {
		return nil, errors.New("confirmation URL is required")
	}
	if s.ValidationURL == "" {
		return nil, errors.New("validation URL is required")
	}

	data := map[string]interface{}{
		"ShortCode":       s.Config.GetBusinessCode(),
		"ResponseType":    s.ResponseType,
		"ConfirmationURL": s.ConfirmationURL,
		"ValidationURL":   s.ValidationURL,
	}

	resp, err := s.Client.ExecuteRequest(data, "/mpesa/c2b/v1/registerurl")
	if err != nil {
		return nil, err
	}
	s.Response = resp
	return s, nil
}

// Simulate simulates a C2B payment transaction.
func (s *CustomerToBusinessService) Simulate(phoneNumber *string, amount *string) (*CustomerToBusinessService, error) {
	if phoneNumber != nil {
		s.SetPhoneNumber(*phoneNumber)
	}
	if amount != nil {
		s.SetAmount(*amount)
	}

	if s.Config.GetBusinessCode() == "" {
		return nil, errors.New("business code is required")
	}
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
		"BillRefNumber": s.BillRefNumber,
	}

	resp, err := s.Client.ExecuteRequest(data, "/mpesa/c2b/v1/simulate")
	if err != nil {
		return nil, err
	}

	s.Response = resp
	return s, nil
}

// GetResponse returns the last C2B API response.
func (s *CustomerToBusinessService) GetResponse() map[string]interface{} {
	return s.Response
}
