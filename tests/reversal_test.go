package tests

import (
	"errors"
	"testing"

	abstracts "github.com/venomous-maker/go-mpesa/Abstracts"
	"github.com/venomous-maker/go-mpesa/Services"
)

type mockClient struct {
	capturedPayload  any
	capturedEndpoint string
}

func (m *mockClient) ExecuteRequest(payload any, endpoint string) (map[string]any, error) {
	m.capturedPayload = payload
	m.capturedEndpoint = endpoint
	// Simulate success response
	return map[string]any{"ResponseCode": "0"}, nil
}

func buildTestConfig() *abstracts.MpesaConfig {
	cfg, _ := abstracts.NewMpesaConfig("ck", "cs", abstracts.Sandbox, nil, nil, nil, nil, nil)
	cfg.SetBusinessCode("603021")
	cfg.SetQueueTimeoutURL("https://example.com/reversal/queue")
	cfg.SetResultURL("https://example.com/reversal/result")
	cfg.OverrideSecurityCredential("FAKE_SECURITY_CREDENTIAL")
	return cfg
}

func TestReversalService_SuccessReverse(t *testing.T) {
	cfg := buildTestConfig()
	client := &mockClient{}
	service := Services.NewReversalService(cfg, client).
		SetInitiator("apiop37").
		SetTransactionID("PDU91HIVIT").
		SetAmount(200).
		SetReceiverIdentifierType("11").
		SetRemarks("Payment reversal")

	resp, err := service.Reverse()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response, got nil")
	}
	if client.capturedEndpoint != "/mpesa/reversal/v1/request" {
		t.Errorf("unexpected endpoint: %s", client.capturedEndpoint)
	}
	payloadMap, ok := client.capturedPayload.(map[string]interface{})
	if !ok {
		t.Fatalf("expected payload map, got %T", client.capturedPayload)
	}
	// Basic field assertions
	if payloadMap["Amount"] != "200" {
		t.Errorf("expected Amount '200', got %v", payloadMap["Amount"])
	}
	if payloadMap["ReceiverParty"] != "603021" {
		t.Errorf("expected ReceiverParty '603021', got %v", payloadMap["ReceiverParty"])
	}
	if payloadMap["RecieverIdentifierType"] != "11" {
		t.Errorf("expected RecieverIdentifierType '11', got %v", payloadMap["RecieverIdentifierType"])
	}
	if payloadMap["Remarks"] != "Payment reversal" {
		t.Errorf("expected Remarks 'Payment reversal', got %v", payloadMap["Remarks"])
	}
}

func TestReversalService_ValidationErrors(t *testing.T) {
	cfg := buildTestConfig()
	client := &mockClient{}
	service := Services.NewReversalService(cfg, client)

	// Missing initiator
	_, err := service.SetTransactionID("TX123").SetAmount(100).SetReceiverIdentifierType("11").SetRemarks("Test").Reverse()
	if err == nil || err.Error() != "initiator is required" {
		t.Errorf("expected initiator validation error, got %v", err)
	}

	// Missing amount
	service = Services.NewReversalService(cfg, client)
	_, err = service.SetInitiator("user").SetTransactionID("TX123").SetReceiverIdentifierType("11").SetRemarks("Test").Reverse()
	if err == nil || err.Error() != "amount must be greater than 0" {
		t.Errorf("expected amount validation error, got %v", err)
	}

	// Missing remarks
	service = Services.NewReversalService(cfg, client)
	_, err = service.SetInitiator("user").SetTransactionID("TX123").SetAmount(10).SetReceiverIdentifierType("11").Reverse()
	if err == nil || err.Error() != "remarks are required" {
		t.Errorf("expected remarks validation error, got %v", err)
	}

	// Missing security credential (override clears credential for this test)
	cfgNoSec := buildTestConfig()
	cfgNoSec.OverrideSecurityCredential("")
	service = Services.NewReversalService(cfgNoSec, client)
	_, err = service.SetInitiator("user").SetTransactionID("TX123").SetAmount(10).SetReceiverIdentifierType("11").SetRemarks("Test").Reverse()
	if err == nil || err.Error() != "security credential is required; set via SetSecurityCredential or OverrideSecurityCredential on config" {
		t.Errorf("expected security credential validation error, got %v", err)
	}

	// Missing business code
	cfgNoBiz, _ := abstracts.NewMpesaConfig("ck", "cs", abstracts.Sandbox, nil, nil, nil, nil, nil)
	cfgNoBiz.SetQueueTimeoutURL("https://example.com/reversal/queue")
	cfgNoBiz.SetResultURL("https://example.com/reversal/result")
	cfgNoBiz.OverrideSecurityCredential("FAKE")
	service = Services.NewReversalService(cfgNoBiz, client)
	_, err = service.SetInitiator("user").SetTransactionID("TX123").SetAmount(10).SetReceiverIdentifierType("11").SetRemarks("Test").Reverse()
	if err == nil || err.Error() != "business shortcode (ReceiverParty) is required; call SetBusinessCode on mpesa config" {
		t.Errorf("expected business code validation error, got %v", err)
	}

	// Missing queue timeout URL
	cfgNoQueue := buildTestConfig()
	cfgNoQueue.SetQueueTimeoutURL("")
	service = Services.NewReversalService(cfgNoQueue, client)
	_, err = service.SetInitiator("user").SetTransactionID("TX123").SetAmount(10).SetReceiverIdentifierType("11").SetRemarks("Test").Reverse()
	if err == nil || err.Error() != "queue timeout URL is required; call SetQueueTimeoutURL on config" {
		t.Errorf("expected queue timeout URL validation error, got %v", err)
	}

	// Missing result URL
	cfgNoResult := buildTestConfig()
	cfgNoResult.SetResultURL("")
	service = Services.NewReversalService(cfgNoResult, client)
	_, err = service.SetInitiator("user").SetTransactionID("TX123").SetAmount(10).SetReceiverIdentifierType("11").SetRemarks("Test").Reverse()
	if err == nil || err.Error() != "result URL is required; call SetResultURL on config" {
		t.Errorf("expected result URL validation error, got %v", err)
	}

	// Receiver identifier type missing
	cfgOK := buildTestConfig()
	service = Services.NewReversalService(cfgOK, client)
	_, err = service.SetInitiator("user").SetTransactionID("TX123").SetAmount(10).SetRemarks("Test").Reverse()
	if err == nil || !errors.Is(err, errors.New("receiver identifier type is required")) {
		// compare string due to distinct error instances
		if err == nil || err.Error() != "receiver identifier type is required" {
			t.Errorf("expected receiver identifier type error, got %v", err)
		}
	}
}
