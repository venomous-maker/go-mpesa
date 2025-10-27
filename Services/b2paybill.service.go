package Services

import (
	"errors"

	abstracts "github.com/venomous-maker/go-mpesa/Abstracts"
)

// BusinessToPayBillService handles Business-to-PayBill (B2B PayBill) payments.
// This service allows a business to pay directly to a PayBill number or store on behalf of a consumer.
type BusinessToPayBillService struct {
	Config                  *abstracts.MpesaConfig
	Client                  abstracts.MpesaInterface
	initiator               string
	commandID               string
	senderIdentifierType    string
	recipientIdentifierType string
	amount                  float64
	partyA                  string
	partyB                  string
	accountReference        string
	requester               string
	remarks                 string
	occasion                string
	response                map[string]any
}

// NewBusinessToPayBillService creates a new B2B PayBill service instance.
func NewBusinessToPayBillService(cfg *abstracts.MpesaConfig, client abstracts.MpesaInterface) *BusinessToPayBillService {
	return &BusinessToPayBillService{
		Config:                  cfg,
		Client:                  client,
		commandID:               "BusinessPayBill",
		senderIdentifierType:    "4",
		recipientIdentifierType: "4",
	}
}

// SetInitiator sets the initiator (operator username) for the transaction.
func (s *BusinessToPayBillService) SetInitiator(name string) *BusinessToPayBillService {
	s.initiator = name
	return s
}

// SetSecurityCredential encrypts and sets the security credential (initiator password).
func (s *BusinessToPayBillService) SetSecurityCredential(password string) error {
	return s.Config.SetSecurityCredential(password)
}

// SetAmount sets the transaction amount in KES.
func (s *BusinessToPayBillService) SetAmount(amount float64) *BusinessToPayBillService {
	s.amount = amount
	return s
}

// SetPartyA sets the shortcode from which money will be deducted.
func (s *BusinessToPayBillService) SetPartyA(code string) *BusinessToPayBillService {
	s.partyA = code
	s.Config.SetBusinessCode(code)
	return s
}

// SetPartyB sets the destination shortcode (paybill) to which money will be credited.
func (s *BusinessToPayBillService) SetPartyB(code string) *BusinessToPayBillService {
	s.partyB = code
	return s
}

// SetAccountReference sets an account/reference associated with the payment.
func (s *BusinessToPayBillService) SetAccountReference(ref string) *BusinessToPayBillService {
	s.accountReference = ref
	return s
}

// SetRequester sets the optional consumer mobile number on whose behalf the payment is made.
func (s *BusinessToPayBillService) SetRequester(msisdn string) *BusinessToPayBillService {
	s.requester = msisdn
	return s
}

// SetRemarks sets transaction remarks.
func (s *BusinessToPayBillService) SetRemarks(r string) *BusinessToPayBillService {
	s.remarks = r
	return s
}

// SetOccasion sets the occasion for the transaction.
func (s *BusinessToPayBillService) SetOccasion(o string) *BusinessToPayBillService {
	s.occasion = o
	return s
}

// SetQueueTimeoutURL updates the config queue timeout URL.
func (s *BusinessToPayBillService) SetQueueTimeoutURL(url string) *BusinessToPayBillService {
	s.Config.SetQueueTimeoutURL(url)
	return s
}

// SetResultURL updates the config result URL.
func (s *BusinessToPayBillService) SetResultURL(url string) *BusinessToPayBillService {
	s.Config.SetResultURL(url)
	return s
}

// Send constructs and sends the B2B BusinessPayBill payment request to M-Pesa using shared helper.
func (s *BusinessToPayBillService) Send() (map[string]any, error) {
	// Validate required fields
	if s.initiator == "" {
		return nil, errors.New("initiator is required")
	}
	if s.Config.GetSecurityCredential() == "" {
		return nil, errors.New("security credential is required; call SetSecurityCredential")
	}
	if s.amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	if s.partyA == "" && s.Config.GetBusinessCode() == "" {
		return nil, errors.New("partyA (business shortcode) is required")
	}
	if s.partyB == "" {
		return nil, errors.New("partyB (destination shortcode/paybill) is required")
	}

	req := B2BRequest{
		Initiator:              s.initiator,
		SecurityCredential:     s.Config.GetSecurityCredential(),
		CommandID:              s.commandID,
		SenderIdentifierType:   s.senderIdentifierType,
		RecieverIdentifierType: s.recipientIdentifierType,
		Amount:                 s.amount,
		PartyA:                 s.getPartyA(),
		PartyB:                 s.partyB,
		AccountReference:       s.accountReference,
		Requester:              s.requester,
		Remarks:                s.remarks,
		QueueTimeOutURL:        s.Config.GetQueueTimeoutURL(),
		ResultURL:              s.Config.GetResultURL(),
		Occasion:               s.occasion,
	}

	resp, err := ExecuteB2BRequest(s.Config, s.Client, req)
	if err != nil {
		return nil, err
	}

	s.response = resp
	return resp, nil
}

// ParseCallback parses a received callback payload using the shared ParseB2BCallback helper.
func (s *BusinessToPayBillService) ParseCallback(payload map[string]any) (*B2BCallbackResult, error) {
	return ParseB2BCallback(payload)
}

func (s *BusinessToPayBillService) getPartyA() string {
	if s.partyA != "" {
		return s.partyA
	}
	return s.Config.GetBusinessCode()
}

// GetResponse returns the last API response stored by the service.
func (s *BusinessToPayBillService) GetResponse() map[string]any {
	return s.response
}
