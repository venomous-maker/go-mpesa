package Services

import "github.com/venomous-maker/go-mpesa/Abstracts"

type AbstractService struct {
	*BaseService
	response map[string]any
}

// NewAbstractService initializes an abstract service with config and client
func NewAbstractService(cfg *Abstracts.MpesaConfig, client Abstracts.MpesaInterface) *AbstractService {
	return &AbstractService{
		BaseService: NewBaseService(cfg, client),
	}
}

// GetResponse returns the last response as a generic map
func (s *AbstractService) GetResponse() map[string]any {
	return s.response
}

// setResponse is a helper to store the last response
func (s *AbstractService) setResponse(resp map[string]any) {
	s.response = resp
}
