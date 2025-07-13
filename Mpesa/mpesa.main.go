package Mpesa

import (
	"github.com/venomous-maker/go-mpesa/Abstracts"
	"github.com/venomous-maker/go-mpesa/Services"
)

type Mpesa struct {
	Config *Abstracts.MpesaConfig
	Client *Abstracts.ApiClient
}

// New creates a new Mpesa instance with optional credentials and environment
func New(consumerKey, consumerSecret, environment string) (*Mpesa, error) {
	cfg, err := Abstracts.NewMpesaConfig(
		consumerKey,
		consumerSecret,
		Abstracts.Environment(environment),
		nil, nil, nil, nil, nil,
	)
	if err != nil {
		return nil, err
	}

	client := Abstracts.NewApiClient(cfg)

	return &Mpesa{
		Config: cfg,
		Client: client,
	}, nil
}

// SetCredentials updates the credentials and returns the modified Mpesa instance
func (m *Mpesa) SetCredentials(consumerKey, consumerSecret, environment string) error {
	cfg, err := Abstracts.NewMpesaConfig(
		consumerKey,
		consumerSecret,
		Abstracts.Environment(environment),
		nil, nil, nil, nil, nil,
	)
	if err != nil {
		return err
	}
	m.Config = cfg
	m.Client = Abstracts.NewApiClient(cfg)
	return nil
}

// SetBusinessCode sets the business shortcode
func (m *Mpesa) SetBusinessCode(code string) {
	m.Config.SetBusinessCode(code)
}

// SetPassKey sets the passkey
func (m *Mpesa) SetPassKey(key string) {
	m.Config.SetPassKey(key)
}

// Stk returns a new STK Push service
func (m *Mpesa) Stk() *Services.StkService {
	return Services.NewStkService(m.Config, m.Client)
}

// CustomerToBusiness returns a new C2B service (stub, add implementation)
func (m *Mpesa) CustomerToBusiness() *Services.CustomerToBusinessService {
	return Services.NewCustomerToBusinessService(m.Config, m.Client)
}

// BusinessToCustomer returns a new B2C service (stub, add implementation)
func (m *Mpesa) BusinessToCustomer() *Services.BusinessToCustomerService {
	return Services.NewBusinessToCustomerService(m.Config, m.Client)
}

// AccountBalance returns a new account balance service (stub, add implementation)
func (m *Mpesa) AccountBalance() *Services.AccountBalanceService {
	return Services.NewAccountBalanceService(m.Config, m.Client)
}

// TransactionStatus returns a new transaction status service (stub, add implementation)
func (m *Mpesa) TransactionStatus() *Services.TransactionStatusService {
	return Services.NewTransactionStatusService(m.Config, m.Client)
}

// Reversal returns a new reversal service (stub, add implementation)
func (m *Mpesa) Reversal() *Services.ReversalService {
	return Services.NewReversalService(m.Config, m.Client)
}
