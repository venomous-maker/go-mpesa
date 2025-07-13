package Services

import (
	"errors"
	abstracts "github.com/venomous-maker/go-mpesa/Abstracts"
)

// ReversalService handles M-Pesa transaction reversals.
type ReversalService struct {
	Config                 *abstracts.MpesaConfig
	Client                 abstracts.MpesaInterface
	Initiator              string
	TransactionID          string
	ReceiverIdentifierType string
	Remarks                string
	Occasion               string
	Response               map[string]interface{}
}

// NewReversalService creates a new instance of ReversalService.
func NewReversalService(cfg *abstracts.MpesaConfig, client abstracts.MpesaInterface) *ReversalService {
	return &ReversalService{
		Config: cfg,
		Client: client,
	}
}

// SetInitiator sets the API operator initiator name.
func (s *ReversalService) SetInitiator(initiator string) *ReversalService {
	s.Initiator = initiator
	return s
}

// SetTransactionID sets the transaction ID to be reversed.
func (s *ReversalService) SetTransactionID(txID string) *ReversalService {
	s.TransactionID = txID
	return s
}

// SetReceiverIdentifierType sets the type of identifier (ShortCode, MSISDN, TillNumber).
func (s *ReversalService) SetReceiverIdentifierType(identifierType string) *ReversalService {
	s.ReceiverIdentifierType = identifierType
	return s
}

// SetRemarks sets any additional remarks for the reversal.
func (s *ReversalService) SetRemarks(remarks string) *ReversalService {
	s.Remarks = remarks
	return s
}

// SetOccasion sets the transaction occasion.
func (s *ReversalService) SetOccasion(occasion string) *ReversalService {
	s.Occasion = occasion
	return s
}

// Reverse initiates a transaction reversal.
func (s *ReversalService) Reverse(
	initiator *string,
	initiatorPassword *string,
	remarks *string,
	receiverParty *string,
	transactionID *string,
	receiverIdentifierType *string,
	queueTimeoutURL *string,
	resultURL *string,
	occasion *string,
) (map[string]interface{}, error) {
	if initiator != nil {
		s.SetInitiator(*initiator)
	}
	if remarks != nil {
		s.SetRemarks(*remarks)
	}
	if receiverParty != nil {
		s.Config.SetBusinessCode(*receiverParty)
	}
	if transactionID != nil {
		s.SetTransactionID(*transactionID)
	}
	if receiverIdentifierType != nil {
		s.SetReceiverIdentifierType(*receiverIdentifierType)
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
		err := s.Config.SetSecurityCredential(*initiatorPassword)
		if err != nil {
			return nil, err
		}
	}

	if s.TransactionID == "" || s.Initiator == "" || s.ReceiverIdentifierType == "" {
		return nil, errors.New("missing required reversal parameters")
	}

	requestData := map[string]interface{}{
		"Initiator":              s.Initiator,
		"SecurityCredential":     s.Config.GetSecurityCredential(),
		"CommandID":              "TransactionReversal",
		"TransactionID":          s.TransactionID,
		"ReceiverParty":          s.Config.GetBusinessCode(),
		"ReceiverIdentifierType": s.ReceiverIdentifierType,
		"ResultURL":              s.Config.GetResultURL(),
		"QueueTimeOutURL":        s.Config.GetQueueTimeoutURL(),
		"Remarks":                s.Remarks,
		"Occassion":              s.Occasion,
	}

	resp, err := s.Client.ExecuteRequest(requestData, "/mpesa/reversal/v1/request")
	if err != nil {
		return nil, err
	}

	s.Response = resp
	return resp, nil
}
