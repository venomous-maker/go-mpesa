package Services

import (
	"errors"
	"fmt"
	"github.com/venomous-maker/go-mpesa/Abstracts"
	"strconv"
)

type StkService struct {
	*BaseService

	transactionType  string
	amount           string
	phoneNumber      string
	callbackUrl      string
	accountReference string
	transactionDesc  string

	response map[string]any
}

// NewStkService creates a new STK Push service
func NewStkService(cfg *Abstracts.MpesaConfig, client Abstracts.MpesaInterface) *StkService {
	return &StkService{
		BaseService: NewBaseService(cfg, client),
	}
}

func (s *StkService) SetTransactionType(t string) *StkService {
	s.transactionType = t
	return s
}

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

func (s *StkService) SetPhoneNumber(phone string) (*StkService, error) {
	cleaned, err := s.CleanPhoneNumber(phone, "254")
	if err != nil {
		return s, err
	}
	s.phoneNumber = cleaned
	return s, nil
}

func (s *StkService) SetCallbackUrl(url string) *StkService {
	s.callbackUrl = url
	return s
}

func (s *StkService) SetAccountReference(ref string) *StkService {
	s.accountReference = ref
	return s
}

func (s *StkService) SetTransactionDesc(desc string) *StkService {
	s.transactionDesc = desc
	return s
}

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

// Push initiates the STK Push request
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

// GetCheckoutRequestID returns the CheckoutRequestID from the response
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

// Query checks the status of an STK push transaction by CheckoutRequestID
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

// GetResponse returns the raw response map
func (s *StkService) GetResponse() map[string]any {
	return s.response
}
