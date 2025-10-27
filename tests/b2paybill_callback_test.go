package tests

import (
	"testing"

	"github.com/venomous-maker/go-mpesa/Services"
)

func TestParseCallback_Success(t *testing.T) {
	svc := Services.NewBusinessToPayBillService(nil, nil)
	payload := map[string]any{
		"Result": map[string]any{
			"ResultType":               "0",
			"ResultCode":               "0",
			"ResultDesc":               "The service request is processed successfully",
			"OriginatorConversationID": "626f6ddf-ab37-4650-b882-b1de92ec9aa4",
			"ConversationID":           "12345677dfdf89099B3",
			"TransactionID":            "QKA81LK5CY",
			"ResultParameters": map[string]any{
				"ResultParameter": []any{
					map[string]any{"Key": "Amount", "Value": "190.00"},
					map[string]any{"Key": "TransCompletedTime", "Value": "20221110110717"},
				},
			},
			"ReferenceData": map[string]any{
				"ReferenceItem": []any{
					map[string]any{"Key": "BillReferenceNumber", "Value": "19008"},
				},
			},
		},
	}
	res, err := svc.ParseCallback(payload)
	if err != nil {
		t.Fatalf("ParseCallback error: %v", err)
	}
	if !res.Success {
		t.Fatalf("expected success true got false, code=%s", res.ResultCode)
	}
	if res.TransactionID != "QKA81LK5CY" {
		t.Fatalf("unexpected tx id: %s", res.TransactionID)
	}
	if val, ok := res.ResultParameters["Amount"]; !ok || val != "190.00" {
		t.Fatalf("Amount not parsed: %v", res.ResultParameters)
	}
	if val, ok := res.ReferenceData["BillReferenceNumber"]; !ok || val != "19008" {
		t.Fatalf("Ref not parsed: %v", res.ReferenceData)
	}
}

func TestParseCallback_Failure(t *testing.T) {
	svc := Services.NewBusinessToPayBillService(nil, nil)
	payload := map[string]any{
		"Result": map[string]any{
			"ResultType":               0,
			"ResultCode":               2001,
			"ResultDesc":               "The initiator information is invalid.",
			"OriginatorConversationID": "12337-23509183-5",
			"ConversationID":           "AG_20200120_0000657265d5fa9ae5c0",
			"TransactionID":            "OAK0000000",
			"ResultParameters": map[string]any{
				"ResultParameter": map[string]any{"Key": "BOCompletedTime", "Value": 20200120164825},
			},
			"ReferenceData": map[string]any{
				"ReferenceItem": map[string]any{"Key": "QueueTimeoutURL", "Value": "https://internalapi.safaricom.co.ke/mpesa/abresults/v1/submit"},
			},
		},
	}
	res, err := svc.ParseCallback(payload)
	if err != nil {
		t.Fatalf("ParseCallback error: %v", err)
	}
	if res.Success {
		t.Fatalf("expected success false got true, code=%s", res.ResultCode)
	}
	if res.ResultCode != "2001" {
		t.Fatalf("expected code 2001 got %s", res.ResultCode)
	}
	if val, ok := res.ResultParameters["BOCompletedTime"]; !ok || val != "20200120164825" {
		t.Fatalf("BOCompletedTime not parsed: %v", res.ResultParameters)
	}
}
