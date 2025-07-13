package Services

import "github.com/venomous-maker/go-mpesa/Abstracts"

type TransactionStatusService struct {
	*AbstractService

	initiator      string
	transactionID  string
	identifierType string
	remarks        string
	occasion       string
}

func NewTransactionStatusService(cfg *Abstracts.MpesaConfig, client Abstracts.MpesaInterface) *TransactionStatusService {
	return &TransactionStatusService{
		AbstractService: NewAbstractService(cfg, client),
	}
}

func (s *TransactionStatusService) SetInitiator(initiator string) *TransactionStatusService {
	s.initiator = initiator
	return s
}

func (s *TransactionStatusService) SetTransactionID(id string) *TransactionStatusService {
	s.transactionID = id
	return s
}

func (s *TransactionStatusService) SetIdentifierType(idType string) *TransactionStatusService {
	s.identifierType = idType
	return s
}

func (s *TransactionStatusService) SetRemarks(remarks string) *TransactionStatusService {
	s.remarks = remarks
	return s
}

func (s *TransactionStatusService) SetOccasion(occasion string) *TransactionStatusService {
	s.occasion = occasion
	return s
}

// CheckTransactionStatus performs the transaction status query with optional parameters.
func (s *TransactionStatusService) CheckTransactionStatus(
	initiator *string,
	initiatorPassword *string,
	remarks *string,
	partyA *string,
	transactionID *string,
	identifierType *string,
	queueTimeoutURL *string,
	resultURL *string,
	occasion *string,
) (map[string]any, error) {
	// Set optional fields if provided
	if initiator != nil {
		s.SetInitiator(*initiator)
	}
	if remarks != nil {
		s.SetRemarks(*remarks)
	}
	if partyA != nil {
		s.Config.SetBusinessCode(*partyA)
	}
	if transactionID != nil {
		s.SetTransactionID(*transactionID)
	}
	if identifierType != nil {
		s.SetIdentifierType(*identifierType)
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
		s.Config.SetSecurityCredential(*initiatorPassword)
	}

	requestData := map[string]any{
		"Initiator":          s.initiator,
		"SecurityCredential": s.Config.GetSecurityCredential(),
		"CommandID":          "TransactionStatusQuery",
		"TransactionID":      s.transactionID,
		"PartyA":             s.Config.GetBusinessCode(),
		"IdentifierType":     s.identifierType,
		"ResultURL":          s.Config.GetResultURL(),
		"QueueTimeOutURL":    s.Config.GetQueueTimeoutURL(),
		"Remarks":            s.remarks,
		"Occassion":          s.occasion,
	}

	resp, err := s.Client.ExecuteRequest(requestData, "/mpesa/transactionstatus/v1/query")
	if err != nil {
		return nil, err
	}

	s.setResponse(resp)
	return resp, nil
}
