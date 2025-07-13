package Services

import (
	"errors"
	abstracts "github.com/venomous-maker/go-mpesa/Abstracts"
)

// AccountBalanceService handles account balance requests.
type AccountBalanceService struct {
	Config         *abstracts.MpesaConfig
	Client         abstracts.MpesaInterface
	initiator      string
	identifierType string
	remarks        string
}

// NewAccountBalanceService creates a new AccountBalanceService instance.
func NewAccountBalanceService(cfg *abstracts.MpesaConfig, client abstracts.MpesaInterface) *AccountBalanceService {
	return &AccountBalanceService{
		Config: cfg,
		Client: client,
	}
}

// SetInitiator sets the username of the M-Pesa API operator.
func (s *AccountBalanceService) SetInitiator(initiator string) *AccountBalanceService {
	s.initiator = initiator
	return s
}

// SetIdentifierType sets the identifier type for the account balance request.
func (s *AccountBalanceService) SetIdentifierType(identifierType string) *AccountBalanceService {
	s.identifierType = identifierType
	return s
}

// SetRemarks sets any additional information to be associated with the transaction.
func (s *AccountBalanceService) SetRemarks(remarks string) *AccountBalanceService {
	s.remarks = remarks
	return s
}

// AccountBalance initiates an account balance request to the M-Pesa API.
// All parameters are optional. If provided, they override the existing fields.
func (s *AccountBalanceService) AccountBalance(
	initiator, initiatorPassword, partyA, identifierType, remarks, queueURL, resultURL *string,
) (map[string]interface{}, error) {
	if initiator != nil {
		s.SetInitiator(*initiator)
	}
	if remarks != nil {
		s.SetRemarks(*remarks)
	}
	if partyA != nil {
		s.Config.SetBusinessCode(*partyA)
	}
	if identifierType != nil {
		s.SetIdentifierType(*identifierType)
	}
	if queueURL != nil {
		s.Config.SetQueueTimeoutURL(*queueURL)
	}
	if resultURL != nil {
		s.Config.SetResultURL(*resultURL)
	}
	if initiatorPassword != nil {
		if err := s.Config.SetSecurityCredential(*initiatorPassword); err != nil {
			return nil, err
		}
	}

	if s.initiator == "" || s.Config.GetSecurityCredential() == "" || s.Config.GetBusinessCode() == "" {
		return nil, errors.New("initiator, security credential and business code are required")
	}

	requestData := map[string]interface{}{
		"Initiator":          s.initiator,
		"SecurityCredential": s.Config.GetSecurityCredential(),
		"CommandID":          "AccountBalance",
		"PartyA":             s.Config.GetBusinessCode(),
		"IdentifierType":     s.identifierType,
		"Remarks":            s.remarks,
		"QueueTimeOutURL":    s.Config.GetQueueTimeoutURL(),
		"ResultURL":          s.Config.GetResultURL(),
	}

	// Assuming the client returns a map and error
	response, err := s.Client.ExecuteRequest(requestData, "/mpesa/accountbalance/v1/query")
	if err != nil {
		return nil, err
	}

	return response, nil
}
