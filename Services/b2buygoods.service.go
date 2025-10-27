package Services

import (
	"errors"

	abstracts "github.com/venomous-maker/go-mpesa/Abstracts"
)

// BusinessBuyGoodsService handles Business-Buy-Goods (B2B BusinessBuyGoods) payments.
// This service allows a business to pay merchants (till, store, HO) on behalf of a consumer.
type BusinessBuyGoodsService struct {
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

// NewBusinessBuyGoodsService creates a new BusinessBuyGoodsService instance.
func NewBusinessBuyGoodsService(cfg *abstracts.MpesaConfig, client abstracts.MpesaInterface) *BusinessBuyGoodsService {
	return &BusinessBuyGoodsService{
		Config:                  cfg,
		Client:                  client,
		commandID:               "BusinessBuyGoods",
		senderIdentifierType:    "4",
		recipientIdentifierType: "4",
	}
}

// SetInitiator sets the initiator (operator username) for the transaction.
func (s *BusinessBuyGoodsService) SetInitiator(name string) *BusinessBuyGoodsService {
	s.initiator = name
	return s
}

// SetSecurityCredential encrypts and sets the security credential (initiator password).
func (s *BusinessBuyGoodsService) SetSecurityCredential(password string) error {
	return s.Config.SetSecurityCredential(password)
}

// SetAmount sets the transaction amount in KES.
func (s *BusinessBuyGoodsService) SetAmount(amount float64) *BusinessBuyGoodsService {
	s.amount = amount
	return s
}

// SetPartyA sets the shortcode from which money will be deducted.
func (s *BusinessBuyGoodsService) SetPartyA(code string) *BusinessBuyGoodsService {
	s.partyA = code
	s.Config.SetBusinessCode(code)
	return s
}

// SetPartyB sets the destination shortcode (merchant) to which money will be credited.
func (s *BusinessBuyGoodsService) SetPartyB(code string) *BusinessBuyGoodsService {
	s.partyB = code
	return s
}

// SetAccountReference sets an account/reference associated with the payment.
func (s *BusinessBuyGoodsService) SetAccountReference(ref string) *BusinessBuyGoodsService {
	s.accountReference = ref
	return s
}

// SetRequester sets the optional consumer mobile number on whose behalf the payment is made.
func (s *BusinessBuyGoodsService) SetRequester(msisdn string) *BusinessBuyGoodsService {
	s.requester = msisdn
	return s
}

// SetRemarks sets transaction remarks.
func (s *BusinessBuyGoodsService) SetRemarks(r string) *BusinessBuyGoodsService {
	s.remarks = r
	return s
}

// SetOccasion sets the occasion for the transaction.
func (s *BusinessBuyGoodsService) SetOccasion(o string) *BusinessBuyGoodsService {
	s.occasion = o
	return s
}

// SetQueueTimeoutURL updates the config queue timeout URL.
func (s *BusinessBuyGoodsService) SetQueueTimeoutURL(url string) *BusinessBuyGoodsService {
	s.Config.SetQueueTimeoutURL(url)
	return s
}

// SetResultURL updates the config result URL.
func (s *BusinessBuyGoodsService) SetResultURL(url string) *BusinessBuyGoodsService {
	s.Config.SetResultURL(url)
	return s
}

// Send constructs and sends the BusinessBuyGoods payment request to M-Pesa using shared helper.
func (s *BusinessBuyGoodsService) Send() (map[string]any, error) {
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
		return nil, errors.New("partyB (destination shortcode/merchant) is required")
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
func (s *BusinessBuyGoodsService) ParseCallback(payload map[string]any) (*B2BCallbackResult, error) {
	return ParseB2BCallback(payload)
}

func (s *BusinessBuyGoodsService) getPartyA() string {
	if s.partyA != "" {
		return s.partyA
	}
	return s.Config.GetBusinessCode()
}

// GetResponse returns the last API response stored by the service.
func (s *BusinessBuyGoodsService) GetResponse() map[string]any {
	return s.response
}
