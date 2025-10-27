package Services

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	abstracts "github.com/venomous-maker/go-mpesa/Abstracts"
)

// B2BRequest represents a generic B2B payment request.
type B2BRequest struct {
	Initiator              string
	SecurityCredential     string
	CommandID              string
	SenderIdentifierType   string
	RecieverIdentifierType string
	Amount                 float64
	PartyA                 string
	PartyB                 string
	AccountReference       string
	Requester              string
	Remarks                string
	QueueTimeOutURL        string
	ResultURL              string
	Occasion               string
}

// ExecuteB2BRequest builds the request payload from B2BRequest and executes the API call.
func ExecuteB2BRequest(cfg *abstracts.MpesaConfig, client abstracts.MpesaInterface, req B2BRequest) (map[string]any, error) {
	if cfg == nil || client == nil {
		return nil, errors.New("cfg and client are required")
	}

	if req.Initiator == "" {
		return nil, errors.New("initiator is required")
	}
	if req.SecurityCredential == "" {
		return nil, errors.New("security credential is required")
	}
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	if req.PartyA == "" && cfg.GetBusinessCode() == "" {
		return nil, errors.New("partyA (business shortcode) is required")
	}
	if req.PartyB == "" {
		return nil, errors.New("partyB (destination shortcode) is required")
	}

	payload := map[string]any{
		"Initiator":              req.Initiator,
		"SecurityCredential":     req.SecurityCredential,
		"CommandID":              req.CommandID,
		"SenderIdentifierType":   req.SenderIdentifierType,
		"RecieverIdentifierType": req.RecieverIdentifierType,
		"Amount":                 math.Round(req.Amount),
		"PartyA":                 choosePartyA(req.PartyA, cfg),
		"PartyB":                 req.PartyB,
		"AccountReference":       req.AccountReference,
		"Requester":              req.Requester,
		"Remarks":                req.Remarks,
		"QueueTimeOutURL":        chooseString(req.QueueTimeOutURL, cfg.GetQueueTimeoutURL()),
		"ResultURL":              chooseString(req.ResultURL, cfg.GetResultURL()),
		"Occasion":               req.Occasion,
	}

	return client.ExecuteRequest(payload, "/mpesa/b2b/v1/paymentrequest")
}

func choosePartyA(partyA string, cfg *abstracts.MpesaConfig) string {
	if partyA != "" {
		return partyA
	}
	return cfg.GetBusinessCode()
}

func chooseString(incoming, fallback string) string {
	if incoming != "" {
		return incoming
	}
	return fallback
}

// B2BCallbackResult represents a parsed B2B callback payload shared across B2B services.
type B2BCallbackResult struct {
	ResultCode               string
	ResultDesc               string
	TransactionID            string
	OriginatorConversationID string
	ConversationID           string
	ResultParameters         map[string]string
	ReferenceData            map[string]string
	Raw                      map[string]any
	Success                  bool
}

// ParseB2BCallback parses a generic B2B callback payload (PayBill, BuyGoods, etc.).
func ParseB2BCallback(payload map[string]any) (*B2BCallbackResult, error) {
	res := &B2BCallbackResult{
		ResultParameters: make(map[string]string),
		ReferenceData:    make(map[string]string),
		Raw:              payload,
	}

	var resultNode any
	if v, ok := payload["Result"]; ok {
		resultNode = v
	} else if v, ok := payload["result"]; ok {
		resultNode = v
	} else {
		return nil, errors.New("payload missing Result node")
	}

	resultMap, ok := resultNode.(map[string]any)
	if !ok {
		if m2, ok2 := resultNode.(map[string]interface{}); ok2 {
			resultMap = make(map[string]any)
			for k, v := range m2 {
				resultMap[k] = v
			}
		} else {
			return nil, errors.New("Result node is not an object")
		}
	}

	res.ResultCode = toString(resultMap["ResultCode"]) // may be string or number
	res.ResultDesc = toString(resultMap["ResultDesc"])
	res.TransactionID = toString(resultMap["TransactionID"])
	res.OriginatorConversationID = toString(resultMap["OriginatorConversationID"])
	res.ConversationID = toString(resultMap["ConversationID"])

	// Parse ResultParameters -> ResultParameter (array or single)
	if rpRaw, ok := resultMap["ResultParameters"]; ok {
		switch rp := rpRaw.(type) {
		case map[string]any:
			parseResultParameterArray(rp["ResultParameter"], res.ResultParameters)
		default:
			parseResultParameterArray(rpRaw, res.ResultParameters)
		}
	}

	// Parse ReferenceData -> ReferenceItem
	if rdRaw, ok := resultMap["ReferenceData"]; ok {
		switch rd := rdRaw.(type) {
		case map[string]any:
			parseReferenceItem(rd["ReferenceItem"], res.ReferenceData)
		default:
			parseReferenceItem(rdRaw, res.ReferenceData)
		}
	}

	if i, err := strconv.Atoi(res.ResultCode); err == nil {
		res.Success = i == 0
	} else {
		res.Success = res.ResultCode == "0"
	}

	return res, nil
}

// helpers (copied from previous implementation)
func toString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		if t == math.Trunc(t) {
			return strconv.FormatInt(int64(t), 10)
		}
		return strconv.FormatFloat(t, 'f', -1, 64)
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	case nil:
		return ""
	default:
		return fmt.Sprint(t)
	}
}

func parseResultParameterArray(input any, out map[string]string) {
	if input == nil {
		return
	}

	if arr, ok := input.([]any); ok {
		for _, item := range arr {
			if m, ok := item.(map[string]any); ok {
				k := toString(m["Key"])
				v := toString(m["Value"])
				if k != "" {
					out[k] = v
				}
			}
		}
		return
	}

	if arr2, ok := input.([]interface{}); ok {
		for _, item := range arr2 {
			if m, ok := item.(map[string]interface{}); ok {
				k := toString(m["Key"])
				v := toString(m["Value"])
				if k != "" {
					out[k] = v
				}
			}
		}
		return
	}

	if m, ok := input.(map[string]any); ok {
		k := toString(m["Key"])
		v := toString(m["Value"])
		if k != "" {
			out[k] = v
		}
		return
	}
	if m2, ok := input.(map[string]interface{}); ok {
		k := toString(m2["Key"])
		v := toString(m2["Value"])
		if k != "" {
			out[k] = v
		}
		return
	}
}

func parseReferenceItem(input any, out map[string]string) {
	if input == nil {
		return
	}

	if arr, ok := input.([]any); ok {
		for _, item := range arr {
			if m, ok := item.(map[string]any); ok {
				k := toString(m["Key"])
				v := toString(m["Value"])
				if k != "" {
					out[k] = v
				}
			}
		}
		return
	}

	if arr2, ok := input.([]interface{}); ok {
		for _, item := range arr2 {
			if m, ok := item.(map[string]interface{}); ok {
				k := toString(m["Key"])
				v := toString(m["Value"])
				if k != "" {
					out[k] = v
				}
			}
		}
		return
	}

	if m, ok := input.(map[string]any); ok {
		k := toString(m["Key"])
		v := toString(m["Value"])
		if k != "" {
			out[k] = v
		}
		return
	}

	if m2, ok := input.(map[string]interface{}); ok {
		k := toString(m2["Key"])
		v := toString(m2["Value"])
		if k != "" {
			out[k] = v
		}
		return
	}
}
